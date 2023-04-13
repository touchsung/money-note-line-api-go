package main

import (
	"github.com/Touchsung/money-note-line-api-go/router"
	_ "github.com/lib/pq"
)

func main() {
    r := router.Router()
    r.Logger.Fatal(r.Start(":80"))
}
