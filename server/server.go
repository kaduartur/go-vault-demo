package server

import (
	"net/http"
	"time"
)

type CreditCardRequest struct {
	Number         string `json:"number"`
	Holder         string `json:"holder"`
	ExpirationDate string `json:"expirationDate"`
	Brand          string `json:"brand"`
}

type PaymentRequest struct {
	CardId      string `json:"cardId"`
	Amount      int    `json:"amount"`
	Description string `json:"description"`
	Cvv         string `json:"cvv"`
}

type PaymentEvent struct {
	PaymentId      string    `json:"paymentId"`
	CardId         string    `json:"cardId"`
	Holder         string    `json:"holder"`
	Number         string    `json:"number"`
	Cvv            string    `json:"cvv"`
	Brand          string    `json:"brand"`
	ExpirationDate time.Time `json:"expirationDate"`
	Amount         int       `json:"amount"`
	Description    string    `json:"description"`
}

type CreditCardHandler interface {
	HandleCard() http.Handler
}

type PaymentHandler interface {
	HandlePay() http.Handler
}

type VaultManager interface {
	Write(key string, data map[string]interface{}) error
	Read(key string) (map[string]interface{}, error)
	Delete(key string) error
	Encrypt(data string) (string, error)
}

type KafkaManager interface {
	SendEvent(p PaymentEvent) error
}
