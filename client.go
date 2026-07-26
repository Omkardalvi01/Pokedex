package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	Search string
	Curr   Pokemon
}

func initialModel(P Pokemon) Model {
	return Model{
		Search: "",
		Curr:   P,
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {

}
