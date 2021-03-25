package pkmnquizbot

import (
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func init() {
	file, _ := os.Open("pokemon.csv")
	pokelist, _ := csv.NewReader(file).ReadAll()
	for _, p := range pokelist {
		pokemonId, _ := strconv.Atoi(p[0])
		allPokemon = append(allPokemon, Pokemon{id: pokemonId, name: p[1]})
	}

	rand.Seed(time.Now().UnixNano())
}

//AllPokemon will be initialized by the main function from the csv file
var allPokemon PokemonList
var wg sync.WaitGroup

//WhosThatPokemon sends a message with a shadow of a Pokemon image
func WhosThatPokemon(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	chatId := update.Message.Chat.ID
	wg.Add(1)
	go deleteShadowMessage(bot, getGameState(chatId))
	log.Printf("Rolling Pokémon for chat %d", chatId)
	randomPokemon := allPokemon.getRandom()
	log.Printf("Chose %s for chat %d", randomPokemon.name, chatId)
	msg := tgbotapi.NewMessage(chatId, fmt.Sprintf("[Who's that Pokémon?!](%s)", randomPokemon.getShadowUrl()))
	msg.ParseMode = "markdown"
	response, _ := bot.Send(msg)
	log.Printf("Shadow image for %s sent to %d", randomPokemon.name, chatId)
	saveGameState(chatId, response.MessageID, randomPokemon.id)
	wg.Wait()
}

//Its checks if the answer is the one stored for the current chat or is equal to "...", then reveals the answer.
func Its(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	chatId := update.Message.Chat.ID
	chatConfig := getGameState(chatId)
	if chatConfig.Id != 0 {
		pokemon := allPokemon.get(chatConfig.CurrentPokemon)
		if strings.EqualFold(update.Message.CommandArguments(), pokemon.name) || update.Message.CommandArguments() == "..." {
			msg := tgbotapi.NewMessage(chatId, fmt.Sprintf("[It's %s!](%s)", pokemon.name, pokemon.getPokemonUrl()))
			msg.ParseMode = "markdown"
			msg.ReplyToMessageID = update.Message.MessageID
			wg.Add(1)
			go deleteShadowMessage(bot, chatConfig)
			bot.Send(msg)
			wg.Wait()
		}
	}

}

func deleteShadowMessage(bot *tgbotapi.BotAPI, chatConfig ChatConfig) {
	defer wg.Done()
	if chatConfig.ShadowMessageId != 0 {
		bot.Send(&tgbotapi.DeleteMessageConfig{
			ChatID:    chatConfig.Id,
			MessageID: chatConfig.ShadowMessageId,
		})
	}
}
