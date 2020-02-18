package payment

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/hashicorp/go-uuid"

	"github.com/kaduartur/go-vault-demo/server"
	"github.com/kaduartur/go-vault-demo/server/database/repository"
)

type Handler struct {
	vaultManager   server.VaultManager
	kafkaManager   server.KafkaManager
	cardRepository repository.CreditCard
}

func NewHandler(v server.VaultManager, k server.KafkaManager, r repository.CreditCard) *Handler {
	return &Handler{vaultManager: v, kafkaManager: k, cardRepository: r}
}

func (h *Handler) HandlePay() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			h.processPost(w, r)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})
}

func (h *Handler) processPost(w http.ResponseWriter, r *http.Request) {
	var p server.PaymentRequest
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		log.Println("Failed to process request", err)
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err.Error())
		return
	}

	card, err := h.cardRepository.FindById(p.CardId)
	if err != nil {
		log.Println(err)
		_, _ = w.Write([]byte("Unknown error"))
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	if card == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	cvvCipher, err := h.vaultManager.Encrypt(p.Cvv)
	if err != nil {
		log.Println(err)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		log.Println(err)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	payId := fmt.Sprintf("PAY-%s", id)

	event := server.PaymentEvent{
		PaymentId:      payId,
		CardId:         card.CardId,
		Holder:         card.Holder,
		Number:         card.Number,
		Cvv:            cvvCipher,
		Brand:          card.Brand,
		ExpirationDate: card.ExpirationDate,
		Amount:         p.Amount,
		Description:    p.Description,
	}

	err = h.kafkaManager.SendEvent(event)
	if err != nil {
		log.Println(err)
		_, _ = w.Write([]byte("Payment failed"))
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	resp := make(map[string]interface{})
	resp["paymentId"] = payId
	resp["status"] = "PROCESSING"

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)

}
