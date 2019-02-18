package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	pk "github.com/99heitor/pokemon-quiz-bot/pkmnquizbot"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_TOKEN"))
	if err != nil {
		log.Panic(err)
	}
	rand.Seed(time.Now().UnixNano())

	file, _ := os.Open("pokemon.csv")
	pk.AllPokemon, _ = csv.NewReader(file).ReadAll()
	pk.StoredAnswers = make(map[int64]pk.Pokemon)

	bot.Debug = false

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {

		if update.Message == nil {
			continue
		}

		command := update.Message.Command()
		switch {

		case strings.EqualFold(command, "who"):
			pk.WhosThatPokemon(bot, update)

		case strings.EqualFold(command, "its"):
			pk.Its(bot, update)

		case strings.EqualFold(command, "debug") && update.Message.Chat.ID == 36992723:
			rsp := fmt.Sprintf("Switching debug mode to %t", !bot.Debug)
			log.Printf(rsp)
			bot.Debug = !bot.Debug

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, rsp)
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		}
	}
}
