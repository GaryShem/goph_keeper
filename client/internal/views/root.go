package views

import (
	tea "github.com/charmbracelet/bubbletea"
)

type RootModel struct {
	currentView tea.Model
	text        string
}

func NewRootModel() *RootModel {
	return &RootModel{
		currentView: NewLoginModel(),
	}
}

func (m *RootModel) Init() tea.Cmd {
	return nil
}

func (m *RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msgKey := msg.(type) {
	case tea.KeyMsg:
		switch msgKey.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}
	model, message := m.currentView.Update(msg)
	m.currentView = model
	return m, message
}

func (m *RootModel) View() string {
	return m.currentView.View() + "\n\nPress ctrl+c to exit" + m.text
}
