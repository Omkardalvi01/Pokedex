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
	Normal bool
	Cursor string
}

func InitialModel() Model {
	return searchPokemon(Model{Search: "lucario", Cursor: "_", Normal: true})
}

func (m Model) Init() tea.Cmd {
	return nil // No initial side-effects needed
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "enter":
			return searchPokemon(m), nil
		case "backspace":
			if len(m.Search) > 0 {
				m.Search = m.Search[:len(m.Search)-1]
			}
		case "tab":
			m.Normal = !m.Normal
		default:
			if len(msg.String()) > 0 {
				m.Search += msg.String()
			}
		}
	}
	return m, nil
}

func (m Model) View() string {
	top := "Pokedex\n"
	img_width := 55
	Front_img := ""
	Back_img := ""
	TextWidth := 38
	if m.Normal {
		Front_img = ToString(img_width, m.Curr.Pic.Front_img)
		Back_img = ToString(img_width, m.Curr.Pic.Back_img)

	} else {
		Front_img = ToString(img_width, m.Curr.Pic.Shinyf_img)
		Back_img = ToString(img_width, m.Curr.Pic.Shinyb_img)
	}

	img_style := lipgloss.NewStyle().Width(img_width + 2).Border(lipgloss.NormalBorder())
	Front_render := img_style.Render(Front_img)
	Back_render := img_style.Render(Back_img)

	Leftsection := lipgloss.JoinHorizontal(lipgloss.Top, Front_render, Back_render)
	leftStyle := lipgloss.NewStyle().Width(2*img_width + 10).Border(lipgloss.NormalBorder())
	leftrender := leftStyle.Render(Leftsection)

	search := fmt.Sprintf("> %s%s\n", m.Search, m.Cursor)
	searchStyle := lipgloss.NewStyle().Width(TextWidth).Border(lipgloss.RoundedBorder()).Height(0)
	searchBlock := searchStyle.Render(search)

	name := strings.ToUpper(m.Curr.Name)
	nameStyle := lipgloss.NewStyle().Bold(true).Width(TextWidth).Italic(true).AlignHorizontal(lipgloss.Center)
	NameBlock := nameStyle.Render(name)

	ModeBlock := ""
	Shiny_mode := "Shiny"
	Normal_mode := "Normal"
	selected := lipgloss.NewStyle().Background(lipgloss.Color("#689d6a")).Border(lipgloss.HiddenBorder()).Foreground(lipgloss.Color("#f9f5d7"))
	nselected := lipgloss.NewStyle().Border(lipgloss.HiddenBorder())
	center_pos := lipgloss.NewStyle().Width(TextWidth).AlignHorizontal(lipgloss.Center)
	if m.Normal {
		shiny_render := nselected.Render(Shiny_mode)
		normal_render := selected.Render(Normal_mode)
		modes := lipgloss.JoinHorizontal(lipgloss.Top, normal_render, shiny_render)
		ModeBlock += center_pos.Render(modes)
	} else {
		shiny_render := selected.Render(Shiny_mode)
		normal_render := nselected.Render(Normal_mode)
		modes := lipgloss.JoinHorizontal(lipgloss.Top, normal_render, shiny_render)
		ModeBlock += center_pos.Render(modes)
	}

	TextBlock := ""
	TextBlock += "Abilities\n"
	for _, ability := range m.Curr.Abilities {
		TextBlock += fmt.Sprintf("  - %s\n", ability.Details.Name)
	}
	TextBlock += "Types\n"
	for _, t := range m.Curr.Types {
		TextBlock += fmt.Sprintf("  - %s\n", t.Details.Name)
	}

	TextBlock += "Stats\n"
	for _, stat := range m.Curr.Stats {
		TextBlock += fmt.Sprintf(" - %s:%d\n", stat.Stat.Name, stat.Base)
	}

	TextBlock += "Moves\n"
	for i, move := range m.Curr.MoveSet {
		if i > 5 {
			break
		}
		TextBlock += fmt.Sprintf(" - %s\n", move.Move.Name)
	}

	Rightsection := lipgloss.JoinVertical(lipgloss.Left, searchBlock, NameBlock, ModeBlock, TextBlock)
	rightStyle := lipgloss.NewStyle().Width(40).Border(lipgloss.NormalBorder())
	rightrender := rightStyle.Render(Rightsection)

	bottom := lipgloss.JoinHorizontal(lipgloss.Top, leftrender, rightrender)
	return lipgloss.JoinVertical(lipgloss.Center, top, bottom)
}

func ToString(width int, img image.Image) string {
	img = imaging.Resize(img, width, 0, imaging.NearestNeighbor)
	b := img.Bounds()
	imageWidth := b.Max.X
	h := b.Max.Y
	str := strings.Builder{}

	for heightCounter := 0; heightCounter < h; heightCounter += 2 {
		for x := imageWidth; x < width; x += 2 {
			str.WriteString(" ")
		}

		for x := 0; x < imageWidth; x++ {
			c1, _ := colorful.MakeColor(img.At(x, heightCounter))
			color1 := lipgloss.Color(c1.Hex())
			c2, _ := colorful.MakeColor(img.At(x, heightCounter+1))
			color2 := lipgloss.Color(c2.Hex())
			str.WriteString(lipgloss.NewStyle().Foreground(color1).
				Background(color2).Render("▀"))
		}

		str.WriteString("\n")
	}

	return str.String()
}
