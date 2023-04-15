package main

import (
	"log"
	"os"

	"github.com/Touchsung/money-note-line-api-go/router"
	"github.com/joho/godotenv"
)

func main() {
    r := router.Router()
	err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }
	port := os.Getenv("PORT")
	if port == "" {
        port = "8080"
    }
    r.Logger.Fatal(r.Start(":" + port))
}
