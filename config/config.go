package config

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	witai "github.com/wit-ai/wit-go/v2"
)

// LINEConfig returns a new LINE SDK client
func LineClient() (*linebot.Client) {
  	CHANNEL_ACCESS_TOKEN := os.Getenv("CHANNEL_ACCESS_TOKEN")
	  CHANNEL_SECRET := os.Getenv("CHANNEL_SECRET")

    bot, err := linebot.New(
		CHANNEL_SECRET,
		CHANNEL_ACCESS_TOKEN,
	)
   
    if err != nil {
        return nil
    }

    return bot
}

// Connect Wit.ai
func ConnectWitAI(msg string) *witai.MessageResponse {
  WIT_AI_TOKEN := os.Getenv("WIT_AI_TOKEN")

  client := witai.NewClient(WIT_AI_TOKEN)

	resp, _ := client.Parse(&witai.MessageRequest{
		Query: msg,
	})

  return resp
}

// Connect to DB
func ConnectDB() *sql.DB{
   connStr := os.Getenv("DB_URL")
    // Connect to database
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal(err)
    }
    return db
}

 



