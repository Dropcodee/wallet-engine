package models

import (
	"github.com/Dropcodee/wallet-engine/pkg/config"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB

type Wallet struct {
	gorm.Model
	Username     string        `json:"username" gorm:"unique"`
	Identity     string        `json:"identity" gorm:"unique_index,unique"`
	Total        float64       `json:"total"`
	Status       bool          `json:"status"`
	StatusLabel  string        `json:"statusLabel" `
	Transactions []Transaction `gorm:"foreignKey:WalletIdentity, references:Identity"`
}

func init() {
	config.ConnectDb()
	db = config.GetDB()
	db.AutoMigrate(&Wallet{})
}

// database services functions
// create wallet db abstraction function
// this creates wallets
func (wal *Wallet) CreateWallet() *Wallet {
	db.NewRecord(wal)
	db.Create(&wal)
	return wal
}

func GetWallets() []Wallet {
	var wallets []Wallet
	db.Find(&wallets)
	return wallets
}

func GetWalletByIdentity(identity string) (*Wallet, *gorm.DB) {
	var getWallet Wallet
	db := db.Where("Identity=?", identity).Find(&getWallet)
	return &getWallet, db
}

func ToggleWallet(reqWallet *Wallet) *Wallet {
	db.Save(&reqWallet)
	return reqWallet
}

func (wal Wallet) ValidateWallet() error {
	return validation.ValidateStruct(&wal,
		validation.Field(&wal.Username, validation.Required.Error("wallet name is required"), validation.Length(3, 50)),
		validation.Field(&wal.Total, validation.Required.Error("Initial wallet amount is required and must be of type int or float")),
	)
}
