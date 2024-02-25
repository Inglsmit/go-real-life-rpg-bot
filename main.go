package main

import (
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

var (
	gTgToken string
	gBot     *tgbotapi.BotAPI
	gChatId  int64
	err      error
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	if gTgToken = os.Getenv("TG_TOKEN"); gTgToken == "" {
		panic(fmt.Errorf(`failed to load token`))
	}

	gBot, err = tgbotapi.NewBotAPI(gTgToken)
	if err != nil {
		log.Panic(err)
	}

	gBot.Debug = true
}

func isStartMessage(update *tgbotapi.Update) bool {
	return update.Message != nil && update.Message.Text == "/start"
}

func main() {
	log.Printf("Authorized on account %s", gBot.Self.UserName)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = UPDATE_CONFIG_TIMEOUT

	updates := gBot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if isStartMessage(&update) {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID

			_, err := gBot.Send(msg)
			if err != nil {
				return
			}
		}
		//if update.Message != nil { // If we got a message
		//	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		//
		//	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		//	msg.ReplyToMessageID = update.Message.MessageID
		//
		//	gBot.Send(msg)
		//}
	}

}
