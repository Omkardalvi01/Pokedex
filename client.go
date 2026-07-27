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
		case "ctrl+c":
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

	Front_img := ""
	Back_img := ""
	if m.Normal {
		Front_img = ToString(55, m.Curr.Pic.Front_img)
		Back_img = ToString(55, m.Curr.Pic.Back_img)

	} else {
		Front_img = ToString(55, m.Curr.Pic.Shinyf_img)
		Back_img = ToString(55, m.Curr.Pic.Shinyb_img)
	}

	img_style := lipgloss.NewStyle().Width(60).Border(lipgloss.NormalBorder())
	Front_render := img_style.Render(Front_img)
	Back_render := img_style.Render(Back_img)

	Leftsection := lipgloss.JoinHorizontal(lipgloss.Top, Front_render, Back_render)
	leftStyle := lipgloss.NewStyle().Width(125).Border(lipgloss.NormalBorder())
	leftrender := leftStyle.Render(Leftsection)

	search := fmt.Sprintf("> %s%s\n", m.Search, m.Cursor)
	searchStyle := lipgloss.NewStyle().Width(38).Border(lipgloss.RoundedBorder()).Height(0)
	Rightsection := searchStyle.Render(search)

	name := fmt.Sprintf("%s", m.Curr.Name)
	nameStyle := lipgloss.NewStyle().AlignHorizontal(lipgloss.Center).Bold(true).Width(38)
	Rightsection += nameStyle.Render(name) + "\n"

	Shiny_mode := "Shiny"
	Normal_mode := "Normal"
	selected := lipgloss.NewStyle().Background(lipgloss.Color("#689d6a")).Border(lipgloss.HiddenBorder()).Foreground(lipgloss.Color("#f9f5d7"))
	nselected := lipgloss.NewStyle().Border(lipgloss.HiddenBorder())
	center_pos := lipgloss.NewStyle().AlignHorizontal(lipgloss.Center)
	if m.Normal {
		shiny_render := nselected.Render(Shiny_mode)
		normal_render := selected.Render(Normal_mode)
		modes := lipgloss.JoinHorizontal(lipgloss.Top, normal_render, shiny_render)
		Rightsection += center_pos.Render(modes) + "\n"
	} else {
		shiny_render := selected.Render(Shiny_mode)
		normal_render := nselected.Render(Normal_mode)
		modes := lipgloss.JoinHorizontal(lipgloss.Top, normal_render, shiny_render)
		Rightsection += center_pos.Render(modes) + "\n"
	}

	Rightsection += "Abilities\n"
	for _, ability := range m.Curr.Abilities {
		Rightsection += fmt.Sprintf("  - %s\n", ability.Details.Name)
	}
	Rightsection += "Types\n"
	for _, t := range m.Curr.Types {
		Rightsection += fmt.Sprintf("  - %s\n", t.Details.Name)
	}

	Rightsection += "Stats\n"
	for _, stat := range m.Curr.Stats {
		Rightsection += fmt.Sprintf(" - %s:%d\n", stat.Stat.Name, stat.Base)
	}

	Rightsection += "Moves\n"
	for i, move := range m.Curr.MoveSet {
		if i > 5 {
			break
		}
		Rightsection += fmt.Sprintf(" - %s\n", move.Move.Name)
	}

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
