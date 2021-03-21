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
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var bot *tgbotapi.BotAPI

func init() {
	file, _ := os.Open("pokemon.csv")
	pk.AllPokemon, _ = csv.NewReader(file).ReadAll()
	pk.StoredAnswers = make(map[int64]pk.Pokemon)
	bot, _ = tgbotapi.NewBotAPI(getToken())
	rand.Seed(time.Now().UnixNano())
}

func handleUpdate(update tgbotapi.Update) {
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

func getToken() string {
	token := os.Getenv("PKMN_TELEGRAM_TOKEN")
	if token != "" {
		return token
	} else {
		mySession := session.Must(session.NewSession())
		svc := ssm.New(mySession)
		param, err := svc.GetParameter(&ssm.GetParameterInput{
			Name:           aws.String("/telegram/token/pokemon-quiz-bot"),
			WithDecryption: aws.Bool(true),
		})

		if err != nil {
			panic(err)
		}
		return *param.Parameter.Value
	}
}

func sqsHandler(ctx context.Context, sqsEvent events.SQSEvent) error {
	for _, message := range sqsEvent.Records {
		var update tgbotapi.Update
		json.Unmarshal([]byte(message.Body), &update)
		fmt.Printf("The message %s for event source %s = %s \n", message.MessageId, message.EventSource, message.Body)
		handleUpdate(update)
	}

	return nil
}

func lambdaHandler() {
	lambda.Start(sqsHandler)
}

// used for running the bot locally, must define PKM_TELEGRAM_TOKEN
func main() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	bot.Request(&tgbotapi.DeleteWebhookConfig{
		DropPendingUpdates: false,
	})

	for update := range bot.GetUpdatesChan(u) {
		handleUpdate(update)
	}
}
