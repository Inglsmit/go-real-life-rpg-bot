package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
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

func sendStringMessage(msg string) {
	_, err := gBot.Send(tgbotapi.NewMessage(gChatId, msg))
	if err != nil {
		return
	}
}
func delay(sec uint8) {
	time.Sleep(time.Second * time.Duration(sec))
}

func sendMessageWithDelay(delayInSec uint8, message string) {
	sendStringMessage(message)
	delay(delayInSec)
}

func printIntro() {
	sendMessageWithDelay(2, "Hello! "+EMOJI_SUNGLASSES)
	sendMessageWithDelay(7, "There are numerous beneficial actions that, by performing regularly, we improve the quality of our life. However, often it's more fun, easier, or tastier to do something harmful. Isn't that so?")
	sendMessageWithDelay(7, "With greater likelihood, we'll prefer to get lost in YouTube Shorts instead of an English lesson, buy M&M's instead of vegetables, or lie in bed instead of doing yoga.")
	sendMessageWithDelay(1, EMOJI_SAD)
	sendMessageWithDelay(10, "Everyone has played at least one game where you need to level up a character, making them stronger, smarter, or more beautiful. It's enjoyable because each action brings results. In real life, though, systematic actions over time start to become noticeable. Let's change that, shall we?")
	sendMessageWithDelay(1, EMOJI_SMILE)
	sendMessageWithDelay(14, `Before you are two tables: "Useful Activities" and "Rewards". The first table lists simple short activities, and for completing each of them, you'll earn the specified amount of coins. In the second table, you'll see a list of activities you can only do after paying for them with coins earned in the previous step.`)
	sendMessageWithDelay(1, EMOJI_COIN)
	sendMessageWithDelay(10, `For example, you spend half an hour doing yoga, for which you get 2 coins. After that, you have 2 hours of programming study, for which you get 8 coins. Now you can watch 1 episode of "Interns" and break even. It's that simple!`)
	sendMessageWithDelay(6, `Mark completed useful activities to not lose your coins. And don't forget to "purchase" the reward before actually doing it.`)

}

func main() {
	log.Printf("Authorized on account %s", gBot.Self.UserName)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = UPDATE_CONFIG_TIMEOUT

	updates := gBot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if isStartMessage(&update) {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			gChatId = update.Message.Chat.ID
			printIntro()
		}
	}

}
