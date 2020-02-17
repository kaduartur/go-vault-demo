package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/ZupIT/go-vault-session/pkg/login"
	"github.com/ZupIT/go-vault-session/pkg/token"
	"github.com/hashicorp/vault/api"
	kf "github.com/segmentio/kafka-go"

	"github.com/kaduartur/go-vault-demo/server"
	"github.com/kaduartur/go-vault-demo/server/database"
	"github.com/kaduartur/go-vault-demo/server/database/repository"
	"github.com/kaduartur/go-vault-demo/server/http/creditcard"
	"github.com/kaduartur/go-vault-demo/server/http/payment"
	"github.com/kaduartur/go-vault-demo/server/kafka"
	"github.com/kaduartur/go-vault-demo/server/vault"
)

var creditCardHandle server.CreditCardHandler
var paymentHandle server.PaymentHandler

func init() {
	client := vaultConfig()
	vaultStarter(client)
	vaultManager := vault.NewManager(client)

	cfg := kafkaConfig()
	kafkaManager := kafka.NewManager(cfg)

	dbManager := vault.NewDatabaseManager(client)
	dbManager.CreateDynamicCredential()

	pgManager := database.NewPgManager()
	cardRepository := repository.NewCardRepository(pgManager)

	creditCardHandle = creditcard.NewHandler(vaultManager, cardRepository)
	paymentHandle = payment.NewHandler(vaultManager, kafkaManager, cardRepository)
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

func kafkaConfig() kf.WriterConfig {
	return kf.WriterConfig{
		Brokers:  strings.Split(os.Getenv("KAFKA_BOOTSTRAP_SERVERS"), ","),
		Topic:    "payment_events",
		Balancer: &kf.LeastBytes{},
	}
}
