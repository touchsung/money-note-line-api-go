package main

import (
	"github.com/Touchsung/money-note-line-api-go/config"
	"github.com/Touchsung/money-note-line-api-go/router"
)

func main() {
    r := router.Router()
	  port := config.LoadEnvVariable("PORT")
	if port == "" {
        port = "8080"
    }
    r.Logger.Fatal(r.Start(":" + port))
}
