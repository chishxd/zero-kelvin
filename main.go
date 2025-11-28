package main

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// The Main State of the simulator. I learnt it from Bubble Tea's docs. They call it Elm Structure
// React also seems to use similar architecture
type model struct {
	Aura        int
	Temperature int
	Discipline  int
	Logs        []string
}

type TickMsg time.Time

// Suggest some better UI changes y'all
func (m model) View() string {
	s := "PROJECT ZERO KELVIN\n----------\n"

	s += fmt.Sprintf("Temperature: %d C \nAura: %d\nDiscipline: %d\n----------\n", m.Temperature, m.Aura, m.Discipline)

	s += "Log: \n"

	for log := range m.Logs {
		s += fmt.Sprintf("%s\n", m.Logs[log])

	}

	s += "Press c to Cold Plunge and q to exit"
	return s
}

// I guess this can be called the LOGIC part of the code. What to update and under which conditions should the update occur
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "c":
			m.Aura += 100
			m.Temperature -= 3
			m.Discipline += 1
			m.Logs = append(m.Logs, "Did a cold plunge. Stay Hard")
		}
		return m, nil

	case TickMsg:
		m.Temperature -= 1
		m.Discipline -= 1

		return m, waitForTick()
	}

	return m, nil
}

func waitForTick() tea.Cmd {
	return tea.Tick(time.Second * 2, func(t time.Time) tea.Msg { return TickMsg(t) })
}

// The Command to be executed as soon as program starts is written here
func (m model) Init() tea.Cmd {
	return waitForTick()
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
