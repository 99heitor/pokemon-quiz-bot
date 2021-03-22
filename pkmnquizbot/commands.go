package pkmnquizbot

import (
	"bytes"
	"fmt"
	png "image/png"
	"log"
	"math/rand"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

//AllPokemon will be initialized by the main function from the csv file
var AllPokemon PokemonList

//StoredAnswers holds the current Pokemon for any given chat
var StoredAnswers map[int64]Pokemon

//WhosThatPokemon sends a message with a shadow of a Pokemon image
func WhosThatPokemon(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	log.Printf("Rolling Pokémon for chat %d", update.Message.Chat.ID)
	r := rand.Intn(801)
	randomPokemon := AllPokemon.getPokemon(r + 1)
	StoredAnswers[update.Message.Chat.ID] = randomPokemon
	log.Printf("Chose %s for chat %d", randomPokemon.name, update.Message.Chat.ID)
	log.Printf("Generating shadow image for %s", randomPokemon.name)
	shadow := shadowImage{randomPokemon.img}
	shadowPNG := new(bytes.Buffer)
	png.Encode(shadowPNG, shadow)
	fileReader := tgbotapi.FileReader{Name: "Name", Reader: shadowPNG}
	log.Printf("Shadow image for %s generated", randomPokemon.name)
	log.Printf("Uploading shadow image for chat %d", update.Message.Chat.ID)
	msg := tgbotapi.NewPhoto(update.Message.Chat.ID, fileReader)
	msg.Caption = "Who's that Pokémon?"
	bot.Send(msg)
	log.Printf("Shadow image for %s sent to %d", randomPokemon.name, update.Message.Chat.ID)
}

//Its checks if the answer is the one stored for the current chat or is equal to "...", then reveals the answer.
func Its(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if answer, ok := StoredAnswers[update.Message.Chat.ID]; ok {
		if strings.EqualFold(update.Message.CommandArguments(), answer.name) || update.Message.CommandArguments() == "..." {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("[It's %s!](%s)", answer.name, answer.url))
			msg.ParseMode = "markdown"
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		}
	}

}
