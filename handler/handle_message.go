package handler

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"log"
	"net/http"

	"github.com/Touchsung/money-note-line-api-go/config"
	"github.com/labstack/echo/v4"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)


type MsgValues struct {
	Text	 string
    Category string
    Class    string
    Type     string
	imgUrl	 string
}

// Handler
func Hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

var msgValues MsgValues = MsgValues{}

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
				if message.Text == "ยืนยัน"{
					if msgValues.Text == "" || msgValues.Class == "" || msgValues.Category == "" || msgValues.Type == "" {
						bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("ไม่พบรายการที่จะบันทึก")).Do();
						msgValues = MsgValues{}	
					}

					pattern := `(\d+)`
					r := regexp.MustCompile(pattern)
					amountStr := r.FindString(msgValues.Text)
					numberInt, _ := strconv.Atoi(amountStr)
					
					db := config.ConnectDB()
					defer db.Close()

					// Check if the user exists
					var userExists bool
					err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE user_id = $1)", event.Source.UserID).Scan(&userExists)
					if err != nil {
						log.Fatal(err)
					}

					if !userExists {
						// Insert new user into the users table
						_, err = db.Exec("INSERT INTO users (user_id) VALUES ($1)", event.Source.UserID)
						if err != nil {
							log.Fatal(err)
						}
					}

					// Insert a new money tracked entry for the user
					_, err = db.Exec("INSERT INTO money_tracked (user_id, text, amount, class, type, category, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)",
						event.Source.UserID, msgValues.Text, numberInt, msgValues.Class, msgValues.Type, msgValues.Category, time.Now())
						
					if err != nil {
						log.Fatal(err)
					}

					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("เพิ่มลงฐานข้อมูลเรียบร้อย")).Do(); err != nil {
						log.Print(err)
					}

					msgValues = MsgValues{}

				} else if message.Text == "ยกเลิก"{
					if msgValues.Text == "" || msgValues.Class == "" || msgValues.Category == "" || msgValues.Type == "" {
						bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("ไม่พบรายการที่จะบันทึก")).Do(); 
					}
					messageReply := fmt.Sprintf("รายการ %s ถูกยกเลิกเรียบร้อย", msgValues.Text)
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(messageReply)).Do(); err != nil {
						log.Print(err)
					}
					msgValues = MsgValues{}
				} else{
					msgValues = extractMsgValues(message.Text)
					jsonTemplate := createLineTemplate(msgValues)
					contents, err := linebot.UnmarshalFlexMessageJSON([]byte(jsonTemplate))
					if err != nil {
						return err
					}
					_, err = bot.ReplyMessage(event.ReplyToken, linebot.NewFlexMessage("Tracked Money Template", contents),).Do()

					if err != nil {
						log.Println(err)
					}
				}
			}
		}
	}
	return nil
}

func extractMsgValues(text string) MsgValues {
	resp := config.ConnectWitAI(text)
	var msgValues MsgValues
	msgValues.Text = resp.Text
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

	jsonStringTH := `{
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
	"flexible": "ผันแปร",
	"salary": "เงินเดือน"
	}`

	data := make(map[string]string)
	err := json.Unmarshal([]byte(jsonStringTH), &data)
	if err != nil {
		panic(err)
	}

	msgValues.Category = data[msgValues.Category]
	msgValues.Type = data[msgValues.Type]
	msgValues.Class = data[msgValues.Class]

	if msgValues.Category == "" {
		msgValues.Category = "ไม่พบข้อมูล"
	}
	
	if msgValues.Type == "" {
		msgValues.Type = "ไม่พบข้อมูล"
	}

	if msgValues.Class == "" {
		msgValues.Class = "ไม่พบข้อมูล"
	}

	if msgValues.Class == "รายจ่าย"{
		msgValues.imgUrl = "https://media.istockphoto.com/id/1054309772/photo/abstract-red-gradient-color-background-christmas-valentine-wallpaper.jpg?b=1&s=170667a&w=0&k=20&c=d__OJwDP-aaeRAszoZa2AIxj0XFLTYUcgnmSl4ZY4wY="
	} else {
		msgValues.imgUrl = "https://images.unsplash.com/photo-1617957796155-72d8717ac882?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxzZWFyY2h8MXx8Z3JlZW4lMjBibHVycnklMjBiYWNrZ3JvdW5kfGVufDB8fDB8fA%3D%3D&w=1000&q=80"
	}

	return msgValues
}

func createLineTemplate(extractedMsg MsgValues) string {

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
			"url": "`+ extractedMsg.imgUrl +`"
		},
		"body": {
			"type": "box",
			"layout": "vertical",
			"contents": [
			{
				"type": "text",
				"text": "`+ extractedMsg.Class +`",
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
						"flex": 2
					},
					{
						"type": "text",
						"text": "`+ extractedMsg.Text +`",
						"wrap": true,
						"color": "#666666",
						"size": "sm",
						"flex": 5
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
						"flex": 2
					},
					{
						"type": "text",
						"text": "`+ extractedMsg.Type + `",
						"wrap": true,
						"color": "#666666",
						"size": "sm",
						"flex": 5
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
						"flex": 2
					},
					{
						"type": "text",
						"text": "`+ extractedMsg.Category +`",
						"wrap": true,
						"color": "#666666",
						"size": "sm",
						"flex": 5
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

	return jsonTemplate
}