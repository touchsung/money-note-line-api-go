package main

import (
	"github.com/Touchsung/money-note-line-api-go/router"
)

func main() {
	r := router.Router()
  	r.Logger.Fatal(r.Start(":80"))
}