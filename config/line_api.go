package config

import (
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

// NewLINEConfig returns a new LINE SDK client
func LINEClient(secretKey string, accessToken string) (*linebot.Client) {
    bot, err := linebot.New(
		secretKey,
		accessToken,
	)
   

    if err != nil {
        return nil
    }

    return bot
}



