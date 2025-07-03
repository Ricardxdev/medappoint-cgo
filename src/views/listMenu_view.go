package views

import (
	"ffi-test/global"
	"ffi-test/src/models"
	"ffi-test/src/utils"
	"fmt"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type listMenuModel struct {
	choices []string
	cursor  int
	BaseModel
}

// "List Patients by Appointment Date", "List Patients by Specialty", "List Patients Under Age"

type customTableModel struct {
	filterInput textinput.Model
	filterType  string
	err         error
	tableModel  tea.Model
	BaseModel
}

func NewCustomTableModel(filterType string, patients []models.Patient, parent tea.Model, parentBase BaseModel) customTableModel {
	ti := textinput.New()
	ti.Focus()
	ti.Width = 20
	ti.PromptStyle = focusedStyle
	ti.TextStyle = focusedStyle
	ti.Cursor.Style = cursorStyle

	return customTableModel{
		filterInput: ti,
		filterType:  filterType,
		tableModel:  NewTableModel(patients, parent, parentBase),
		BaseModel: BaseModel{
			Parent:     parent,
			Breadcrumb: append(parentBase.Breadcrumb, "Filter by "+filterType),
			Width:      parentBase.Width,
			Height:     parentBase.Height,
		},
		err: nil,
	}
}

func (m customTableModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m customTableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			return m.Parent, nil
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			if m.filterInput.Value() == "" {
				m.err = fmt.Errorf("input cannot be empty")
			} else {
				m.err = nil
				var patients []models.Patient
				var err error
				switch m.filterType {
				case "Appointment Date":
					// Validate date format (e.g., YYYY-MM-DD)
					if _, err := time.Parse(time.DateOnly, m.filterInput.Value()); err != nil {
						m.err = fmt.Errorf("invalid date format, expected YYYY-MM-DD")
						return m, nil
					}

					patients, err = global.PatientsService.ListPatientsByAppointmentDate(m.filterInput.Value())
				case "Specialty":
					patients, err = global.PatientsService.ListPatientsBySpecialty(m.filterInput.Value())
					if err != nil {
						m.err = err
						return m, nil
					}
				case "Under Age":
					// Validate age input (e.g., positive integer)
					var age int
					if age, err = strconv.Atoi(m.filterInput.Value()); err != nil {
						m.err = fmt.Errorf("invalid age format, expected positive integer")
						return m, nil
					}

					patients, err = global.PatientsService.ListPatientsUnderAge(age)
					if err != nil {
						m.err = err
						return m, nil
					}
				default:
					m.err = fmt.Errorf("unsupported filter type: %s", m.filterType)
					return m, nil
				}

				m.tableModel = NewTableModel(patients, m.Parent, m.BaseModel)

				return m, nil
			}
		}
	}

	m.filterInput, cmd = m.filterInput.Update(msg)
	if cmd != nil {
		return m, cmd
	}
	m.tableModel, cmd = m.tableModel.Update(msg)
	if cmd != nil {
		return m, cmd
	}
	return m, nil
}

func (m customTableModel) View() string {
	s := utils.BreadcrumbView(m.Breadcrumb) + "\n\n"               // 3 lines
	s += utils.AlignW("Filter by "+m.filterType, m.Width) + "\n\n" // 3 lines
	s += utils.AlignW(m.filterInput.View(), m.Width) + "\n"        // 2 lines
	if m.err != nil {
		s += utils.AlignW(errorStyle.Render(m.err.Error()), m.Width) + "\n" // 2 lines
	}
	s += "\n" // 1 line
	s += utils.Center(m.tableModel.(tableModel).Table.View(), m.Width, m.Height-11) + "\n"
	s += utils.AlignW("Row Count: "+strconv.Itoa(len(m.tableModel.(tableModel).Table.Rows())), m.Width) + "\n"
	return s
}

func NewListMenuModel(parent tea.Model, parentBase BaseModel) listMenuModel {
	breadCrumb := append(parentBase.Breadcrumb, "Select Patient List Menu")
	return listMenuModel{
		choices: []string{
			"List All Patients",
			"List Disabled Patients",
			"List Patients by Appointment Date",
			"List Patients by Specialty",
			"List Female Patients",
			"List Male Patients",
			"List Patients Under Age",
		},
		cursor: 0,
		BaseModel: BaseModel{
			Parent:     parent,
			Breadcrumb: breadCrumb,
			Width:      parentBase.Width,
			Height:     parentBase.Height,
		},
	}
}

func (m listMenuModel) Init() tea.Cmd {
	return tea.Batch(tea.SetWindowTitle("PatientPatient Management System - List Menu"), tea.EnterAltScreen)
}

func (m listMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		case "esc", "q":
			return m.Parent, nil
		case "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m listMenuModel) View() string {
	s := utils.BreadcrumbView(m.Breadcrumb) + "\n\n"
	menuStr := "Select Patient List Menu:\n\n"
	for i, choice := range m.choices {
		cursor := " " // no cursor
		row := fmt.Sprintf("%s %s", cursor, choice)
		if m.cursor == i {
			cursor = ">" // cursor
			row = global.SelectedStyle.Render(fmt.Sprintf("%s %s", cursor, choice))
		}
		menuStr += row + "\n"
	}

	s += utils.Center(menuStr, m.Width, m.Height-6)
	return s
}

func (m *listMenuModel) handleSelection() (tea.Model, tea.Cmd) {
	var err error
	var patientList []models.Patient
	switch m.cursor {
	case 0:
		patientList, err = global.PatientsService.ListPatients()
		if err != nil {
			return m, nil
		}
	case 1:
		patientList, err = global.PatientsService.ListDisabledPatients()
		if err != nil {
			return m, nil
		}
	case 2:
		t := NewCustomTableModel("Appointment Date", patientList, m, m.BaseModel)
		return t, t.Init()
	case 3:
		t := NewCustomTableModel("Specialty", patientList, m, m.BaseModel)
		return t, t.Init()

	case 4:
		patientList, err = global.PatientsService.ListFemalePatients()
		if err != nil {
			return m, nil
		}
	case 5:
		patientList, err = global.PatientsService.ListMalePatients()
		if err != nil {
			return m, nil
		}
	case 6:
		t := NewCustomTableModel("Under Age", patientList, m, m.BaseModel)
		return t, t.Init()
	default:
		return m, tea.Printf("Invalid selection: %d\n", m.cursor)
	}

	return NewTableModel(patientList, m, m.BaseModel), nil
}
