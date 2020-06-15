package main

import (
	"fmt"
	"log"
	"net/http"
	"notes/models"
	"notes/router"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	var (
		port             = os.Getenv("PORT")
		connectionString = mustGetenv("CONNECTION_STRING")
	)

	models.InitDB(connectionString)
	router := router.Router()

	log.Printf("Starting server on port %s...", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}

func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Panicf("%s environment variable not set.", k)
	}
	return v
}
