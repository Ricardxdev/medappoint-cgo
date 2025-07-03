package views

import (
	"ffi-test/src/models"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
)

func PatientAddFormView(input []textinput.Model, focusIndex int) string {
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
		"                    " + titleStyle.Render("Patient Form"),
		"",
	}
	// Create the lines for the patient form
	for i, input := range input {
		lines = append(lines, fmt.Sprintf("%s %s", labelStyle.Render(fmt.Sprintf("%-*s", maxLabelWidth, labels[i])), input.View()))
	}

	button := blurredButton
	if focusIndex == len(input) {
		button = focusedButton
	}
	lines = append(lines, "", "                    "+button)

	content := strings.Join(lines, "\n")
	return boxStyle.Render(content)
}

type PatientInput struct {
	ID              textinput.Model
	Name            textinput.Model
	Age             textinput.Model
	Diagnosis       textinput.Model
	Gender          textinput.Model
	Disability      textinput.Model
	DocSpecialty    textinput.Model
	AppointmentDate textinput.Model
}

func NewPatientInput() *PatientInput {
	i := &PatientInput{}
	i.ID = textinput.New()
	i.ID.Placeholder = "CI"
	i.ID.CharLimit = 8
	i.ID.Width = 30
	i.ID.Cursor.Style = cursorStyle

	i.Name = textinput.New()
	i.Name.Placeholder = "Name"
	i.Name.CharLimit = 50
	i.Name.Width = 30
	i.Name.Cursor.Style = cursorStyle

	i.Age = textinput.New()
	i.Age.Placeholder = "Age"
	i.Age.CharLimit = 3
	i.Age.Width = 30
	i.Age.Cursor.Style = cursorStyle

	i.Diagnosis = textinput.New()
	i.Diagnosis.Placeholder = "Diagnosis"
	i.Diagnosis.CharLimit = 100
	i.Diagnosis.Width = 30
	i.Diagnosis.Cursor.Style = cursorStyle

	i.Gender = textinput.New()
	i.Gender.Placeholder = "Gender"
	i.Gender.CharLimit = 10
	i.Gender.Width = 30
	i.Gender.Cursor.Style = cursorStyle

	i.Disability = textinput.New()
	i.Disability.Placeholder = "Disability"
	i.Disability.CharLimit = 100
	i.Disability.Width = 30
	i.Disability.Cursor.Style = cursorStyle

	i.DocSpecialty = textinput.New()
	i.DocSpecialty.Placeholder = "Doctor's Specialty"
	i.DocSpecialty.CharLimit = 100
	i.DocSpecialty.Width = 30
	i.DocSpecialty.Cursor.Style = cursorStyle

	i.AppointmentDate = textinput.New()
	i.AppointmentDate.Placeholder = "Appointment Date"
	i.AppointmentDate.CharLimit = 10
	i.AppointmentDate.Width = 30
	i.AppointmentDate.Cursor.Style = cursorStyle

	return i
}

func (i *PatientInput) AsList() []textinput.Model {
	return []textinput.Model{
		i.ID,
		i.Name,
		i.Age,
		i.Gender,
		i.Diagnosis,
		i.Disability,
		i.DocSpecialty,
		i.AppointmentDate,
	}
}

func (i *PatientInput) FromList(inputs []textinput.Model) {
	if len(inputs) != 8 {
		return
	}
	i.ID = inputs[0]
	i.Name = inputs[1]
	i.Age = inputs[2]
	i.Gender = inputs[3]
	i.Diagnosis = inputs[4]
	i.Disability = inputs[5]
	i.DocSpecialty = inputs[6]
	i.AppointmentDate = inputs[7]
}

func (i *PatientInput) Validate() error {
	// validate id
	if i.ID.Value() == "" {
		return fmt.Errorf("ID cannot be empty")
	}
	// Validate specific formats
	if len(i.ID.Value()) > 8 {
		return fmt.Errorf("ID must be up to 8 characters long")
	}
	// ID should be numeric
	if _, err := strconv.Atoi(i.ID.Value()); err != nil {
		return fmt.Errorf("ID must be a valid number")
	}

	// validate name
	if i.Name.Value() == "" {
		return fmt.Errorf("Name cannot be empty")
	}
	if len(i.Name.Value()) > 25 {
		return fmt.Errorf("Name must be up to 25 characters long")
	}

	// validate age
	if i.Age.Value() == "" {
		return fmt.Errorf("Age cannot be empty")
	}
	if v, err := strconv.Atoi(i.Age.Value()); err != nil {
		return fmt.Errorf("Age must be a valid number")
	} else if v <= 0 || v > 120 {
		return fmt.Errorf("Age must be between 1 and 120")
	}

	// validate gender
	if i.Gender.Value() == "" {
		return fmt.Errorf("Gender cannot be empty")
	}
	if len(i.Gender.Value()) != 1 {
		return fmt.Errorf("Gender must be a single character")
	}
	if i.Gender.Value() != "M" && i.Gender.Value() != "F" {
		return fmt.Errorf("Gender must be 'M' or 'F'")
	}

	// validate diagnosis
	if i.Diagnosis.Value() == "" {
		return fmt.Errorf("Diagnosis cannot be empty")
	}
	if len(i.Diagnosis.Value()) > 50 {
		return fmt.Errorf("Diagnosis must be up to 50 characters long")
	}

	// validate disability
	if i.Disability.Value() == "" {
		return fmt.Errorf("Disability cannot be empty")
	}
	yesValues := []string{"yes", "true", "1", "y", "si", "s"}
	noValues := []string{"no", "false", "0", "n"}
	var validDisability bool
	if i.Disability.Value() != "" {
		if strings.Contains(strings.Join(yesValues, ","), strings.ToLower(i.Disability.Value())) {
			validDisability = true
		} else if strings.Contains(strings.Join(noValues, ","), strings.ToLower(i.Disability.Value())) {
			validDisability = true
		}
	}
	if !validDisability {
		return fmt.Errorf("Disability must be 'yes' or 'no'")
	}

	// validate doctor's specialty
	if i.DocSpecialty.Value() == "" {
		return fmt.Errorf("Doctor's Specialty cannot be empty")
	}

	if len(i.DocSpecialty.Value()) > 50 {
		return fmt.Errorf("Doctor's Specialty must be up to 50 characters long")
	}

	// validate appointment date
	if i.AppointmentDate.Value() == "" {
		return fmt.Errorf("Appointment Date cannot be empty")
	}
	if _, err := time.Parse(time.DateOnly, i.AppointmentDate.Value()); err != nil {
		return fmt.Errorf("Appointment Date must be a valid date in YYYY-MM-DD format")
	}

	return nil
}

func (i *PatientInput) ToPatient() models.Patient {
	var age int
	if parsedValue, err := strconv.ParseInt(i.Age.Value(), 10, 64); err == nil {
		age = int(parsedValue)
	}

	yesValues := []string{"yes", "true", "1", "y", "si", "s"}
	var disability bool
	if i.Disability.Value() != "" {
		if strings.Contains(strings.Join(yesValues, ","), strings.ToLower(i.Disability.Value())) {
			disability = true
		}
	}

	// pad the ID to 8 digits if necessary
	id := i.ID.Value()
	if len(i.ID.Value()) < 8 {
		id = fmt.Sprintf("%0*s", 8, i.ID.Value())
	}
	return models.Patient{
		ID:              id,
		Name:            i.Name.Value(),
		Age:             age,
		Diagnosis:       i.Diagnosis.Value(),
		Gender:          i.Gender.Value()[0],
		Disability:      disability,
		DocSpecialty:    i.DocSpecialty.Value(),
		AppointmentDate: i.AppointmentDate.Value(),
	}
}

func (i *PatientInput) FromPatient(p models.Patient) {
	i.ID.SetValue(p.ID)
	i.Name.SetValue(p.Name)
	i.Age.SetValue(strconv.Itoa(p.Age))
	i.Diagnosis.SetValue(p.Diagnosis)
	i.Gender.SetValue(string(p.Gender))
	i.Disability.SetValue(fmt.Sprintf("%t", p.Disability))
	i.DocSpecialty.SetValue(p.DocSpecialty)
	i.AppointmentDate.SetValue(p.AppointmentDate)
}
