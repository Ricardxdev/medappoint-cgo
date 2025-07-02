package src

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type FormModel struct {
	cursor   int
	inputs   []string
	focused  int
	complete bool
}

func NewFormModel() FormModel {
	return FormModel{
		cursor:  0,
		inputs:  make([]string, 7), // 7 fields for the patient form
		focused: 0,
	}
}

func (m FormModel) Init() tea.Cmd {
	return nil
}

func (m FormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if m.focused > 0 {
				m.focused--
			}
		case "down":
			if m.focused < len(m.inputs)-1 {
				m.focused++
			}
		case "enter":
			if m.focused == len(m.inputs)-1 {
				m.complete = true
			} else {
				m.focused++
			}
		case "backspace":
			if len(m.inputs[m.focused]) > 0 {
				m.inputs[m.focused] = m.inputs[m.focused][:len(m.inputs[m.focused])-1]
			}
		default:
			if len(msg.String()) == 1 {
				m.inputs[m.focused] += msg.String()
			}
		}
	}
	return m, nil
}

func (m FormModel) View() string {
	if m.complete {
		return "Form submitted! Press q to quit.\n"
	}

	var b strings.Builder
	b.WriteString("Add Patient Form:\n\n")

	fields := []string{
		"ID (CI):",
		"Name:",
		"Age:",
		"Diagnosis:",
		"Gender (M/F):",
		"Disability (0/1):",
		"Doctor Specialty:",
	}

	for i, field := range fields {
		cursor := " " // No cursor
		if m.focused == i {
			cursor = ">" // Cursor
		}
		b.WriteString(fmt.Sprintf("%s %s %s\n", cursor, field, m.inputs[i]))
	}

	b.WriteString("\nUse ↑/↓ to navigate, type to input, and Enter to submit.\n")
	return b.String()
}

func main() {
	p := tea.NewProgram(NewFormModel())
	if err := p.Start(); err != nil {
		fmt.Printf("Error starting program: %v\n", err)
	}
}
