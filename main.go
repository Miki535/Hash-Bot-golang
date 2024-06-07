package main

import (
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/sha3"
	"os"

	"github.com/mymmrac/telego"

	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func main() {
	fmt.Println("Starting Bot")
	var token string
	fmt.Println("Enter Bot Token to authenticate")
	fmt.Scan(&token)
	// Захешований токен
	tokenHash := "fbc537c6657a0970090690f7fe85d2b014c5c4f4ce1d54df9a8112d42759cc1a"
	data1 := []byte(token)
	hash := sha3.Sum256(data1)
	// Переводимо hash з типу byte до типу string
	hashString1 := hex.EncodeToString(hash[:])
	// Перевірка на справжність токену
	if tokenHash == hashString1 {
		fmt.Println("Bot is authenticated")
	} else {
		fmt.Println("Bot is not authenticated")
	}

	bot, err := telego.NewBot(token, telego.WithDefaultDebugLogger())

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	updates, _ := bot.UpdatesViaLongPolling(nil)

	bh, _ := th.NewBotHandler(bot, updates)

	defer bh.Stop()
	defer bot.StopLongPolling()

	bh.Handle(func(bot *telego.Bot, update telego.Update) {
		chatID := tu.ID(update.Message.Chat.ID)

		message := tu.Message(
			chatID,
			"Введіть якусь інформацію щоб захешувати. Ми викоритовуєм SHA-3 це найновіший алгоритм хешування та дуже надійний",
		)
		bot.SendMessage(message)

	}, th.CommandEqual("start"))

	bh.Handle(func(bot *telego.Bot, update telego.Update) {
		chatID := tu.ID(update.Message.Chat.ID)
		newMessage := update.Message.Text
		hashedMessage := sha3.Sum256([]byte(newMessage))
		hashedString := hex.EncodeToString(hashedMessage[:])
		bot.SendMessage(tu.Message(chatID, hashedString))
	}, th.AnyMessageWithText())

	bh.Start()
}
