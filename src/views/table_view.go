package views

import (
	"ffi-test/src/models"
	"ffi-test/src/utils"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type tableModel struct {
	Table table.Model
	BaseModel
}

func NewTableModel(patients []models.Patient, parent tea.Model, parentBase BaseModel) tableModel {
	breadcrumb := append(parentBase.Breadcrumb, "Patient List")
	t := NewPatientTable(patients)
	return tableModel{
		BaseModel: BaseModel{
			Parent:     parent,
			Breadcrumb: breadcrumb,
			Width:      parentBase.Width,
			Height:     parentBase.Height,
		},
		Table: t,
	}
}

func (m tableModel) Init() tea.Cmd {
	return nil
}

func (m tableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.Table.Focused() {
				m.Table.Blur()
			} else {
				m.Table.Focus()
			}
		case "q", "ctrl+c":
			return m.Parent, nil
		case "enter":
			return m, tea.Batch(
				tea.Printf("Let's go to %s!", m.Table.SelectedRow()[1]),
			)
		}
	}
	m.Table, cmd = m.Table.Update(msg)
	return m, cmd
}

func (m tableModel) View() string {
	s := utils.BreadcrumbView(m.Breadcrumb) + "\n\n"
	s += "Patient List:\n\n"
	actH := m.Height - 6 // Previous lines
	tableStr := utils.Center(baseStyle.Render(m.Table.View()), m.Width, actH)

	s += tableStr + "\n"
	s += utils.AlignW("Row Count: "+strconv.Itoa(len(m.Table.Rows())), m.Width) + "\n"
	return s
}

func NewPatientTable(patients []models.Patient) table.Model {
	columns := []table.Column{
		{Title: "ID", Width: 8},
		{Title: "Name", Width: 25},
		{Title: "Age", Width: 3},
		{Title: "Diagnosis", Width: GetMaxWidth(patients, "Diagnosis")},
		{Title: "Gender", Width: 8},
		{Title: "Disability", Width: 10},
		{Title: "Doc Specialty", Width: GetMaxWidth(patients, "Doc Specialty")},
		{Title: "Appointment Date", Width: 10},
	}

	rows := make([]table.Row, 0, len(patients))
	for _, p := range patients {
		rows = append(rows, PatientToRow(&p))
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(15),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return t
}
