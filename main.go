package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strings"
	"time"
)

var (
	gTgToken string
	gBot     *tgbotapi.BotAPI
	gChatId  int64
	err      error

	gUsersInChat      Users
	gUsefulActivities = Activities{
		// Self-Development
		{"yoga", "Yoga (15 minutes)", 1},
		{"meditation", "Meditation (15 minutes)", 1},
		{"language", "Learning a foreign language (15 minutes)", 1},
		{"swimming", "Swimming (15 minutes)", 1},
		{"walk", "Walk (15 minutes)", 1},
		{"chores", "Chores (1 hour)", 1},

		// Work
		{"work_learning", "Studying work materials (15 minutes)", 1},
		{"portfolio_work", "Working on a portfolio project (15 minutes)", 1},
		{"resume_edit", "Resume editing (15 minutes)", 1},

		// Creativity
		{"creative", "Creative creation (15 minutes)", 1},
		{"reading", "Reading fiction literature (15 minutes)", 1},
	}

	gRewards = Activities{
		// Entertainment
		{"watch_series", "Watching a series (1 episode)", 10},
		{"watch_movie", "Watching a movie (1 item)", 30},
		{"social_nets", "Browsing social networks (30 minutes)", 10},

		// Food
		{"eat_sweets", "300 kcal of sweets", 60},
	}
)

type (
	User struct {
		id    int64
		name  string
		coins uint16
	}
	Users []*User

	Activity struct {
		code, name string
		coins      uint16
	}

	//List of the Activities. Which is slice of the pointers Activity
	Activities []*Activity
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

func isCallbackQuery(update *tgbotapi.Update) bool {
	return update.CallbackQuery != nil && update.CallbackQuery.Data != ""
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

func getKeyboardRow(buttonText, buttonCode string) []tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(buttonText, buttonCode))
}

func askToPrintIntro() {
	msg := tgbotapi.NewMessage(gChatId, "In the introductory messages, you can find the purpose of this bot and the rules of the game. What do you think?")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		getKeyboardRow(BUTTON_TEXT_PRINT_INTRO, BUTTON_CODE_PRINT_INTRO),
		getKeyboardRow(BUTTON_TEXT_SKIP_INTRO, BUTTON_CODE_SKIP_INTRO),
	)
	_, err := gBot.Send(msg)
	if err != nil {
		return
	}
}

func showMenu() {
	msg := tgbotapi.NewMessage(gChatId, "Choose option:")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		getKeyboardRow(BUTTON_TEXT_BALANCE, BUTTON_CODE_BALANCE),
		getKeyboardRow(BUTTON_TEXT_USEFUL_ACTIVITIES, BUTTON_CODE_USEFUL_ACTIVITIES),
		getKeyboardRow(BUTTON_TEXT_REWARDS, BUTTON_CODE_REWARDS),
	)

	_, err := gBot.Send(msg)
	if err != nil {
		return
	}
}

func showBalance(user *User) {
	msg := fmt.Sprintf("%s, your wallet is empty %s. Use useful activities to get more coins", user.name, EMOJI_DONT_KNOW)
	if coins := user.coins; coins > 0 {
		msg = fmt.Sprintf("%s, you have %d %s", user.name, EMOJI_COIN)
	}
	_, err := gBot.Send(tgbotapi.NewMessage(gChatId, msg))
	if err != nil {
		return
	}

	showMenu()
}

func isCallbackQueryMissing(update *tgbotapi.Update) bool {
	return update.CallbackQuery == nil || update.CallbackQuery.From == nil
}

func getUserFromUpdate(update *tgbotapi.Update) (user *User, found bool) {
	if isCallbackQueryMissing(update) {
		return
	}

	userId := update.CallbackQuery.From.ID
	for _, userInChat := range gUsersInChat {
		if userId == userInChat.id {
			return userInChat, true
		}
	}
	return
}

func storeUserFromUpdate(update *tgbotapi.Update) (user *User, found bool) {
	if isCallbackQueryMissing(update) {
		return
	}

	from := update.CallbackQuery.From
	// [?] what next string mean?
	user = &User{id: from.ID, name: strings.TrimSpace(from.FirstName + " " + from.LastName), coins: 0}
	gUsersInChat = append(gUsersInChat, user)
	return user, true
}

func showActivities(activities Activities, message string, isUseful bool) {
	activitiesButtonsRows := make([][]tgbotapi.InlineKeyboardButton, 0, len(activities)+1)
	for _, activity := range activities {
		activityDescription := ""
		if isUseful {
			activityDescription = fmt.Sprintf("+ %d %s: %s", activity.coins, EMOJI_COIN, activity.name)
		} else {
			activityDescription = fmt.Sprintf("- %d %s: %s", activity.coins, EMOJI_COIN, activity.name)
		}
		activitiesButtonsRows = append(activitiesButtonsRows, getKeyboardRow(activityDescription, activity.code))
	}
	activitiesButtonsRows = append(activitiesButtonsRows, getKeyboardRow(BUTTON_TEXT_PRINT_MENU, BUTTON_CODE_PRINT_MENU))

	msg := tgbotapi.NewMessage(gChatId, message)
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(activitiesButtonsRows...)
	_, err := gBot.Send(msg)
	if err != nil {
		return
	}
}

func showUsefulActivities() {
	showActivities(gUsefulActivities, "Track your actions or comeback to menu:", true)
}

func showRewards() {
	showActivities(gRewards, "Spend your coins or comeback to menu:", false)
}

func identifyActivity(activities Activities, choiceCode string) (activity *Activity, found bool) {
	for _, activity := range activities {
		if choiceCode == activity.code {
			return activity, true
		}
	}
	return nil, false
}

func processUsefulActivity(activity *Activity, user *User) {
	errorMsg := ""
	if activity.coins == 0 {
		errorMsg = fmt.Sprintf(`the activity "%s" doesn't have a specified cost`, activity.name)
	} else if user.coins+activity.coins > MAX_USER_COINS {
		errorMsg = fmt.Sprintf("you cannot have more than %d %s", MAX_USER_COINS, EMOJI_COIN)
	}

	resultMessage := ""
	if errorMsg != "" {
		resultMessage = fmt.Sprintf("%s, I'm sorry, but %s %s Your balance remains unchanged.", user.name, errorMsg, EMOJI_SAD)
	} else {
		user.coins += activity.coins
		resultMessage = fmt.Sprintf(`%s, the "%s" activity is completed! %d %s has been added to your account. Keep it up! %s%s Now you have %d %s`,
			user.name, activity.name, activity.coins, EMOJI_COIN, EMOJI_BICEPS, EMOJI_SUNGLASSES, user.coins, EMOJI_COIN)
	}
	sendStringMessage(resultMessage)
}

func processRewards(activity *Activity, user *User) {
	errorMsg := ""
	if activity.coins == 0 {
		errorMsg = fmt.Sprintf(`the reward "%s" doesn't have a specified cost`, activity.name)
	} else if user.coins < activity.coins {
		errorMsg = fmt.Sprintf(`you currently have %d %s. You cannot afford "%s" for %d %s`, user.coins, EMOJI_COIN, activity.name, activity.coins, EMOJI_COIN)
	}

	resultMessage := ""
	if errorMsg != "" {
		resultMessage = fmt.Sprintf("%s, I'm sorry, but %s %s Your balance remains unchanged, the reward is unavailable %s", user.name, errorMsg, EMOJI_SAD, EMOJI_DONT_KNOW)
	} else {
		user.coins -= activity.coins
		resultMessage = fmt.Sprintf(`%s, the reward "%s" has been paid for, get started! %d %s has been deducted from your account. Now you have %d %s`, user.name, activity.name, activity.coins, EMOJI_COIN, user.coins, EMOJI_COIN)
	}
	sendStringMessage(resultMessage)
}

func updateProcessing(update *tgbotapi.Update) {
	user, found := getUserFromUpdate(update)
	if !found {
		if user, found = storeUserFromUpdate(update); !found {
			_, err := gBot.Send(tgbotapi.NewMessage(gChatId, "Can't identify user"))
			if err != nil {
				return
			}
			return
		}
	}

	choiceCode := update.CallbackQuery.Data
	log.Printf("[%T] %s", time.Now(), choiceCode)

	switch choiceCode {
	case BUTTON_CODE_BALANCE:
		showBalance(user)
	case BUTTON_CODE_USEFUL_ACTIVITIES:
		showUsefulActivities()
	case BUTTON_CODE_REWARDS:
		showRewards()
	case BUTTON_CODE_PRINT_INTRO:
		printIntro()
		showMenu()
	case BUTTON_CODE_SKIP_INTRO:
		showMenu()
	case BUTTON_CODE_PRINT_MENU:
		showMenu()
	default:
		if usefulActivity, isFound := identifyActivity(gUsefulActivities, choiceCode); isFound {
			processUsefulActivity(usefulActivity, user)

			delay(2)
			showUsefulActivities()
			return
		}
		if reward, isFound := identifyActivity(gRewards, choiceCode); isFound {
			processRewards(reward, user)

			delay(2)
			showRewards()
			return
		}

		log.Printf(`[%T] !ERROR!: Unknown code "%s"`, time.Now(), choiceCode)
		msg := fmt.Sprintf("%s, sorry, I don't code '%s' %s. Ask administrator about support", user.name, choiceCode, EMOJI_SAD)
		_, err := gBot.Send(tgbotapi.NewMessage(gChatId, msg))
		if err != nil {
			return
		}
	}
}

func main() {
	log.Printf("Authorized on account %s", gBot.Self.UserName)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = UPDATE_CONFIG_TIMEOUT

	updates := gBot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if isCallbackQuery(&update) {
			updateProcessing(&update)
		} else if isStartMessage(&update) {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			gChatId = update.Message.Chat.ID
			askToPrintIntro()
			//printIntro()
		}
	}

}
