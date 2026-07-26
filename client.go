package main

import (
	"fmt"
	"image"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/disintegration/imaging"
	"github.com/lucasb-eyer/go-colorful"
)

type Model struct {
	Search string
	Curr   Pokemon
}

func InitialModel() Model {
	return searchPokemon(Model{Search: "pikachu"})
}

func (m Model) Init() tea.Cmd {
	return nil // No initial side-effects needed
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			return searchPokemon(m), nil
		case "backspace":
			if len(m.Search) > 0 {
				m.Search = m.Search[:len(m.Search)-1]
			}
		default:
			if len(msg.String()) > 0 {
				m.Search += msg.String()
			}
		}
	}
	return m, nil
}

func (m Model) View() string {
	v := "Pokedex\n"
	v += fmt.Sprintf("> %s\n", m.Search)

	v += fmt.Sprintf("%s\n", m.Curr.Name)
	v += ToString(m.Curr.Pic.Front_img)
	// v += ToString(30, m.Curr.Pic.Back_img)

	v += "Abilities\n"
	for _, ability := range m.Curr.Abilities {
		v += fmt.Sprintf("  - %s\n", ability.Details.Name)
	}
	v += "Types\n"
	for _, t := range m.Curr.Types {
		v += fmt.Sprintf("  - %s\n", t.Details.Name)
	}
	return v
}
func ToString(img image.Image) string {
	width := 50

	img = imaging.Resize(img, width, 0, imaging.NearestNeighbor)
	b := img.Bounds()
	str := strings.Builder{}

	for y := b.Min.Y; y < b.Max.Y; y += 2 {
		for x := b.Min.X; x < b.Max.X; x++ {
			// Top pixel → foreground
			c1, _ := colorful.MakeColor(img.At(x, y))
			fgColor := lipgloss.Color(c1.Hex())

			// Bottom pixel → background
			bgHex := "#000000"
			if y+1 < b.Max.Y {
				c2, _ := colorful.MakeColor(img.At(x, y+1))
				bgHex = c2.Hex()
			}
			bgColor := lipgloss.Color(bgHex)

			str.WriteString(lipgloss.NewStyle().
				Foreground(fgColor).
				Background(bgColor).
				Render("▀"))
		}
		str.WriteString("\n")
	}

	return str.String()
}
