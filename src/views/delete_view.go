package views

import (
	"ffi-test/global"
	"ffi-test/src/models"
	"ffi-test/src/utils"
	"unicode"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type DeleteModel struct {
	BaseModel
	focusSearchBar  bool
	searchInput     textinput.Model
	errSearch       error
	errDelete       error
	patientToDelete *models.Patient
	deleted         bool // This will be set to true if the patient was deleted successfully
}

func NewDeleteModel(parent tea.Model, parentBase BaseModel) DeleteModel {
	ti := textinput.New()
	ti.Placeholder = "Input ID to search"
	ti.Focus()
	ti.CharLimit = 8
	ti.Width = 10

	return DeleteModel{
		searchInput: ti,
		BaseModel: BaseModel{
			Parent:     parent,
			Width:      parentBase.Width,
			Height:     parentBase.Height,
			Breadcrumb: append(parentBase.Breadcrumb, "Delete Patient"),
		},
		focusSearchBar:  true,
		errSearch:       nil,
		errDelete:       nil,
		patientToDelete: nil,
		deleted:         false,
	}
}

func (m DeleteModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m DeleteModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle search input first
	var cmds []tea.Cmd
	if m.focusSearchBar {
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
						m.patientToDelete = nil
						m.errSearch = error(err)
						cmds = append(cmds, textinput.Blink)
						return m, tea.Batch(cmds...)
					}

					m.patientToDelete = &patient.Patient
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
			if m.deleted {
				m.deleted = false // Reset after displaying
			}
			if m.errSearch != nil {
				m.errSearch = nil // Reset error after displaying
			}
			if m.patientToDelete != nil {
				s := msg.String()
				// Did the user press enter while the submit button was focused?
				// If so, exit.
				if s == "enter" {
					if !m.focusSearchBar {
						err := global.PatientsService.DeletePatient(m.patientToDelete.ID)
						if err != nil {
							m.errDelete = err
						} else {
							m.deleted = true
							m.patientToDelete = nil
							m.focusSearchBar = true
							return m, m.Init()
						}
					}
				}

				m.focusSearchBar = !m.focusSearchBar

				if !m.focusSearchBar {
					// Remove focused state
					m.searchInput.Blur()
					m.searchInput.PromptStyle = noStyle
					m.searchInput.TextStyle = noStyle
				} else {
					// Set focused state
					cmds = append(cmds, m.searchInput.Focus())
					m.searchInput.PromptStyle = focusedStyle
					m.searchInput.TextStyle = focusedStyle
				}

				return m, tea.Batch(cmds...)
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m DeleteModel) View() string {
	// Create the breadcrumb view
	breadcrumbStr := utils.BreadcrumbView(m.Breadcrumb)
	s := breadcrumbStr + "\n\n" + utils.AlignW("Delete Patient Form", m.Width) + "\n"
	s += utils.AlignW(m.searchInput.View(), m.Width) + "\n\n"

	if m.deleted {
		s += utils.AlignW(valueStyle.Render("Patient deleted successfully!"), m.Width)
		s += "\n\n"
	}

	if m.errSearch != nil {
		s += utils.AlignW(errorStyle.Render(m.errSearch.Error()), m.Width) + "\n\n"
	}

	if m.patientToDelete != nil {
		input := NewPatientInput()
		input.FromPatient(*m.patientToDelete)
		focusIndex := -1
		if !m.focusSearchBar {
			focusIndex = 8
		}
		s += utils.AlignW(PatientAddFormView(input.AsList(), focusIndex), m.Width) + "\n"
	}

	if m.errDelete != nil {
		s += utils.AlignW(errorStyle.Render(m.errDelete.Error()), m.Width) + "\n"
	}

	return s
}
