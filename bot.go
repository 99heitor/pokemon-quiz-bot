package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	pk "github.com/99heitor/pokemon-quiz-bot/pkmnquizbot"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/ssm"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var bot *tgbotapi.BotAPI
var mySession *session.Session
var localPtr *bool
var shadowGenerationPtr *bool

func init() {
	localPtr = flag.Bool("local", false, "Run the bot locally with long polling.")
	shadowGenerationPtr = flag.Bool("generate", false, "Generate shadow images, saves them to ./shadow")
	flag.Parse()
	if *shadowGenerationPtr {
		log.Printf("Generating shadow images...")
		pk.GenerateShadow()
		return
	}
	log.Printf("Initializing...")
	mySession = session.Must(session.NewSession())
	pk.DynamoClient = dynamodb.New(mySession)
	bot, _ = tgbotapi.NewBotAPI(getToken())
}

func handleUpdate(update tgbotapi.Update) {
	log.Printf("Handling telegram update %d.", update.UpdateID)
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
		log.Printf("Getting token from environment variable.")
		return token
	} else {
		log.Printf("Getting token from SSM.")

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
	log.Printf("Handling SQS event with %d record", len(sqsEvent.Records))
	for _, message := range sqsEvent.Records {
		var update tgbotapi.Update
		json.Unmarshal([]byte(message.Body), &update)
		handleUpdate(update)
	}

	return nil
}

// To run the bot locally use flag --local
func main() {

	if *shadowGenerationPtr {
		return
	} else if !*localPtr {
		lambda.Start(sqsHandler)
	} else {
		log.Printf("Running bot locally, initializing long polling channel.")

		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60

		bot.Request(&tgbotapi.DeleteWebhookConfig{
			DropPendingUpdates: false,
		})

		for update := range bot.GetUpdatesChan(u) {
			handleUpdate(update)
		}
	}
}
