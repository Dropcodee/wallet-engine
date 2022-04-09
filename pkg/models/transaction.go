package models

import (
	"github.com/Dropcodee/wallet-engine/pkg/config"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/jinzhu/gorm"
)

type Transaction struct {
	gorm.Model
	Amount         float64 `json:"amount"`
	Type           string  `json:"transactionType"`
	WalletIdentity string  `json:"walletIdentity"`
}

func init() {
	config.ConnectDb()
	db = config.GetDB()
	db.AutoMigrate(&Transaction{})
}

func (trans *Transaction) CreateTransaction() *Transaction {
	db.NewRecord(trans)
	db.Create(&trans)
	return trans
}

func (trans Transaction) ValidateTransaction() error {
	return validation.ValidateStruct(&trans,
		validation.Field(&trans.Amount, validation.Required.Error("transaction amount is required & should be either a number or decimal")),
		validation.Field(&trans.Type, validation.Required.Error("Transaction Type must be either credit or debit and cannot be empty.")),
		validation.Field(&trans.WalletIdentity, validation.Required.Error("wallet identity code cannot be empty.")),
	)
}
