package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

type pokemon struct {
	id   int
	name string
}

var p [][]string

func getPokemon(id int) pokemon {
	return pokemon{id, p[id][30]}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	b, err := tb.NewBot(tb.Settings{
		Token:  "655858914:AAGyujNYdGtbfmQUwcCq6FI7H_lXmgfsNaE",
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	file, _ := os.Open("pokemon.csv")
	p, _ = csv.NewReader(file).ReadAll()
	storedAnswers := make(map[int64]pokemon)

	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle("/who", func(m *tb.Message) {
		r := rand.Intn(801)
		pic := &tb.Photo{
			File:    tb.FromURL(fmt.Sprintf("https://assets.pokemon.com/assets/cms2/img/pokedex/full/%03d.png", r+1)),
			Caption: "Who's that Pokémon?"}
		storedAnswers[m.Chat.ID] = getPokemon(r + 1)
		b.Send(m.Sender, pic)
	})

	b.Handle(tb.OnText, func(m *tb.Message) {
		if answer, ok := storedAnswers[m.Chat.ID]; ok {
			if m.Text == answer.name {
				b.Reply(m, "It's "+answer.name+"!")
			}
		}
	})

	b.Start()
}
