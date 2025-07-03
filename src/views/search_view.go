package views

import (
	"ffi-test/global"
	"ffi-test/src/models"
	"ffi-test/src/utils"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type SearchModel struct {
	BaseModel
	textInput textinput.Model
	err       error
	patient   *models.Patient
}

func NewSearchModel(parent tea.Model, parentBase BaseModel) SearchModel {
	ti := textinput.New()
	ti.Placeholder = "Input ID to search"
	ti.Focus()
	ti.CharLimit = 8
	ti.Width = 10

	breadcrumb := append(parentBase.Breadcrumb, "Search")
	return SearchModel{
		BaseModel: BaseModel{
			Parent:     parent,
			Width:      parentBase.Width,
			Height:     parentBase.Height,
			Breadcrumb: breadcrumb,
		},
		textInput: ti,
		err:       nil,
	}
}

func (m SearchModel) Init() tea.Cmd {
	return tea.Batch(
		tea.SetWindowTitle("Patient Management System"),
		tea.EnterAltScreen, // This enables alternate screen buffer
		textinput.Blink,
	)
}

func (m SearchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	//var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+c", "q":
			return m.Parent, nil
		case "enter":
			if m.textInput.Value() == "" {
				m.err = error(fmt.Errorf("input cannot be empty"))
			}

			patient, err := global.PatientsService.GetPatient(m.textInput.Value())
			if err != nil {
				m.err = error(err)
				return m, textinput.Blink
			}
			m.patient = &patient.Patient
			m.textInput.Reset()

			return m, m.Init()
		case "backspace":
			if len(m.textInput.Value()) != 0 {
				m.textInput.SetValue(m.textInput.Value()[:len(m.textInput.Value())-1])
			}
		case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
			m.textInput, _ = m.textInput.Update(msg)
		default:
			m.err = error(fmt.Errorf("invalid input: %s", msg))
		}
	}

	return m, textinput.Blink
}
func (m SearchModel) View() string {
	breadcrumbStr := utils.BreadcrumbView(m.Breadcrumb)
	s := breadcrumbStr + "\n\n" + utils.AlignW("Patient Search", m.Width) + "\n"
	s += utils.AlignW(m.textInput.View(), m.Width) + "\n\n"

	patientStr := ""
	if m.patient != nil {
		patientStr += PatientSummaryView(m.patient) + "\n"
	}
	err := ""
	if m.err != nil {
		err = errorStyle.Render(m.err.Error())
	}

	s += utils.AlignW(err, m.Width) + "\n"
	m.Height -= 7 // Adjust height to fit the screen

	s += utils.Center(patientStr, m.Width, m.Height)
	return s
}

func PatientSummaryView(p *models.Patient) string {
	// Find the max label width
	labels := []string{
		"CI:", "Name:", "Age:", "Gender:", "Diagnosis:", "Disability:", "Doc Speciality:", "Appointment-date:",
	}
	maxLabelWidth := 0
	for _, l := range labels {
		if len(l) > maxLabelWidth {
			maxLabelWidth = len(l)
		}
	}

	lines := []string{
		titleStyle.Render("Patient Summary"),
		"",
		fmt.Sprintf("%s %s", labelStyle.Render(fmt.Sprintf("%-*s", maxLabelWidth, "CI:")), valueStyle.Render(p.ID)),
		fmt.Sprintf("%s %s", labelStyle.Render(fmt.Sprintf("%-*s", maxLabelWidth, "Name:")), valueStyle.Render(p.Name)),
		fmt.Sprintf("%s %s", labelStyle.Render(fmt.Sprintf("%-*s", maxLabelWidth, "Age:")), valueStyle.Render(fmt.Sprintf("%d", p.Age))),
		fmt.Sprintf("%s %s", labelStyle.Render(fmt.Sprintf("%-*s", maxLabelWidth, "Gender:")), valueStyle.Render(string(p.Gender))),
		fmt.Sprintf("%s %s", labelStyle.Render(fmt.Sprintf("%-*s", maxLabelWidth, "Diagnosis:")), valueStyle.Render(p.Diagnosis)),
		fmt.Sprintf("%s %s", labelStyle.Render(fmt.Sprintf("%-*s", maxLabelWidth, "Disability:")), valueStyle.Render(fmt.Sprintf("%v", p.Disability))),
		fmt.Sprintf("%s %s", labelStyle.Render(fmt.Sprintf("%-*s", maxLabelWidth, "Doc Speciality:")), valueStyle.Render(p.DocSpecialty)),
		fmt.Sprintf("%s %s", labelStyle.Render(fmt.Sprintf("%-*s", maxLabelWidth, "Appointment-date:")), valueStyle.Render(p.AppointmentDate)),
	}
	content := strings.Join(lines, "\n")
	return boxStyle.Render(content)
}
