package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Touchsung/money-note-line-api-go/config"
	"github.com/labstack/echo/v4"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

// Handler
func Hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

// CallbackHandler handles LINE webhook callbacks
func LineCallbackHandler(c echo.Context) error {
	CHANNEL_ACCESS_TOKEN := config.LoadEnvVariable("CHANNEL_ACCESS_TOKEN")
	CHANNEL_SECRET := config.LoadEnvVariable("CHANNEL_SECRET")

	bot := config.LINEClient(CHANNEL_SECRET,CHANNEL_ACCESS_TOKEN)
	
	req := c.Request()
	res := c.Response()

	events, err := bot.ParseRequest(req)

	if err != nil {
		if err == linebot.ErrInvalidSignature {
			res.WriteHeader(400)
		} else {
			res.WriteHeader(500)
		}
		return nil
	}
	
	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text)).Do(); err != nil {
					log.Print(err)
				}
			case *linebot.StickerMessage:
				replyMessage := fmt.Sprintf(
					"sticker id is %s, stickerResourceType is %s", message.StickerID, message.StickerResourceType)
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
					log.Print(err)
				}
			}
		}
	}
	return nil
}
