package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/BourgeoisBear/rasterm"
)

var (
	pokemon_search = "https://pokeapi.co/api/v2/pokemon/"
)

func RenderImg(img_body io.Reader, idf string) {
	img, _, err := image.Decode(img_body)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	_ = rasterm.KittyWriteImage(os.Stdout, img, rasterm.KittyImgOpts{DstCols: 20, DstRows: 10})
	fmt.Println()

}
func main() {
	client := http.Client{Timeout: time.Second * 1}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter a pokemon name: ")
	scanner.Scan()
	pokemonName := strings.ToLower(scanner.Text())

	resp, err := client.Get(pokemon_search + pokemonName)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Pokemon not found")
		os.Exit(1)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	var pokemon Pokemon
	if err := json.Unmarshal(data, &pokemon); err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	front_img, err := client.Get(pokemon.Imgpaths.Front)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	defer front_img.Body.Close()

	back_img, err := client.Get(pokemon.Imgpaths.Back)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	defer back_img.Body.Close()

	RenderImg(front_img.Body, "front.png")
	RenderImg(back_img.Body, "back.png")

	fmt.Printf("Name: %s\n", pokemon.Name)

	fmt.Printf("Physical:\n")
	fmt.Printf("\t- Height: %d\n", pokemon.Height)
	fmt.Printf("\t- Weight: %d\n", pokemon.Weight)

	fmt.Println("Types")
	for _, Type := range pokemon.Types {
		fmt.Printf("\t- %s\n", Type.Details.Name)
	}

	fmt.Println("Abilities")
	for _, ability := range pokemon.Abilities {
		fmt.Printf("\t- %s\n", ability.Details.Name)
	}

	fmt.Println("Moves")
	for i, move := range pokemon.MoveSet {
		if i > 5 {
			break
		}
		fmt.Printf("\t- %s\n", move.Move.Name)
	}

	fmt.Println("Stat")
	for i, stat := range pokemon.Stats {
		if i > 5 {
			break
		}
		fmt.Printf("\t- %s:%d\n", stat.Stat.Name, stat.Base)
	}
}
