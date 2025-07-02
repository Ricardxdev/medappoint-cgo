package main

import (
	"ffi-test/global"
	"ffi-test/src"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// patientsService := src.NewPatientService()
	// err := patientsService.Run()
	// if err != nil {
	// 	fmt.Printf("Error running patient service: %v\n", err)
	// 	return
	// }

	// a menu entry for realize operations on patients
	// 1. Consult a specific patient by CI
	// 2. Add a new patient
	// 3. Modify the data of a specific patient
	// 4. Schedule a new appointment for an existing patient
	// 5. Remove a specific patient
	// 6. Lists, a menu entry for listing with a sub menu entry for each signature in src/metrics.go
	// 7. See indexes

	Run()
}

type Model struct {
	choices []string
	cursor  int
}

func NewModel() Model {
	return Model{
		choices: []string{
			"List Patients",
			"Add Patient",
			"Update Patient",
			"Delete Patient",
			"Schedule Appointment",
			"See Index",
			"Exit",
		},
		cursor: 0,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "down":
			m.cursor++
			if m.cursor >= len(m.choices) {
				m.cursor = 0
			}
		case "up":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(m.choices) - 1
			}
		case "enter":
			patients, err := global.PatientsService.ListPatients()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error listing patients: %v\n", err)
				return m, tea.Quit
			}
			listM := src.NewModel(patients)
			return listM, listM.Init()
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) View() string {
	s := "Patient Management Menu:\n\n"
	for i, choice := range m.choices {
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor
		}
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}
	s += "\nPress q to quit.\n"
	return s
}

func (m *Model) handleSelection() tea.Cmd {
	switch m.cursor {
	case 0:
		// Call function to list patients
		fmt.Println("Listing patients...")
	case 1:
		// Call function to add a patient
		fmt.Println("Adding a patient...")
	case 2:
		// Call function to update a patient
		fmt.Println("Updating a patient...")
	case 3:
		// Call function to delete a patient
		fmt.Println("Deleting a patient...")
	case 4:
		// Call function to schedule an appointment
		fmt.Println("Scheduling an appointment...")
	case 5:
		return m.Init()
	}
	return nil
}

func Run() {
	p := tea.NewProgram(NewModel())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting program: %v\n", err)
		os.Exit(1)
	}
}
