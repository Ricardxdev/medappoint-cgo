package views

import (
	"ffi-test/global"
	"ffi-test/src/models"
	"ffi-test/src/utils"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type AddModel struct {
	BaseModel
	focusIndex   int
	inputPatient PatientInput
	err          error
	addedPatient *models.Patient // This will hold the patient added after validation
}

func NewAddModel(parent tea.Model, parentBase BaseModel) AddModel {
	inputs := NewPatientInput()
	inputs.ID.Focus()
	inputs.ID.PromptStyle = focusedStyle
	inputs.ID.TextStyle = focusedStyle

	breadcrumb := append(parentBase.Breadcrumb, "Add Patient")
	return AddModel{
		BaseModel: BaseModel{
			Parent:     parent,
			Width:      parentBase.Width,
			Height:     parentBase.Height,
			Breadcrumb: breadcrumb,
		},
		focusIndex: 0,
		//cursorMode: cursor.CursorBlink,
		inputPatient: *inputs,
		err:          nil,
	}
}

func (m AddModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m AddModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m.Parent, nil
		case "ctrl+c":
			global.PatientsService.Save()
			return m, tea.Quit
			// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			if err := m.inputPatient.Validate(); err != nil {
				m.err = err
			} else {
				m.err = nil
			}
			inputList := m.inputPatient.AsList()
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" {
				if m.err == nil && m.focusIndex == len(inputList) {
					patient := m.inputPatient.ToPatient()
					err := global.PatientsService.AddPatient(patient)
					if err != nil {
						m.err = err
					} else {
						m.addedPatient = &patient
						m.inputPatient = *NewPatientInput() // Reset inputs
						m.inputPatient.ID.Focus()
						m.inputPatient.ID.PromptStyle = focusedStyle
						m.inputPatient.ID.TextStyle = focusedStyle
						m.focusIndex = 0 // Reset focus index
						return m, m.Init()
					}
				}
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(inputList) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(inputList)
			}

			cmds := make([]tea.Cmd, len(inputList))
			for i := 0; i <= len(inputList)-1; i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = inputList[i].Focus()
					inputList[i].PromptStyle = focusedStyle
					inputList[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				inputList[i].Blur()
				inputList[i].PromptStyle = noStyle
				inputList[i].TextStyle = noStyle
			}

			// update the patients inputs
			m.inputPatient.FromList(inputList)

			return m, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *AddModel) updateInputs(msg tea.Msg) tea.Cmd {
	inputList := m.inputPatient.AsList()
	cmds := make([]tea.Cmd, len(inputList))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range inputList {
		inputList[i], cmds[i] = inputList[i].Update(msg)
	}

	m.inputPatient.FromList(inputList)

	return tea.Batch(cmds...)
}

func (m AddModel) View() string {
	// Create the breadcrumb view
	breadcrumbStr := utils.BreadcrumbView(m.Breadcrumb)
	s := breadcrumbStr + "\n\n" + utils.AlignW("Add Patient Form", m.Width) + "\n"
	s += utils.AlignW(PatientAddFormView(m.inputPatient.AsList(), m.focusIndex), m.Width) + "\n"
	if m.err != nil {
		s += utils.AlignW(errorStyle.Render(m.err.Error()), m.Width) + "\n"
	}

	if m.addedPatient != nil {
		s += utils.AlignW("Patient added successfully!", m.Width) + "\n"
		m.addedPatient = nil // Reset after displaying the message
	}

	return s
}
