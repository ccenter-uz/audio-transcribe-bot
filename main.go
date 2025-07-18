package main

import (
	"log"
	"telegram-bot/audio"
	"telegram-bot/auth"
	"telegram-bot/config"
	"telegram-bot/handlers"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	bot, err := tgbotapi.NewBotAPI(cfg.TG_API.Token)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true
	log.Printf("Start bot: %s", bot.Self.UserName)

	auth.LoadTokens()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	// go handlers.YTLinks(bot, updates)

	for update := range updates {
		if update.CallbackQuery != nil {
			audio.HandleCallback(bot, update.CallbackQuery)
			continue
		}
		if update.Message == nil {
			continue
		}

		chatID := update.Message.Chat.ID
		msgText := update.Message.Text

		if handlers.HandleLogin(bot, chatID, msgText) {
			continue
		}

		if msgText == "/start" {
			handlers.HandleStart(bot, chatID)
			continue
		}

		if msgText == "/logout" {
			auth.Logout(chatID)
			bot.Send(tgbotapi.NewMessage(chatID, "✅ Siz tizimdan chiqdingiz."))
			continue
		}

		audio.HandleText(bot, update.Message.Chat.ID, update.Message.Text, update.Message.MessageID)
	}
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or failed to load")
	}
}
