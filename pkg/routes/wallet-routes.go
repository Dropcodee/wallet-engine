package routes

import (
	"github.com/Dropcodee/wallet-engine/pkg/handlers"
	"github.com/gorilla/mux"
)

var WalletRoutes = func(router *mux.Router) {
	// creates a wallet identity
	// endpoint localhost:9010/wallets/
	// eg payload
	/*
		{
		    "username": "dropcode",
		    "total": 100000.50
		}
	*/
	router.HandleFunc("/wallets/", handlers.CreateWalletHandler).Methods("POST")
	// gets all created wallets
	// endpoint localhost:9010/wallets/
	router.HandleFunc("/wallets/", handlers.RetrieveAllWalletsHandler).Methods("GET")
	// gets single wallet by wallet identity
	// e.g endpoint localhost:9010/wallets/7f5980f2-603e-45ae-93c6-d43ccc603e65/
	router.HandleFunc("/wallets/{walletIdentity}/", handlers.RetrieveWalletHandler).Methods("GET")
	// activate a wallet just pass the wallet identity to the url
	// eg wallet identity "7f5980f2-603e-45ae-93c6-d43ccc603e65"
	// e.g endpoint localhost:9010/wallets/activate/7f5980f2-603e-45ae-93c6-d43ccc603e65/
	router.HandleFunc("/wallets/activate/{walletIdentity}/", handlers.ActivateWalletHandler).Methods("PATCH")
	// deactivate a wallet just pass the wallet identity to the url
	// eg wallet identity "7f5980f2-603e-45ae-93c6-d43ccc603e65"
	router.HandleFunc("/wallets/deactivate/{walletIdentity}/", handlers.DeactivateWalletHandler).Methods("PATCH")

	// creates a transaction for a wallet either credit/debit
	// eg payload
	// e.g endpoint localhost:9010/wallets/transactions/create/
	/*
		{
		    "amount": 10000, should be in number not string
		    "transactionType": "debit", or "transactionType": "credit"
		    "walletIdentity": "7f5980f2-603e-45ae-93c6-d43ccc603e65"
		}
	*/
	router.HandleFunc("/wallets/transactions/create/", handlers.CreateWalletTransaction).Methods("POST")
}
