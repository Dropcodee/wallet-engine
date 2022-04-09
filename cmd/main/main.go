package main

import (
	"log"
	"net/http"

	"github.com/Dropcodee/wallet-engine/pkg/routes"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	routes.WalletRoutes(router)
	http.Handle("/", router)
	log.Fatal(http.ListenAndServe("localhost:9010", router))
}
