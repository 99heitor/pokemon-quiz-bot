package pkmnquizbot

import (
	"fmt"
	"image"
	"image/color"
	"io"
	"net/http"
)

type Pokemon struct {
	id   int
	name string
	img  io.Reader
}

type PokemonList [][]string

const pokemonAssets = "https://assets.pokemon.com/assets/cms2/img/pokedex/full/%03d.png"

func (p PokemonList) getPokemon(id int) Pokemon {
	resp, _ := http.Get(fmt.Sprintf(pokemonAssets, id))
	return Pokemon{id, p[id][30], resp.Body}
}

type shadowImage struct {
	originalImage image.Image
}

func (i shadowImage) ColorModel() color.Model {
	return i.originalImage.ColorModel()
}

func (i shadowImage) Bounds() image.Rectangle {
	return i.originalImage.Bounds()
}

func (i shadowImage) At(x, y int) color.Color {
	_, _, _, a := i.originalImage.At(x, y).RGBA()
	if a == 0 {
		return i.originalImage.At(x, y)
	}
	return color.RGBA{0, 0, 0, 255}
}
