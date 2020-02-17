package repository

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-uuid"

	"github.com/kaduartur/go-vault-demo/server/database"
)

type CreditCard interface {
	Save(card CreditCardEntity) (string, error)
	FindById(id string) (*CreditCardEntity, error)
}

type CreditCardEntity struct {
	CardId         string     `json:"cardId"`
	Holder         string     `json:"holder"`
	Number         string     `json:"number"`
	Brand          string     `json:"brand"`
	Bin            string     `json:"bin"`
	LastDigits     string     `json:"lastDigits"`
	ExpirationDate time.Time  `json:"expirationDate"`
	CreateAt       time.Time  `json:"createAt"`
	UpdateAt       *time.Time `json:"updateAt,omitempty"`
}

type creditCard struct {
	database database.DbConnection
}

func NewCardRepository(d database.DbConnection) *creditCard {
	return &creditCard{database: d}
}

func (c *creditCard) Save(card CreditCardEntity) (string, error) {
	statement := `INSERT INTO public.credit_card 
				  (credit_card_id, "number", holder, expiration_date, brand, bin, last_digits, created)
				  VALUES($1, $2, $3, $4, $5, $6, $7, $8);`

	db := c.database.ConnectHandle()
	defer db.Close()

	id, err := uuid.GenerateUUID()
	if err != nil {
		return "", err
	}
	cardId := fmt.Sprintf("CARD-%s", id)
	_, err = db.Exec(statement, cardId, card.Number, card.Holder, card.ExpirationDate, card.Brand, card.Bin, card.LastDigits, card.CreateAt)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return cardId, nil
}

func (c *creditCard) FindById(id string) (*CreditCardEntity, error) {
	statement := "SELECT * FROM public.credit_card WHERE credit_card_id=$1;"

	db := c.database.ConnectHandle()
	defer db.Close()

	var card CreditCardEntity
	row := db.QueryRow(statement, id)
	switch err := row.Scan(&card.CardId, &card.Number, &card.Holder, &card.ExpirationDate, &card.Brand, &card.Bin, &card.LastDigits, &card.CreateAt, &card.UpdateAt); err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		return &card, nil
	default:
		return nil, err
	}
}
