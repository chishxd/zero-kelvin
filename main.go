package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	Aura        int
	Temperature int
	Discipline  int
	Logs        []string
}

func (m model) View() string {
	s := "PROJECT ZERO KELVIN\n----------\n"

	for log := range m.Logs {
		s += fmt.Sprintf("Temperature: %d C \n Aura: %d C\n Discipline: %d C \n ----------\n Log: %s\n ", m.Temperature, m.Aura, m.Discipline, m.Logs[log])

	}

	s += "Press q to exit"
	return s
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:

		switch msg.String() {
		case "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) Init() tea.Cmd {
	return nil
}

func main() {
	initialModel := model{
		Aura:        0,
		Temperature: 65,
		Discipline:  10,
	}
	p := tea.NewProgram(initialModel)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Uh, oh... Seems like we have some error in the code: %v", err)
	}
}
