package creditcard

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/kaduartur/go-vault-demo/server"
	"github.com/kaduartur/go-vault-demo/server/database/repository"
)

type Handler struct {
	vaultManger    server.VaultManager
	cardRepository repository.CreditCard
}

func NewHandler(v server.VaultManager, r repository.CreditCard) *Handler {
	return &Handler{vaultManger: v, cardRepository: r}
}

func (h *Handler) HandleCard() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			h.processPost(w, r)
		case http.MethodGet:
			h.processGet(w, r)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})
}

func (h *Handler) processPost(w http.ResponseWriter, r *http.Request) {
	var c server.CreditCard
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		log.Println("Failed to process request", err)
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err.Error())
		return
	}

	cipher, err := h.vaultManger.EncryptCardNumber(c.Number)
	if err != nil {
		log.Println(err)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	lastDigits := c.Number[len(c.Number)-4:]
	bin := c.Number[0:6]

	cardEntity := repository.CreditCardEntity{
		Holder:         c.Holder,
		Number:         cipher,
		Brand:          c.Brand,
		Bin:            bin,
		LastDigits:     lastDigits,
		ExpirationDate: expirationDate(c.ExpirationDate),
		CreateAt:       time.Now(),
	}

	cardId, err := h.cardRepository.Save(cardEntity)
	if err != nil {
		log.Println(err)
		_, _ = w.Write([]byte("Unknown error"))
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	resp := make(map[string]interface{})
	resp["cardId"] = cardId

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
}

func (h *Handler) processGet(w http.ResponseWriter, r *http.Request) {
	cardId := cardId(r.URL.Path)

	card, err := h.cardRepository.FindById(cardId)
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

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(card)
}

func cardId(path string) string {
	return strings.Replace(path, "/credit-cards/", "", 1)
}

func expirationDate(d string) time.Time {
	split := strings.Split(d, "/")
	month, _ := strconv.Atoi(split[0])
	year, _ := strconv.Atoi(split[1])

	return time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
}
