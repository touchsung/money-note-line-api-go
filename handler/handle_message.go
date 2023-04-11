package handler

import (
	"encoding/json"

	"log"
	"net/http"

	"github.com/Touchsung/money-note-line-api-go/config"
	"github.com/labstack/echo/v4"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)


type MsgValues struct {
    Text     string
    Category string
    Class    string
    Type     string
}




// Handler
func Hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

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

	var msgValues MsgValues

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				resp := config.ConnectWitAI(message.Text)
				for key, trait := range resp.Traits {
					value := trait[0].Value
					switch key {
					case "category":
						msgValues.Category = value
					case "class":
						msgValues.Class = value
					case "type":
						msgValues.Type = value
					}	
				}
				jsonString := `{
				"food": "อาหาร",
				"expenses": "รายจ่าย",
				"income": "รายรับ",
				"extra": "เงินพิเศษ",
				"give": "การบริจาค",
				"entertainment": "ความบันเทิง",
				"travel": "ท่องเที่ยว",
				"credit_card": "ผ่อนชำระ",
				"cash": "ผ่อนชำระ",
				"health": "สุขภาพ",
				"stock": "หุ้น",
				"tax": "ภาษี",
				"travel_expenses": "ค่าเดินทาง",
				"invest": "การลงทุน",
				"saving": "การออม",
				"fixed": "คงที่",
				"flexible": "ผันแปร"
				}`

				data := make(map[string]string)
				err := json.Unmarshal([]byte(jsonString), &data)
				if err != nil {
					panic(err)
				}

				var imgUrl string
				if msgValues.Class == "expenses"{
					imgUrl = "https://media.istockphoto.com/id/1054309772/photo/abstract-red-gradient-color-background-christmas-valentine-wallpaper.jpg?b=1&s=170667a&w=0&k=20&c=d__OJwDP-aaeRAszoZa2AIxj0XFLTYUcgnmSl4ZY4wY="
				} else {
					imgUrl = "https://images.unsplash.com/photo-1617957796155-72d8717ac882?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxzZWFyY2h8MXx8Z3JlZW4lMjBibHVycnklMjBiYWNrZ3JvdW5kfGVufDB8fDB8fA%3D%3D&w=1000&q=80"
				}

				jsonTemplate := `{
					"type": "bubble",
					"hero": {
						"type": "image",
						"size": "full",
						"aspectRatio": "20:1",
						"aspectMode": "cover",
						"action": {
						"type": "uri",
						"uri": "http://linecorp.com/"
						},
						"url": "`+ imgUrl +`"
					},
					"body": {
						"type": "box",
						"layout": "vertical",
						"contents": [
						{
							"type": "text",
							"text": "`+ data[msgValues.Class] +`",
							"weight": "bold",
							"size": "xl"
						},
						{
							"type": "box",
							"layout": "vertical",
							"margin": "lg",
							"spacing": "sm",
							"contents": [
							{
								"type": "box",
								"layout": "baseline",
								"spacing": "sm",
								"contents": [
								{
									"type": "text",
									"text": "รายการ",
									"color": "#aaaaaa",
									"size": "sm",
									"flex": 1
								},
								{
									"type": "text",
									"text": "`+ message.Text +`",
									"wrap": true,
									"color": "#666666",
									"size": "sm",
									"flex": 4
								}
								]
							},
							{
								"type": "box",
								"layout": "baseline",
								"spacing": "sm",
								"contents": [
								{
									"type": "text",
									"text": "ประเภท",
									"color": "#aaaaaa",
									"size": "sm",
									"flex": 1
								},
								{
									"type": "text",
									"text": "`+ data[msgValues.Type] +`",
									"wrap": true,
									"color": "#666666",
									"size": "sm",
									"flex": 4
								}
								]
							},
							{
								"type": "box",
								"layout": "baseline",
								"spacing": "sm",
								"contents": [
								{
									"type": "text",
									"text": "หมวดหมู่",
									"color": "#aaaaaa",
									"size": "sm",
									"flex": 1
								},
								{
									"type": "text",
									"text": "`+ data[msgValues.Category] +`",
									"wrap": true,
									"color": "#666666",
									"size": "sm",
									"flex": 4
								}
								]
							}
							]
						}
						]
					},
					"footer": {
						"type": "box",
						"layout": "horizontal",
						"contents": [
						{
							"type": "button",
							"style": "link",
							"action": {
							"type": "message",
							"label": "ยืนยัน",
							"text": "ยืนยัน"
							},
							"height": "sm"
						},
						{
							"type": "button",
							"style": "link",
							"height": "sm",
							"action": {
							"type": "message",
							"label": "ยกเลิก",
							"text": "ยกเลิก"
							}
						}
						],
						"flex": 0
					}
				}`

				contents, err := linebot.UnmarshalFlexMessageJSON([]byte(jsonTemplate))
				if err != nil {
					return err
				}
				_, err = bot.ReplyMessage(event.ReplyToken, linebot.NewFlexMessage("Flex message alt text", contents),).Do()

				if err != nil {
					log.Println(err)
				}
			}
		}
	}
	return nil
}



