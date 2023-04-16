package main

import (
	"os"

	"github.com/Touchsung/money-note-line-api-go/router"
)

func main() {
    r := router.Router()
	port := os.Getenv("PORT")
	if port == "" {
        port = "8080"
    }
    r.Logger.Fatal(r.Start(":" + port))
}
