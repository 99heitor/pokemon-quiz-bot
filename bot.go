package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	pk "github.com/99heitor/pokemon-quiz-bot/pkmnquizbot"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var bot *tgbotapi.BotAPI

func handleCommand(update tgbotapi.Update) {

	if update.Message == nil {
		return
	}

	if update.Message.Chat.IsGroup() || update.Message.Chat.IsSuperGroup() {
		log.Printf("Request from chat: %s", update.Message.Chat.Title)
	} else {
		log.Printf("Request from user: %s", update.Message.Chat.UserName)
	}
	if bot.Debug {
		log.Printf("Update: %v", update.Message.Text)
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

func setupBot() {
	bot, _ = tgbotapi.NewBotAPI(os.Getenv("PKMN_QUIZ_BOT_TELEGRAM_TOKEN"))
}

func handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	for _, message := range sqsEvent.Records {
		var update tgbotapi.Update
		json.Unmarshal([]byte(message.Body), &update)
		fmt.Printf("The message %s for event source %s = %s \n", message.MessageId, message.EventSource, message.Body)
		handleCommand(update)
	}

	return nil
}

func init() {
	file, _ := os.Open("pokemon.csv")
	pk.AllPokemon, _ = csv.NewReader(file).ReadAll()
	pk.StoredAnswers = make(map[int64]pk.Pokemon)
	setupBot()
	rand.Seed(time.Now().UnixNano())
}

func main() {
	lambda.Start(handler)
}
