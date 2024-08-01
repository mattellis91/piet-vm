package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {

}

func initModel() model {
	return model{}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
		
	}

	return m, nil
}

func (m model) View() string {
	s := "test view"
	return s
}