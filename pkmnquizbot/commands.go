package pkmnquizbot

import (
	"bytes"
	"image"
	png "image/png"
	"math/rand"
	"strings"

	"gopkg.in/telegram-bot-api.v4"
)

var AllPokemon PokemonList
var StoredAnswers map[int64]Pokemon

func WhosThatPokemon(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	r := rand.Intn(801)
	randomPokemon := AllPokemon.getPokemon(r + 1)

	StoredAnswers[update.Message.Chat.ID] = AllPokemon.getPokemon(r + 1)
	decodedImage, _, _ := image.Decode(randomPokemon.img)
	shadow := shadowImage{decodedImage}
	shadowPNG := new(bytes.Buffer)
	png.Encode(shadowPNG, shadow)
	fileReader := tgbotapi.FileReader{Name: "Name", Reader: shadowPNG, Size: -1}

	msg := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, fileReader)
	msg.Caption = "Who's that Pokémon?"
	bot.Send(msg)
}

func Its(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if answer, ok := StoredAnswers[update.Message.Chat.ID]; ok {
		if strings.EqualFold(update.Message.CommandArguments(), answer.name) {
			fileReader := tgbotapi.FileReader{Name: "Name", Reader: answer.img, Size: -1}
			msg := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, fileReader)
			msg.Caption = "It's " + answer.name + "!"
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		}
	}

}
