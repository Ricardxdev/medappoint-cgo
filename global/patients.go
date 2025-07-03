package global

import "ffi-test/src/models"

var (
	// PatientsService is a global instance of PatientService
	PatientsService = models.NewPatientService()
)

func init() {
	err := PatientsService.LoadPatients()
	if err != nil {
		panic("Failed to load patients: " + err.Error())
	}

	err = PatientsService.CreateIndex()
	if err != nil {
		panic("Failed to create index: " + err.Error())
	}
}
