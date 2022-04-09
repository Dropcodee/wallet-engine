package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Dropcodee/wallet-engine/pkg/models"
	"github.com/Dropcodee/wallet-engine/pkg/utils"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

var newWallet models.Wallet

type ErrorResponse struct {
	Stack   *gorm.DB
	Code    uint
	Message string
}

func RetrieveAllWalletsHandler(res http.ResponseWriter, req *http.Request) {
	allBooks := models.GetWallets()
	jsonResponse, _ := json.Marshal(allBooks)
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(jsonResponse)
}

func RetrieveWalletHandler(res http.ResponseWriter, req *http.Request) {
	reqParams := mux.Vars(req)
	walletIdentity := reqParams["walletIdentity"]

	wallet, _ := models.GetWalletByIdentity(walletIdentity)
	// convert to json for response
	jsonResponse, _ := json.Marshal(wallet)

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(jsonResponse)
}

func CreateWalletHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	wallet := &models.Wallet{}
	// create wallet identity key
	wallet.Identity = uuid.New().String()
	wallet.Status = false
	wallet.StatusLabel = "DISABLED"

	utils.ParseBody(req, wallet)
	err := wallet.ValidateWallet()
	if err != nil {
		jsonError, _ := json.Marshal(err)
		res.WriteHeader(http.StatusUnprocessableEntity)
		res.Write(jsonError)
		return
	}
	wal := wallet.CreateWallet()

	jsonRes, _ := json.Marshal(wal)
	res.WriteHeader(http.StatusCreated)
	res.Write(jsonRes)
}

func ActivateWalletHandler(res http.ResponseWriter, req *http.Request) {
	reqParams := mux.Vars(req)
	walletIdentity := reqParams["walletIdentity"]

	wallet, _ := models.GetWalletByIdentity(walletIdentity)

	wallet.Status = true
	wallet.StatusLabel = "ACTIVE"
	updatedWallet := models.ToggleWallet(wallet)

	jsonResponse, _ := json.Marshal(updatedWallet)

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(jsonResponse)

}
func DeactivateWalletHandler(res http.ResponseWriter, req *http.Request) {
	reqParams := mux.Vars(req)
	walletIdentity := reqParams["walletIdentity"]

	wallet, _ := models.GetWalletByIdentity(walletIdentity)

	wallet.Status = false
	wallet.StatusLabel = "DISABLED"
	updatedWallet := models.ToggleWallet(wallet)

	jsonResponse, _ := json.Marshal(updatedWallet)

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(jsonResponse)

}

func CreateWalletTransaction(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	transaction := &models.Transaction{}
	// parse data to golang format
	utils.ParseBody(req, transaction)
	err := transaction.ValidateTransaction()
	if err != nil {
		jsonError, _ := json.Marshal(err)
		res.WriteHeader(http.StatusUnprocessableEntity)
		res.Write(jsonError)
		return
	}
	// check if wallet exists
	wallet, db := models.GetWalletByIdentity(transaction.WalletIdentity)
	// check if wallet is activated before creating a transaction & if not return error
	if !(wallet.Status) {
		res.WriteHeader(http.StatusUnprocessableEntity)
		response := ErrorResponse{Code: http.StatusUnprocessableEntity, Message: "Please activate wallet to make a transaction"}
		jsonResponse, _ := json.Marshal(response)
		res.Write(jsonResponse)
		return
	}

	// after above checks meaning wallet exists and is active
	// check type of transaction if debit then validate if wallet has enough balance to deduct from it
	if transaction.Type == "debit" {
		// if balance is less than the transaction amount return error
		if wallet.Total < transaction.Amount {
			res.WriteHeader(http.StatusUnprocessableEntity)
			response := ErrorResponse{Stack: nil, Code: http.StatusUnprocessableEntity, Message: "Insufficient Wallet Balance."}
			jsonResponse, _ := json.Marshal(response)
			res.Write(jsonResponse)
			return
		}
		// create transaction
		trans := transaction.CreateTransaction()

		// update wallet total for next transaction
		wallet.Total = wallet.Total - transaction.Amount
		db.Save(&wallet)
		jsonResponse, _ := json.Marshal(trans)
		res.WriteHeader(http.StatusCreated)
		res.Write(jsonResponse)
		return
	}

	// if credit then update total and save Transaction
	if transaction.Type == "credit" {
		utils.ParseBody(req, transaction)
		// create transaction
		trans := transaction.CreateTransaction()

		// update wallet total for next transaction
		wallet.Total = wallet.Total + transaction.Amount
		db.Save(&wallet)
		jsonResponse, _ := json.Marshal(trans)
		res.WriteHeader(http.StatusCreated)
		res.Write(jsonResponse)
		return
	}
	res.WriteHeader(http.StatusUnprocessableEntity)
	response := ErrorResponse{Stack: nil, Code: http.StatusUnprocessableEntity, Message: "Please check your transaction type again it should be (debit or credit)."}
	jsonResponse, _ := json.Marshal(response)
	res.Write(jsonResponse)
}
