package main

import (
	"log"
	"net/http"

	"github.com/ZupIT/go-vault-session/pkg/login"
	"github.com/ZupIT/go-vault-session/pkg/token"
	"github.com/hashicorp/vault/api"

	"github.com/kaduartur/go-vault-demo/server"
	"github.com/kaduartur/go-vault-demo/server/database"
	"github.com/kaduartur/go-vault-demo/server/database/repository"
	"github.com/kaduartur/go-vault-demo/server/http/creditcard"
	"github.com/kaduartur/go-vault-demo/server/http/payment"
	"github.com/kaduartur/go-vault-demo/server/vault"
)

var creditCardHandle server.CreditCardHandler
var paymentHandle server.PaymentHandler

func init() {
	client := vaultConfig()
	vaultStarter(client)
	vaultManager := vault.NewManager(client)

	pgManager := database.NewPgManager()
	cardRepository := repository.NewCardRepository(pgManager)

	creditCardHandle = creditcard.NewHandler(vaultManager, cardRepository)
	paymentHandle = payment.NewHandler(vaultManager)
}

func main() {
	log.Println("Starting server")
	http.Handle("/credit-cards", creditCardHandle.HandleCard())
	http.Handle("/credit-cards/", creditCardHandle.HandleCard())
	http.Handle("/payments", paymentHandle.HandlePay())
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func vaultConfig() *api.Client {
	vaultConfig := api.DefaultConfig()
	_ = vaultConfig.ReadEnvironment()
	client, _ := api.NewClient(vaultConfig)
	return client
}

func vaultStarter(client *api.Client) {
	vaultAuth := login.NewHandler(client)
	secret := vaultAuth.HandleLogin()

	renewal := token.NewRenewalHandler(client, secret)
	renewal.HandleRenewal()
}
