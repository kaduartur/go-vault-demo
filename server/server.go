package server

import (
	"net/http"
)

type CreditCard struct {
	Number         string `json:"number"`
	Holder         string `json:"holder"`
	ExpirationDate string `json:"expirationDate"`
	Brand          string `json:"brand"`
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
	EncryptCardNumber(num string) (string, error)
}


