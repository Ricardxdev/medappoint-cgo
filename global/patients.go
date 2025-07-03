package global

import "ffi-test/src/models"

var (
	// PatientsService is a global instance of PatientService
	PatientsService = models.NewPatientService()
)

func init() {
	err := PatientsService.Load()
	if err != nil {
		panic("Failed to load patients: " + err.Error())
	}
}
