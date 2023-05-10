package service

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/Touchsung/money-note-line-api-go/config"
	"github.com/Touchsung/money-note-line-api-go/model"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func ExtractMsgValues(text string) model.MsgValues {
	resp := config.ConnectWitAI(text)
	var msgValues model.MsgValues
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

	return msgValues
}

func HandleLineTemplate(event *linebot.Event, bot *linebot.Client,extractedMsg model.MsgValues)  {
	jsonStringTH := `{
	"food": "อาหาร",
 "equipment": "อุปกรณ์",
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

	extractedMsg.Category = data[extractedMsg.Category]
	extractedMsg.Type = data[extractedMsg.Type]
	extractedMsg.Class = data[extractedMsg.Class]

	if extractedMsg.Category == "" {
		extractedMsg.Category = "ไม่พบข้อมูล"
	}
	
	if extractedMsg.Type == "" {
		extractedMsg.Type = "ไม่พบข้อมูล"
	}

	if extractedMsg.Class == "" {
		extractedMsg.Class = "ไม่พบข้อมูล"
	}

	if extractedMsg.Class == "รายจ่าย"{
		extractedMsg.ImgUrl = "https://media.istockphoto.com/id/1054309772/photo/abstract-red-gradient-color-background-christmas-valentine-wallpaper.jpg?b=1&s=170667a&w=0&k=20&c=d__OJwDP-aaeRAszoZa2AIxj0XFLTYUcgnmSl4ZY4wY="
	} else {
		extractedMsg.ImgUrl = "https://images.unsplash.com/photo-1617957796155-72d8717ac882?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxzZWFyY2h8MXx8Z3JlZW4lMjBibHVycnklMjBiYWNrZ3JvdW5kfGVufDB8fDB8fA%3D%3D&w=1000&q=80"
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
			"url": "`+ extractedMsg.ImgUrl +`"
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
				"text": "/ยืนยัน"
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
				"text": "/ยกเลิก"
				}
			}
			],
			"flex": 0
		}
	}`

	contents, err := linebot.UnmarshalFlexMessageJSON([]byte(jsonTemplate))

	if err != nil {
		log.Fatal(err)
	}
	_, err = bot.ReplyMessage(event.ReplyToken, linebot.NewFlexMessage("Tracked Money Template", contents),).Do()

	if err != nil {
		log.Println(err)
	}				
}

func HandleCommandMessage(event *linebot.Event, bot *linebot.Client, msgValues *model.MsgValues) {
    switch message := event.Message.(type) {
    case *linebot.TextMessage:
        // Check if message starts with command prefix "/"
        if strings.HasPrefix(message.Text, "/") {
            // Parse command
            command := strings.TrimPrefix(message.Text, "/")
            switch command {
            case "รายงานประจำเดือน":
                HandleMonthSummaryReport(event, bot)
			case "ยืนยัน":
                HandleConfirmationMessage(event, bot, msgValues)
            case "ยกเลิก":
                HandleCancelMessage(event, bot, msgValues)
			case "รายงานประจำปี":
				HandleYearlySummaryReport(event, bot)
            default:
                // Unknown command
                reply := linebot.NewTextMessage("ไม่รู้จักคำสั่ง \"" + command + "\"")
                _, err := bot.ReplyMessage(event.ReplyToken, reply).Do()
                if err != nil {
                    log.Fatal(err)
                }
            }
        }
    }
}