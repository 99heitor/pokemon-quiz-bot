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

//WhosThatPokemon sends a message with a shadow of a Pokemon image
func WhosThatPokemon(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	chatId := update.Message.Chat.ID
	deleteShadowMessage(bot, getGameState(chatId))
	log.Printf("Rolling Pokémon for chat %d", chatId)
	r := rand.Intn(801)
	randomPokemon := AllPokemon.getPokemon(r + 1)
	log.Printf("Chose %s for chat %d", randomPokemon.name, chatId)
	log.Printf("Generating shadow image for %s", randomPokemon.name)
	shadow := shadowImage{randomPokemon.getImage()}
	shadowPNG := new(bytes.Buffer)
	png.Encode(shadowPNG, shadow)
	fileReader := tgbotapi.FileReader{Name: "Name", Reader: shadowPNG}
	log.Printf("Shadow image for %s generated", randomPokemon.name)
	log.Printf("Uploading shadow image for chat %d", chatId)
	msg := tgbotapi.NewPhoto(chatId, fileReader)
	msg.Caption = "Who's that Pokémon?"
	response, _ := bot.Send(msg)
	log.Printf("Shadow image for %s sent to %d", randomPokemon.name, chatId)
	saveGameState(chatId, response.MessageID, randomPokemon.id)
}

//Its checks if the answer is the one stored for the current chat or is equal to "...", then reveals the answer.
func Its(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	chatId := update.Message.Chat.ID
	chatConfig := getGameState(chatId)
	if chatConfig.Id != 0 {
		pokemon := AllPokemon.getPokemon(chatConfig.CurrentPokemon)
		if strings.EqualFold(update.Message.CommandArguments(), pokemon.name) || update.Message.CommandArguments() == "..." {
			msg := tgbotapi.NewMessage(chatId, fmt.Sprintf("[It's %s!](%s)", pokemon.name, pokemon.getAssetUrl()))
			msg.ParseMode = "markdown"
			msg.ReplyToMessageID = update.Message.MessageID
			deleteShadowMessage(bot, chatConfig)
			bot.Send(msg)
		}
	}

}

func deleteShadowMessage(bot *tgbotapi.BotAPI, chatConfig ChatConfig) {
	if chatConfig.ShadowMessageId != 0 {
		bot.Send(&tgbotapi.DeleteMessageConfig{
			ChatID:    chatConfig.Id,
			MessageID: chatConfig.ShadowMessageId,
		})
	}
}
