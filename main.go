package main

import (
	"ffi-test/global"
	"ffi-test/src/utils"
	"ffi-test/src/views"
	"fmt"
	"os"
	"strings"

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
	help    tea.Model
	views.BaseModel
}

func NewModel() Model {
	return Model{
		choices: []string{
			"List Patients",
			"Search Patient",
			"Add Patient",
			"Update Patient",
			"Delete Patient",
			"Exit",
		},
		cursor:    0,
		help:      views.NewHelpModel(),
		BaseModel: views.BaseModel{},
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(tea.SetWindowTitle("Patient Management System"), tea.EnterAltScreen)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Update the help model
	m.help, _ = m.help.Update(msg)
	// Update the breadcrumb for the main menu
	m.Breadcrumb = []string{"Main Menu"}

	// Handle window size messages
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
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
			return m.handleSelection()
		case "ctrl+c", "q":
			global.PatientsService.Save()
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) View() string {
	s := utils.BreadcrumbView(m.Breadcrumb) + "\n\n"
	menuStr := "Patient Management Menu:\n\n"
	for i, choice := range m.choices {
		cursor := " " // no cursor
		row := fmt.Sprintf("%s %s", cursor, choice)
		if m.cursor == i {
			cursor = ">" // cursor
			row = global.SelectedStyle.Render(fmt.Sprintf("%s %s", cursor, choice))
		}
		menuStr += row + "\n"
	}

	// Add help view
	helpStr := m.help.View()
	helpStrHeight := strings.Count(helpStr, "\n")

	s += utils.Center(menuStr, m.Width, m.Height-6-helpStrHeight)
	s += helpStr
	return s
}

func (m *Model) handleSelection() (tea.Model, tea.Cmd) {
	switch m.cursor {
	case 0:
		listM := views.NewListMenuModel(m, m.BaseModel)
		return listM, listM.Init()
	case 1:
		searchM := views.NewSearchModel(m, m.BaseModel)
		return searchM, searchM.Init()
	case 2:
		addM := views.NewAddModel(m, m.BaseModel)
		return addM, addM.Init()
	case 3:
		updateM := views.NewUpdateModel(m, m.BaseModel)
		return updateM, updateM.Init()
	case 4:
		deleteM := views.NewDeleteModel(m, m.BaseModel)
		return deleteM, deleteM.Init()
	case 5:
		global.PatientsService.Save()
		return m, tea.Quit
	}
	return m, nil
}

func Run() {
	p := tea.NewProgram(NewModel())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting program: %v\n", err)
		os.Exit(1)
	}

	// global.PatientsService.ListPatients()
}
