package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/spf13/viper"
	witai "github.com/wit-ai/wit-go/v2"
)

//  Load ENV from .env file
func LoadEnvVariable(key string) string {
  viper.SetConfigFile(".env")

  err := viper.ReadInConfig()

  if err != nil {
    fmt.Printf("Error while reading config file %s", err)
  }

  value, ok := viper.Get(key).(string)

  if !ok {
    fmt.Printf("Invalid type assertion")
  }

  return value
}

// LINEConfig returns a new LINE SDK client
func LineClient() (*linebot.Client) {
  	CHANNEL_ACCESS_TOKEN := LoadEnvVariable("CHANNEL_ACCESS_TOKEN")
	  CHANNEL_SECRET := LoadEnvVariable("CHANNEL_SECRET")

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
  WIT_AI_TOKEN := LoadEnvVariable("WIT_AI_TOKEN")

  client := witai.NewClient(WIT_AI_TOKEN)

	resp, _ := client.Parse(&witai.MessageRequest{
		Query: msg,
	})

  return resp
}

// Connect to DB
func ConnectDB() *sql.DB{
   connStr := LoadEnvVariable("DB_URL")
    // Connect to database
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal(err)
    }
    return db
}

 



