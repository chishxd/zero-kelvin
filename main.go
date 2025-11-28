package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	Aura        int
	Temperature int
	Discipline  int
	Log         []string
}

func (m model) View() string {
	s := "PROJECT ZERO KELVIN"

	s += fmt.Sprintf("----------\n Temperature: %d C \n Aura: %d C\n Discipline: %d C \n", m.Temperature, m.Aura, m.Discipline)

	return s
}
