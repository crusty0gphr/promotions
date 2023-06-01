package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gorilla/mux"

	"github.com/promotions"
)

const (
	host = "0.0.0.0"
	port = "7700"
)

func main() {
	logger := log.New(os.Stdout, "app: ", 1)

	srv := promotions.NewService(logger)

	router := mux.NewRouter()
	if err := promotions.MakeHandlers(
		logger,
		router,
		fmt.Sprintf("%s:%s", host, port),
		srv,
	); err != nil {
		logger.Fatal(err)
	}
}
