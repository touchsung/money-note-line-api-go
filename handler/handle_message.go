package handler

import (
	"strings"

	"net/http"

	"github.com/Touchsung/money-note-line-api-go/config"
	"github.com/Touchsung/money-note-line-api-go/model"
	"github.com/Touchsung/money-note-line-api-go/service"
	"github.com/labstack/echo/v4"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

// Handler
func Hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

var msgValues model.MsgValues = model.MsgValues{}

// CallbackHandler handles LINE webhook callbacks
func LineCallbackHandler(c echo.Context) error {
	bot := config.LineClient()
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
				if strings.HasPrefix(message.Text, "/") {
					service.HandleCommandMessage(event, bot, &msgValues)
				} else{
					msgValues = service.ExtractMsgValues(message.Text)
					service.HandleLineTemplate(event, bot,msgValues)

				}
			}
		}
	}
	return nil
