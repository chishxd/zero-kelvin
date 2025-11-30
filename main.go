package main

import (
	"fmt"
	"math/rand"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

// The Main State of the simulator. I learnt it from Bubble Tea's docs. They call it Elm Architecture
// React also seems to use similar architecture
type model struct {
	Aura         int      //CAN YOU FIX THE BROKEN, CAN YOU FEEL MY HEARRTTT
	Temperature  int      //The Body Temperature of the character, You lose if temp goes below 30 or above 75
	Will         int      //THE HAKI POWAAAA
	Days         int      //The Amount of days passed
	Progress     int      //The Amount of time in a day that has passed
	Logs         []string //IDK If we needs TS. But might look cool
	State        int      //To track state of game like Playing, Lost or Victory
	BusyTimer    int      //The amount of time player can't do any other stuff
	BusyTask     string   //Name of the task
	CurrentEvent GameEvent
	Width        int
	Height       int
}

type GameEvent struct {
	Prompt  string
	OptionA string
	OptionB string

	A_TempMod int
	A_WillMod int
	A_AuraMod int

	B_TempMod int
	B_WillMod int
	B_AuraMod int
}

const (
	TicksPerDay = 12   //The amount of ticks that make a day
	MaxDays     = 30   //The amount of days you need to survive
	WinAura     = 2000 //The required amt of AURA needed to win the game after 5 days
	SafeTemp    = 30   //If model.Temperature goes below 30, the player freezes(game over)

	// States of game.. IDK why I am writing these dumb ahh comments tho
	StatePlaying  = 0
	StateGameOver = 1
	StateWon      = 2
	StateMenu     = 3
	StateEvent    = 4
)

const Logo = `
███████╗███████╗██████╗  ██████╗     ██╗  ██╗███████╗██╗     ██╗   ██╗██╗███╗   ██╗
╚══███╔╝██╔════╝██╔══██╗██╔═══██╗    ██║ ██╔╝██╔════╝██║     ██║   ██║██║████╗  ██║
  ███╔╝ █████╗  ██████╔╝██║   ██║    █████╔╝ █████╗  ██║     ██║   ██║██║██╔██╗ ██║
 ███╔╝  ██╔══╝  ██╔══██╗██║   ██║    ██╔═██╗ ██╔══╝  ██║     ╚██╗ ██╔╝██║██║╚██╗██║
███████╗███████╗██║  ██║╚██████╔╝    ██║  ██╗███████╗███████╗ ╚████╔╝ ██║██║ ╚████║
╚══════╝╚══════╝╚═╝  ╚═╝ ╚═════╝     ╚═╝  ╚═╝╚══════╝╚══════╝  ╚═══╝  ╚═╝╚═╝  ╚═══╝
`

type TickMsg time.Time

func getRndEvent() GameEvent {
	events := []GameEvent{
		{
			Prompt:    "Your GF is shivering. She asks for your hoodie",
			OptionA:   "Give Hoodie (Warm her heart)",
			OptionB:   "Refuse (Cold builds character)",
			A_TempMod: -5, A_WillMod: -10, A_AuraMod: 50, // You get cold, lose will, but gain social aura?
			B_TempMod: 0, B_WillMod: 20, B_AuraMod: 100, // Stay warm, huge will, massive aura
		},
		{
			Prompt:    "Grandma cooked some cookies. Smells like love",
			OptionA:   "Eat One",
			OptionB:   "Don't eat(Reject sugar)",
			A_TempMod: 2, A_WillMod: -20, A_AuraMod: -100,
			B_TempMod: 0, B_WillMod: 50, B_AuraMod: 200,
		},
		{
			Prompt:    "The heater is fixed. The room is 22°C.",
			OptionA:   "Enjoy the warmth",
			OptionB:   "Open the windows (Sub-zero only)",
			A_TempMod: 10, A_WillMod: -30, A_AuraMod: -500,
			B_TempMod: -5, B_WillMod: 30, B_AuraMod: 300,
		},
		{
			Prompt:    "It is 3 AM. The Bed Looks soft and warm.",
			OptionA:   "Sleep on Bed (Recovery)",
			OptionB:   "Sleep on Floor",
			A_TempMod: 5, A_WillMod: 5, A_AuraMod: -50,
			B_TempMod: -2, B_WillMod: -5, B_AuraMod: 150,
		},
		{
			Prompt:    "It is snowing. You need to go to the Gym",
			OptionA:   "Wear Puffer Jacket",
			OptionB:   "Tank Top Only",
			A_TempMod: 5, A_WillMod: 0, A_AuraMod: -100,
			B_TempMod: -10, B_WillMod: -20, B_AuraMod: 600,
		},
		{
			Prompt:    "Your friends invite you to a Holiday Party",
			OptionA:   "GO with them",
			OptionB:   "Focus on yourself",
			A_TempMod: 3, A_WillMod: -10, A_AuraMod: -300,
			B_TempMod: 0, B_WillMod: 10, B_AuraMod: 400,
		},
		{
			Prompt:    "Your Ex texts you: 'I miss you'",
			OptionA:   "Reply",
			OptionB:   "Block & Lift",
			A_TempMod: 2, A_WillMod: -50, A_AuraMod: -1000,
			B_TempMod: 1, B_WillMod: 50, B_AuraMod: 500,
		},
		{
			Prompt:    "You have 1 hour of free time",
			OptionA:   "Scroll Tik Tok",
			OptionB:   "Stare at black wall",
			A_TempMod: 1, A_WillMod: -10, A_AuraMod: -200,
			B_TempMod: 0, B_WillMod: 50, B_AuraMod: 500,
		},
		{
			Prompt:    "You stubbed you toe on a dumbbell",
			OptionA:   "Scream",
			OptionB:   "Silence",
			A_TempMod: 1, A_WillMod: -10, A_AuraMod: -50,
			B_TempMod: 0, B_WillMod: 20, B_AuraMod: 500,
		},
	}

	return events[rand.Intn(len(events))]
}

func renderProgBar(current, max int, color string) string {
	const width = 20

	if max == 0 {
		max = 1
	}

	perct := float64(current) / float64(max)

	filled := min(int(perct*width), width)

	if filled < 0 {
		filled = 0
	}

	bar := ""

	for i := 0; i < filled; i++ {
		bar += "█"
	}
	for i := filled; i < width; i++ {
		bar += "░"
	}

	return lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Render(bar)

}

func viewMenu() string {
	logoStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#4D4DFF")).Bold(true)

	subStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFFF")).MarginBottom(2)

	instrStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#4D4DFF")).
		Padding(1, 2).
		Align(lipgloss.Left)

	instructions := `
GOAL: Survive 30 Days without losing

[c] Cold Plunge  ::  Cools you down.
[g] Gym          ::  Gains Aura (Heat Risk).
[r] Read         ::  Restores Willpower.

WARNING: Do not reach 0 Will, 30- Temp or 60+ Temp.
	`

	prompt := "\n[ PRESS ENTER TO BEGIN ]\n"

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		logoStyle.Render(Logo),
		subStyle.Render("THE ULTIMTE WINTER ARC SIMULATOR"),
		instrStyle.Render(instructions),
		prompt,
	)

	return content
}

func viewGameOver(m model) string {
	//Styling for the borders of container
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(lipgloss.Color("#FF0000")).
		Padding(1, 3).
		Align(lipgloss.Center).
		Width(80)
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
		"\n[ [q] quit | [enter] Retry ]",
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
		"\n[ [q] quit | [enter] Retry ]",
	)

	return boxStyle.Render(content) + "\n"
}

func formatStatDiff(temp, aura, will int) string {
	s := ""

	if temp != 0 {
		s += fmt.Sprintf("[%+d Temp]", temp)
	}
	if aura != 0 {
		s += fmt.Sprintf("[%+d Aura]", aura)
	}
	if will != 0 {
		s += fmt.Sprintf("[%+d Will]", will)
	}

	return s
}

func viewEvent(m model) string {
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color("#FF00FF")).
		Padding(1, 3).
		Align(lipgloss.Center)

	promptStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		MarginBottom(1)

	diffA := formatStatDiff(m.CurrentEvent.A_TempMod, m.CurrentEvent.A_AuraMod, m.CurrentEvent.A_WillMod)
	diffB := formatStatDiff(m.CurrentEvent.B_TempMod, m.CurrentEvent.B_AuraMod, m.CurrentEvent.B_WillMod)

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		promptStyle.Render("TEST OF WILL!"),
		"\n"+m.CurrentEvent.Prompt+"\n",
		"[ a ]"+m.CurrentEvent.OptionA+diffA,
		"[ b ]"+m.CurrentEvent.OptionB+diffB,
		"\n",
	)

	return boxStyle.Render(content)
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
		Padding(1).
		MarginBottom(1).
		Width(48)

	logStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#5C5C5C")).
		Width(48).
		Height(5).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#444444")).
		Padding(0, 1)

	busyStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#f1cb0cff"))

	// THE RENDERING LOGIC

	title := titleStyle.Render("PROJECT ZERO KELVIN\n")

	tempBar := renderProgBar(m.Temperature, 70, "#FF0000")
	willBar := renderProgBar(m.Will, 100, "#00FFFF")

	topRow := fmt.Sprintf("DAY: %d	|	Aura: %d ", m.Days, m.Aura)

	tempRow := fmt.Sprintf("Temp: %s %d", tempBar, m.Temperature)
	willRow := fmt.Sprintf("Will: %s %d", willBar, m.Will)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		topRow,
		"\n",
		tempRow,
		willRow,
	)

	stats := statsStyle.Render(content)

	logstr := "LOGS:\n"
	for _, l := range m.Logs {
		logstr += l + "\n"
	}
	history := logStyle.Render(logstr)

	var footer string
	if m.BusyTimer > 0 {
		footer = busyStyle.Render("[" + m.BusyTask + "]")
	} else {
		footer = "\n[c] Plunge | [g] Gym | [r] Read | [q] Quit"
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		stats,
		history,
		footer,
	) + "\n"
}

func (m model) View() string {

	var content string

	switch m.State {
	case StateMenu:
		content = viewMenu()
	case StateEvent:
		content = viewEvent(m)
	case StateWon:
		content = viewWin(m)
	case StateGameOver:
		content = viewGameOver(m)
	default:
		content = viewDashboard(m)
	}

	if m.Width == 0 {
		return content
	}

	return lipgloss.Place(
		m.Width, m.Height,
		lipgloss.Center, lipgloss.Center,
		content,
	)

}

// All Working logic is handled here

func (m model) performAction(taskName string, auraMod, tempMod, willMod, busyTime int, logMsg string) (tea.Model, tea.Cmd) {
	m.Aura += auraMod
	m.Temperature += tempMod
	m.Will += willMod

	if m.Will > 100 {
		m.Will = 100
	}

	m.BusyTask = taskName
	m.BusyTimer = busyTime

	m.addLog(logMsg)

	return m, nil
}

func (m model) handleKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if msg.String() == "q" || msg.String() == "ctrl+c" {
		return m, tea.Quit
	}

	if m.State == StateMenu {
		if msg.String() == "enter" {
			m.State = StatePlaying
			m.Aura = 0
			m.Temperature = 37
			m.Will = 10
			m.Days = 0
			m.Progress = 0
			m.Logs = []string{"Winter has arrived. Stay Strong"}
			return m, waitForTick()
		}
	}
	if m.State == StateGameOver || m.State == StateWon {
		if msg.String() == "enter" {
			m.State = StateMenu
		}
	}

	if m.State == StatePlaying {
		if m.BusyTimer > 0 {
			return m, nil
		}
		// A'ight, then what is it?
		switch msg.String() {

		case "c":
			return m.performAction("PLUNGED", 10, -5, -1, 0, "Did a cold plunge. Stay Hard")

		case "g":
			return m.performAction("LIFTING", 100, 5, -5, 1, "Lifted Heavy Weights. Feels PEAK")

		case "r":
			return m.performAction("READING", 1, 0, 20, 1, "Knowledge Acquired. Focus restored")
		}
	}

	if m.State == StateEvent {
		if msg.String() == "a" {
			m.Temperature += m.CurrentEvent.A_TempMod
			m.Aura += m.CurrentEvent.A_AuraMod
			m.Will += m.CurrentEvent.A_WillMod

			diffs := formatStatDiff(m.CurrentEvent.A_TempMod, m.CurrentEvent.A_AuraMod, m.CurrentEvent.A_WillMod)
			m.addLog("Chose: " + m.CurrentEvent.OptionA + diffs)

			m.State = StatePlaying
			return m, waitForTick()
		}
		if msg.String() == "b" {
			m.Temperature += m.CurrentEvent.B_TempMod
			m.Aura += m.CurrentEvent.B_AuraMod
			m.Will += m.CurrentEvent.B_WillMod

			diffs := formatStatDiff(m.CurrentEvent.B_TempMod, m.CurrentEvent.B_AuraMod, m.CurrentEvent.B_WillMod)
			m.addLog("Chose: " + m.CurrentEvent.OptionB + diffs)

			m.State = StatePlaying
			return m, waitForTick()
		}
	}

	return m, nil
}

func (m model) handleTick() (tea.Model, tea.Cmd) {
	m.Temperature++
	m.Will--
	m.Progress++

	// INCREMENT DAY
	if m.Progress > TicksPerDay {
		m.Progress = 0
		m.Days += 1
		m.addLog("Day " + fmt.Sprint(m.Days) + " begins.")

		if rand.Intn(100) < 10 {
			m.State = StateEvent
			m.CurrentEvent = getRndEvent()
			return m, nil
		}
	}

	// REDUCE BUSYTIMER ON EVERY TICK
	if m.BusyTimer > 0 {
		m.BusyTimer--
		if m.BusyTimer == 0 {
			m.BusyTask = ""
			m.addLog("Action Complete")
		}
	}

	// CHANGING STATES
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

// I guess this can be called the LOGIC part of the code. What to update and under which conditions should the update occur
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Is Msg a KeyPress?
	case tea.KeyMsg:
		return m.handleKeys(msg)

	// So, A tick was received huh...
	case TickMsg:
		return m.handleTick()

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil
	}

	return m, nil
}

func waitForTick() tea.Cmd {
	return tea.Tick(time.Second*2, func(t time.Time) tea.Msg { return TickMsg(t) })
}

func (m *model) addLog(msg string) {
	m.Logs = append(m.Logs, msg)
	if len(m.Logs) > 5 {
		m.Logs = m.Logs[1:]
	}
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
		State:       StateMenu,
	}
	p := tea.NewProgram(initialModel)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Uh, oh... Seems like we have some error in the code: %v", err)
	}
}
