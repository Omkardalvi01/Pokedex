package main

import (
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

	tea "github.com/charmbracelet/bubbletea"
)

var (
	pokemon_search = "https://pokeapi.co/api/v2/pokemon/"
)

func RenderImg(img_body io.Reader) image.Image {
	img, _, err := image.Decode(img_body)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	return img

}

func searchPokemon(m Model) Model {
	client := http.Client{Timeout: time.Second * 1}
	resp, err := client.Get(pokemon_search + strings.ToLower(m.Search))
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

	front_img, err := client.Get(pokemon.Pic.Front)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	defer front_img.Body.Close()

	pokemon.Pic.Front_img = RenderImg(front_img.Body)

	back_img, err := client.Get(pokemon.Pic.Back)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	defer back_img.Body.Close()

	pokemon.Pic.Back_img = RenderImg(back_img.Body)

	m.Curr = pokemon
	m.Search = ""
	return m
}
func main() {

	p := tea.NewProgram(InitialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

}
