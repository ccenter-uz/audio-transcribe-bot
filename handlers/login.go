package handlers

import (
	"telegram-bot/auth"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var userStates = make(map[int64]string)
var userLogins = make(map[int64]string)

func HandleLogin(bot *tgbotapi.BotAPI, chatID int64, msgText string) bool {
	switch userStates[chatID] {
	case "awaiting_login":
		userLogins[chatID] = msgText
		userStates[chatID] = "awaiting_password"
		bot.Send(tgbotapi.NewMessage(chatID, "🔑 Password:"))
		return true

	case "awaiting_password":
		login := userLogins[chatID]
		password := msgText

		info, err := auth.LoginUser(login, password)
		if err != nil {
			bot.Send(tgbotapi.NewMessage(chatID, "❌ Failed, again login."))
			delete(userStates, chatID)
			return true
		}

		auth.SetAuth(chatID, info)
		delete(userStates, chatID)
		delete(userLogins, chatID)
		startBtn := tgbotapi.NewMessage(chatID, "✅ Successfully logged in!")
		startBtn.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("▶️ Start Transcribe", "start_transcribe"),
			),
		)
		bot.Send(startBtn)
		return false
	}
	return false
}

func StartLogin(bot *tgbotapi.BotAPI, chatID int64) {
	userStates[chatID] = "awaiting_login"
	bot.Send(tgbotapi.NewMessage(chatID, "👤 Login:"))
}
