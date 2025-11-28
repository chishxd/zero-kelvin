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

	s += fmt.Sprintf("Temperature: %d C \nAura: %d C\nDiscipline: %d C \n----------\n", m.Temperature, m.Aura, m.Discipline)

	s += "Log: \n"

	for log := range m.Logs {
		s += fmt.Sprintf("%s\n", m.Logs[log])

	}

	s += "Press c to Cold Plunge and q to exit"
	return s
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "c":
			m.Aura += 100
			m.Temperature -= 3
			m.Logs = append(m.Logs, "Did a cold plunge. Stay Hard")
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
		Temperature: 37,
		Discipline:  10,
	}
	p := tea.NewProgram(initialModel)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Uh, oh... Seems like we have some error in the code: %v", err)
	}
}
