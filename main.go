package main

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

// The Main State of the simulator. I learnt it from Bubble Tea's docs. They call it Elm Architecture
// React also seems to use similar architecture
type model struct {
	Aura        int      //CAN YOU FIX THE BROKEN, CAN YOU FEEL MY HEARRTTT
	Temperature int      //The Body Temperature of the character, You lose if temp goes below 30 or above 75
	Discipline  int      //Might be action based and choice based
	Days        int      //The Amount of days passed
	Progress    int      //The Amount of time in a day that has passed
	Logs        []string //IDK If we needs TS. But might look cool
}

const (
	TicksPerDay = 5    //The amount of ticks that make a day
	MaxDays     = 90   //The amount of days you need to survive
	WinAura     = 5000 //The required amt of AURA needed to win the game after 5 days
	SafeTemp    = 30   //If model.Temperature goes below 30, the player freezes(game over)
)

type TickMsg time.Time

// Suggest some better UI changes y'all
func (m model) View() string {

	// THE LIPGLOSS
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#4D4DFF")).
		Padding(0, 1).
		MarginBottom(1).
		Width(50).
		Align(lipgloss.Center)

	statsStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00FFFF")).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#00FFFF")).
		Padding(0, 1).
		MarginBottom(1).
		Width(48).
		Align(lipgloss.Center)

	logStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#5C5C5C")).
		Width(48).
		Height(5).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#444444")).
		Padding(0, 1)

	// THE RENDERING LOGIC
	title := titleStyle.Render("PROJECT ZERO KELVIN\n")

	stats := statsStyle.Render(fmt.Sprintf(
		"DAY: %d | TEMP: %d C  |  AURA: %d  |  DISC: %d",
		m.Days, m.Temperature, m.Aura, m.Discipline,
	))

	logstr := "LOGS:\n"
	for _, l := range m.Logs {
		logstr += l + "\n"
	}
	history := logStyle.Render(logstr)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		stats,
		history,
		"\nPress 'c' to Plunge. 'q' to quit",
	)
}

// I guess this can be called the LOGIC part of the code. What to update and under which conditions should the update occur
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Is Msg a KeyPress?
	case tea.KeyMsg:
		// A'ight, then what is it?
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "c":
			m.Aura += 50
			m.Temperature -= 3
			m.Discipline += 1
			m.Logs = append(m.Logs, "Did a cold plunge. Stay Hard")

			if len(m.Logs) > 5 {
				m.Logs = m.Logs[1:]
			}
		}
		return m, nil

	case TickMsg:
		m.Temperature++
		m.Discipline--
		m.Progress++
		if m.Progress > TicksPerDay {
			m.Progress = 0
			m.Days += 1
			m.Logs = append(m.Logs, "Day "+fmt.Sprint(m.Days)+" begins.")

			if len(m.Logs) > 5 {
				m.Logs = m.Logs[1:]
			}

		}

		return m, waitForTick()
	}

	return m, nil
}

func waitForTick() tea.Cmd {
	return tea.Tick(time.Second*2, func(t time.Time) tea.Msg { return TickMsg(t) })
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
		Days:        0,
		Progress:    0,
	}
	p := tea.NewProgram(initialModel)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Uh, oh... Seems like we have some error in the code: %v", err)
	}
}
