package pkmnquizbot

import (
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"net/http"
)

//Pokemon holds only the necessary info for the game to work
type Pokemon struct {
	id   int
	name string
}

//PokemonList is just the pokemon CSV
type PokemonList []Pokemon

const pokemonAssets = "https://assets.pokemon.com/assets/cms2/img/pokedex/full/%03d.png"
const shadowAssets = "https://static.heitor.dev/pkmn/shadow/%d.png"
const PokemonAmount = 898

func (p PokemonList) get(id int) Pokemon {
	return p[id-1]
}

func (p PokemonList) getRandom() Pokemon {
	return p.get(rand.Intn(PokemonAmount) + 1)
}

func (p Pokemon) getImage() image.Image {
	resp, _ := http.Get(p.getPokemonUrl())
	decodedImage, _, _ := image.Decode(resp.Body)
	return decodedImage
}

func (p Pokemon) getPokemonUrl() string {
	return fmt.Sprintf(pokemonAssets, p.id)
}

func (p Pokemon) getShadowUrl() string {
	return fmt.Sprintf(shadowAssets, p.id)
}

//shadowImage is the "shadow" version of an image: all non-alpha pixels are changed to black.
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
		return color.RGBA{255, 255, 255, 255}
	}
	return color.RGBA{0, 0, 0, 255}
}
