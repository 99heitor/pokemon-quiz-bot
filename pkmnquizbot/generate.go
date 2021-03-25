package pkmnquizbot

import (
	"fmt"
	"image/png"
	"os"
	"path/filepath"
	"time"
)

func GenerateShadow() {
	newpath := filepath.Join(".", "shadows")
	os.MkdirAll(newpath, os.ModePerm)
	for i := 1; i <= PokemonAmount; i++ {
		shadow := shadowImage{allPokemon.get(i).getImage()}
		f, err := os.Create(fmt.Sprintf("shadows/%d.png", i))
		if err != nil {
			panic(err)
		}
		png.Encode(f, shadow)
		f.Close()
		time.Sleep(200 * time.Millisecond)
	}
}
