package views

import (
	"ffi-test/src/models"
	"fmt"

	"github.com/charmbracelet/bubbles/table"
)

func PatientToRow(p *models.Patient) table.Row {
	return table.Row{
		p.ID,
		p.Name,
		fmt.Sprintf("%d", p.Age),
		p.Diagnosis,
		string(p.Gender),
		fmt.Sprintf("%d", func() int {
			if p.Disability {
				return 1
			} else {
				return 0
			}
		}()),
		p.DocSpecialty,
		p.AppointmentDate,
	}
}

func GetMaxWidth(patients []models.Patient, field string) int {
	maxWidth := 0
	for _, p := range patients {
		var width int
		switch field {
		case "Name":
			width = len(p.Name)
		case "Diagnosis":
			width = len(p.Diagnosis)
		case "Doc Specialty":
			width = len(p.DocSpecialty)
		case "Appointment Date":
			width = len(p.AppointmentDate)
		}

		if width > maxWidth {
			maxWidth = width
		}
	}
	return maxWidth
}
