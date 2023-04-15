package service

import (
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/Touchsung/money-note-line-api-go/config"
	"github.com/Touchsung/money-note-line-api-go/model"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)


func HandleYearlySummaryReport(event *linebot.Event, bot *linebot.Client) {
      // Retrieve user ID from event
    userID := event.Source.UserID

    // Connect to database
    db := config.ConnectDB()
    defer db.Close()

    // Retrieve income and expense data for each month in the year
    query := `
        SELECT
            EXTRACT(YEAR FROM created_at) AS year,
            EXTRACT(MONTH FROM created_at) AS month,
            SUM(CASE WHEN class = 'income' THEN amount ELSE 0 END) AS total_income,
            SUM(CASE WHEN class = 'expenses' THEN amount ELSE 0 END) AS total_expense
        FROM
            money_tracked
        WHERE
            user_id = $1 AND EXTRACT(YEAR FROM created_at) = $2 AND EXTRACT(MONTH FROM created_at) = $3
        GROUP BY
            year, month
        ORDER BY
            year, month
    `
    currentYear := time.Now().Year()
    var messageText string
    for month := 1; month <= 12; month++ {
        var totalIncome int
        var totalExpense int
        err := db.QueryRow(query, userID, currentYear, month).Scan(&currentYear, &month, &totalIncome, &totalExpense)
        if err != nil {
            if err == sql.ErrNoRows {
                continue // skip if no data for the month
            } else {
                log.Fatal(err)
            }
        }
        totalBalance := totalIncome - totalExpense
        messageText += fmt.Sprintf("📅 สรุปยอดรายรับ-รายจ่ายประจำเดือน %d/%d 📅\n", month, currentYear)
        messageText += fmt.Sprintf("💰 ยอดรายรับทั้งหมด: %d บาท\n", totalIncome)
        messageText += fmt.Sprintf("💸 ยอดรายจ่ายทั้งหมด: %d บาท\n", totalExpense)
        if totalBalance >= 0 {
            messageText += fmt.Sprintf("👍 คุณมีรายได้สุทธิ %d บาท\n", totalBalance)
        } else {
            messageText += fmt.Sprintf("👎 คุณมีรายจ่ายเกินรายได้ %d บาท\n", -totalBalance)
        }
        messageText += "\n"
    }
    if messageText == "" {
        messageText = "ไม่พบข้อมูลในปีนี้"
    } else {
        messageText += "โปรดจัดการการเงินอย่างมีสติ\U0001F609"
    }
    // Send summary report text message to user
    _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(messageText)).Do()
    if err != nil {
        log.Fatal(err)
    }
}

func HandleConfirmationMessage(event *linebot.Event, bot *linebot.Client, msgValues *model.MsgValues) {
    if msgValues.Text == "" || msgValues.Class == "" || msgValues.Category == "" || msgValues.Type == "" {
        bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("ไม่พบรายการที่จะบันทึก")).Do();
    } else{
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
	}
	*msgValues = model.MsgValues{}
}

func HandleCancelMessage(event *linebot.Event,bot *linebot.Client, msgValues *model.MsgValues)  {
    var messageReply = fmt.Sprintf("รายการ %s ถูกยกเลิกเรียบร้อย", msgValues.Text)
	 if msgValues.Text == "" || msgValues.Class == "" || msgValues.Category == "" || msgValues.Type == "" {
        messageReply = "ไม่พบรายการที่จะบันทึก"
    }
    _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(messageReply)).Do()

	if err != nil {
        log.Print(err)
    }
    *msgValues = model.MsgValues{}
}

func HandleMonthSummaryReport(event *linebot.Event, bot *linebot.Client) {
    // Retrieve user ID from event
    userID := event.Source.UserID
	
    // Connect to database
    db := config.ConnectDB()
    defer db.Close()

    // Retrieve income and expense data for the current month
    query := `
    SELECT
        COALESCE(SUM(CASE WHEN class = 'income' THEN amount ELSE 0 END), 0) AS total_income,
        COALESCE(SUM(CASE WHEN class = 'expenses' THEN amount ELSE 0 END), 0) AS total_expense
    FROM
        money_tracked
    WHERE
        user_id = $1 AND DATE_TRUNC('month', created_at) = DATE_TRUNC('month', CURRENT_DATE)
	`

    var totalIncome  int
    var totalExpense int

    err := db.QueryRow(query, userID).Scan(&totalIncome, &totalExpense)
    if err != nil {
        log.Fatal(err)
    }

	// Calculate the total balance
	totalBalance := totalIncome - totalExpense

	// Create summary text message with emoji
	// Generate summary report text message with emojis
	messageText := "📊 สรุปยอดรายรับ-รายจ่ายประจำเดือนนี้ 📊\n\n"
	messageText += fmt.Sprintf("💰 ยอดรายรับทั้งหมด: %d บาท\n", totalIncome)
	messageText += fmt.Sprintf("💸 ยอดรายจ่ายทั้งหมด: %d บาท\n", totalExpense)
	if totalBalance >= 0 {
		messageText += fmt.Sprintf("👍 คุณมีรายได้สุทธิ %d บาท\n", totalBalance)
	} else {
		messageText += fmt.Sprintf("👎 คุณมีรายจ่ายเกินรายได้ %d บาท\n", -totalBalance)
	}
	messageText += "\nโปรดจัดการการเงินอย่างมีสติ\U0001F609"

	// Send summary report text message to user
	_, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(messageText)).Do()
	if err != nil {
		log.Fatal(err)
	}
}