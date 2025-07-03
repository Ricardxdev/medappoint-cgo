package views

import (
	"ffi-test/global"
	"ffi-test/src/models"
	"ffi-test/src/utils"
	"unicode"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type UpdateModel struct {
	BaseModel
	focusIndex     int
	inputPatient   PatientInput
	searchInput    textinput.Model // This is used to search for the patient by ID
	errSearch      error
	err            error
	updatedPatient *models.Patient
	wasUpdated     bool // This will be set to true if the patient was updated successfully
}

func NewUpdateModel(parent tea.Model, parentBase BaseModel) UpdateModel {
	ti := textinput.New()
	ti.Placeholder = "Input ID to search"
	ti.Focus()
	ti.CharLimit = 8
	ti.Width = 10

	inputs := NewPatientInput()
	breadcrumb := append(parentBase.Breadcrumb, "Update Patient")
	return UpdateModel{
		searchInput: ti,
		BaseModel: BaseModel{
			Parent:     parent,
			Width:      parentBase.Width,
			Height:     parentBase.Height,
			Breadcrumb: breadcrumb,
		},
		focusIndex:     -1,
		inputPatient:   *inputs,
		errSearch:      nil,
		err:            nil,
		updatedPatient: nil,
		wasUpdated:     false,
	}
}

func (m UpdateModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m UpdateModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle search input first
	var cmds []tea.Cmd
	if m.focusIndex == -1 {
		var cmd tea.Cmd
		cmds = append(cmds, m.searchInput.Focus())
		m.searchInput.PromptStyle = focusedStyle
		m.searchInput.TextStyle = focusedStyle

		switch msg := msg.(type) {
		case tea.KeyMsg:
			key := msg.String()
			if len(key) == 1 && (unicode.IsLetter(rune(key[0])) || key == "_" || key == " ") {
				m.searchInput, cmd = m.searchInput.Update(msg)
			} else {
				switch msg.String() {
				case "enter":
					patient, err := global.PatientsService.GetPatient(m.searchInput.Value())
					if err != nil {
						m.updatedPatient = nil
						m.errSearch = error(err)
						cmds = append(cmds, textinput.Blink)
						return m, tea.Batch(cmds...)
					}

					m.inputPatient.FromPatient(patient.Patient)
					m.updatedPatient = &patient.Patient
					m.searchInput.Reset()

					return m, m.Init()
				case "backspace":
					if len(m.searchInput.Value()) != 0 {
						m.searchInput.SetValue(m.searchInput.Value()[:len(m.searchInput.Value())-1])
					}
				case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
					m.searchInput, cmd = m.searchInput.Update(msg)
				}
			}
		}

		cmds = append(cmds, cmd)
	}

	// Handle update form
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m.Parent, nil
		case "ctrl+c":
			return m, tea.Quit
			// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			if m.wasUpdated {
				m.wasUpdated = false // Reset after displaying
			}
			if m.errSearch != nil {
				m.errSearch = nil // Reset error after displaying
			}
			if m.updatedPatient != nil {
				s := msg.String()
				if err := m.inputPatient.Validate(); err != nil {
					m.err = err
				} else {
					m.err = nil
				}
				inputList := m.inputPatient.AsList()

				// Did the user press enter while the submit button was focused?
				// If so, exit.
				if s == "enter" {
					if m.err == nil && m.focusIndex == len(inputList) {
						patient := m.inputPatient.ToPatient()
						err := global.PatientsService.UpdatePatient(patient)
						if err != nil {
							m.err = err
						} else {
							m.wasUpdated = true
							m.updatedPatient = nil
							m.inputPatient = *NewPatientInput() // Reset inputs
							m.focusIndex = -1                   // Reset focus index
							return m, m.Init()
						}
					}
				}

				if m.focusIndex == -1 {
					// Remove focused state
					m.searchInput.Blur()
					m.searchInput.PromptStyle = noStyle
					m.searchInput.TextStyle = noStyle
				}

				// Cycle indexes
				if s == "up" || s == "shift+tab" {
					m.focusIndex--
					if m.focusIndex == 0 {
						m.focusIndex--
					}
				} else {
					m.focusIndex++
					if m.focusIndex == 0 {
						m.focusIndex++
					}
				}

				if m.focusIndex > len(inputList) {
					m.focusIndex = -1
				} else if m.focusIndex < -1 {
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
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *UpdateModel) updateInputs(msg tea.Msg) tea.Cmd {
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

func (m UpdateModel) View() string {
	// Create the breadcrumb view
	breadcrumbStr := utils.BreadcrumbView(m.Breadcrumb)
	s := breadcrumbStr + "\n\n" + utils.AlignW("Update Patient Form", m.Width) + "\n"
	s += utils.AlignW(m.searchInput.View(), m.Width) + "\n\n"

	if m.wasUpdated {
		s += utils.AlignW(valueStyle.Render("Patient updated successfully!"), m.Width)
		s += "\n\n"
	}

	if m.errSearch != nil {
		s += utils.AlignW(errorStyle.Render(m.errSearch.Error()), m.Width) + "\n\n"
	}

	if m.updatedPatient != nil {
		s += utils.AlignW(PatientAddFormView(m.inputPatient.AsList(), m.focusIndex), m.Width) + "\n"
		if m.err != nil {
			s += utils.AlignW(errorStyle.Render(m.err.Error()), m.Width) + "\n"
		}
	}

	return s
}
