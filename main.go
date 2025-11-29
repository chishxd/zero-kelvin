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
	Will        int      //Might be action based and choice based
	Days        int      //The Amount of days passed
	Progress    int      //The Amount of time in a day that has passed
	Logs        []string //IDK If we needs TS. But might look cool
	State       int      //To track state of game like Playing, Lost or Victory
}

const (
	TicksPerDay = 5    //The amount of ticks that make a day
	MaxDays     = 90   //The amount of days you need to survive
	WinAura     = 5000 //The required amt of AURA needed to win the game after 5 days
	SafeTemp    = 30   //If model.Temperature goes below 30, the player freezes(game over)

	// States of game.. IDK why I am writing these dumb ahh comments tho
	StatePlaying  = 0
	StateGameOver = 1
	StateWon      = 2
)

type TickMsg time.Time

func viewGameOver(m model) string {
	//Styling for the borders of container
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(lipgloss.Color("#FF0000")).
		Padding(1, 3).
		Align(lipgloss.Center).
		Width(50)
	// Styling for big header
	headerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).
		Bold(true).
		Blink(true).
		PaddingBottom(1)
	// Smol Text
	textStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#CC6666"))

	reason := "WEAKNESS DETECTED"
	if m.Temperature > 60 {
		reason = "You Got Softened by Warmth of Luxury"
	} else if m.Temperature < SafeTemp {
		reason = "You caugh frostbite"
	} else if m.Will < 0 {
		reason = "You forgot your motives and installed League Of Legends"
	}

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		headerStyle.Render("YOU FAILED"),
		textStyle.Render(reason),
		textStyle.Render(fmt.Sprintf("FINAL AURA: %d", m.Aura)),
		"\n[ Press q to quit ]",
	)

	return boxStyle.Render(content) + "\n"
}

// The Render Logic for Victory screen
func viewWin(m model) string {
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color("#FFD700")).
		Padding(1, 3).
		Align(lipgloss.Center).
		Width(50)

	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FFFF")).
		Background(lipgloss.Color("#000000")).
		Bold(true).
		Padding(0, 1).
		MarginBottom(1)

	textStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF"))

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		headerStyle.Render("WINTER ARC CONQUERED!"),
		textStyle.Render("You are Great Alpha Male now"),
		"\n"+textStyle.Render(fmt.Sprintf("LEGENDARY AURA: %d", m.Aura)),
		"\n[ Press q to accept destiny ]\n",
	)

	return boxStyle.Render(content) + "\n"
}

// Suggest some better UI changes y'all
func viewDashboard(m model) string {

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
		"DAY: %d | TEMP: %d C  |  AURA: %d  |  WILL: %d",
		m.Days, m.Temperature, m.Aura, m.Will,
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
		"\n[c] Plunge | [g] Gym | [r] Read | [q] Quit",
	) + "\n"
}

func (m model) View() string {
	switch m.State {
	case StateWon:
		return viewWin(m)
	case StateGameOver:
		return viewGameOver(m)
	default:
		return viewDashboard(m)
	}

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
			m.Aura += 10
			m.Temperature -= 5
			m.Will -= 1
			m.Logs = append(m.Logs, "Did a cold plunge. Stay Hard")

			if len(m.Logs) > 5 {
				m.Logs = m.Logs[1:]
			}
		case "g":
			m.Aura += 100
			m.Temperature += 5
			m.Will -= 5
			m.Logs = append(m.Logs, "Lifted Heavy Weights. Feels PEAK")

			if len(m.Logs) > 5 {
				m.Logs = m.Logs[1:]
			}

		case "r":
			m.Aura += 3
			m.Will += 20
			if m.Will > 100 {
				m.Will = 100
			}
			m.Logs = append(m.Logs, "Knowledge Acquired. Focus restored")

			if len(m.Logs) > 5 {
				m.Logs = m.Logs[1:]
			}

		}
		return m, nil

	// So, A tick was received huh...
	case TickMsg:
		m.Temperature++
		m.Will--
		m.Progress++

		if m.Progress > TicksPerDay {
			m.Progress = 0
			m.Days += 1
			m.Logs = append(m.Logs, "Day "+fmt.Sprint(m.Days)+" begins.")

			if len(m.Logs) > 5 {
				m.Logs = m.Logs[1:]
			}
		}

		if m.Temperature > 60 || m.Temperature < SafeTemp || m.Will <= 0 {
			m.State = StateGameOver
		}

		if m.Days >= MaxDays && m.Aura >= WinAura {
			m.State = StateWon
		}

		if m.State == StatePlaying {
			return m, waitForTick()
		}

		return m, nil
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
		Will:        10,
		Days:        0,
		Progress:    0,
	}
	p := tea.NewProgram(initialModel)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Uh, oh... Seems like we have some error in the code: %v", err)
	}
}
