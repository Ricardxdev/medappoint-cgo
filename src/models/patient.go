package models

/*
#cgo CFLAGS: -I${SRCDIR}/../../csrc -DCGO_BUILD
#include "patient.h"
#include "errors.h"
#include "patient_metrics.h"
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"os"
	"unsafe"

	"github.com/sanity-io/litter"
)

type Gender byte

func (g Gender) String() string {
	switch g {
	case 'M':
		return "Male"
	case 'F':
		return "Female"
	default:
		return "Unknown"
	}
}

type Patient struct {
	ID              string
	Name            string
	Age             int
	Diagnosis       string
	Gender          byte
	Disability      bool
	DocSpecialty    string
	AppointmentDate string
}

type PatientIndex struct {
	CI       string // up to 8 characters
	Position uint   // index in the patients slice
}

type PatientService struct {
	patients       [C.MAX_PATIENTS]C.Patient
	count_patients C.size_t
	max_patients   C.size_t
	index          C.Index
	max_index      C.size_t
}

func NewPatientService() PatientService {
	return PatientService{
		patients:       [C.MAX_PATIENTS]C.Patient{},
		count_patients: 0,
		max_patients:   C.MAX_PATIENTS,
		index:          C.Index{},
		max_index:      C.MAX_INDEX,
	}
}

func NewPatient(p Patient) (C.Patient, error) {
	var c_patient C.Patient

	id := C.CString(p.ID)
	defer C.free(unsafe.Pointer(id))
	name := C.CString(p.Name)
	defer C.free(unsafe.Pointer(name))
	diagnosis := C.CString(p.Diagnosis)
	defer C.free(unsafe.Pointer(diagnosis))
	docSpecialty := C.CString(p.DocSpecialty)
	defer C.free(unsafe.Pointer(docSpecialty))

	appointmentDate := C.CString(p.AppointmentDate)
	defer C.free(unsafe.Pointer(appointmentDate))

	var disabilityInt C.int
	if p.Disability {
		disabilityInt = 1
	} else {
		disabilityInt = 0
	}

	errCode := C.NewPatient(&c_patient, id, name, C.int(p.Age), diagnosis, C.char(p.Gender), disabilityInt, docSpecialty, appointmentDate)
	if errCode != 0 {
		errMsg := C.GoString(C.ErrorDescription(errCode))
		return c_patient, fmt.Errorf("error creating patient: %s", errMsg)
	}

	return c_patient, nil
}

type PatientResponse struct {
	Patient Patient
	Index   uint
}

func (s *PatientService) GetPatient(ci string) (*PatientResponse, error) {
	cci := C.CString(ci)
	defer C.free(unsafe.Pointer(cci))

	var c_patient C.Patient
	var c_pIndex C.size_t
	errCode := C.GetPatient(&c_patient, &c_pIndex, &s.index[0], cci) // TODO: use s.index directly
	if errCode != 0 {
		if errCode == C.ERR_NOT_FOUND {
			return nil, fmt.Errorf("patient with CI %s not found", ci)
		}
		errMsg := C.GoString(C.ErrorDescription(errCode))
		return nil, fmt.Errorf("error getting patient: %s", errMsg)
	}

	return &PatientResponse{
		Patient: ParseCPatient(&c_patient),
		Index:   uint(c_pIndex),
	}, nil
}

func (s *PatientService) AddPatient(p Patient) error {
	c_patient, err := NewPatient(p)
	if err != nil {
		return err
	}

	errCode := C.AddPatient(&s.count_patients, &s.index, &c_patient)
	if errCode != 0 {
		errMsg := C.GoString(C.ErrorDescription(errCode))
		return fmt.Errorf("error adding patient: %s", errMsg)
	}

	s.patients = [C.MAX_PATIENTS]C.Patient{}
	s.count_patients = 0
	if err := s.LoadPatients(); err != nil {
		return fmt.Errorf("failed to reload patients after adding: %w", err)
	}

	return nil
}

func (s *PatientService) ListPatients() ([]Patient, error) {
	result := make([]Patient, s.count_patients)
	for i := 0; i < int(s.count_patients); i++ {
		result[i] = ParseCPatient(&s.patients[i])
		//fmt.Printf("Patient %d: %s (%c)\n", i, result[i].Name, result[i].Gender)
	}

	return result, nil
}

func (s *PatientService) UpdatePatient(p Patient) error {
	c_patient, err := NewPatient(p)
	if err != nil {
		return err
	}

	ci := C.CString(p.ID)
	defer C.free(unsafe.Pointer(ci))
	errCode := C.UpdatePatient(&s.index, ci, &c_patient)
	if errCode != 0 {
		errMsg := C.GoString(C.ErrorDescription(errCode))
		return fmt.Errorf("error updating patient: %s", errMsg)
	}

	s.patients = [C.MAX_PATIENTS]C.Patient{}
	s.count_patients = 0
	if err := s.LoadPatients(); err != nil {
		return fmt.Errorf("failed to reload patients after adding: %w", err)
	}

	return nil
}

func (s *PatientService) ScheduleAppointment(ci string, date string) error {
	cci := C.CString(ci)
	defer C.free(unsafe.Pointer(cci))
	cDate := C.CString(date)
	defer C.free(unsafe.Pointer(cDate))

	errCode := C.ScheduleAppointment(&s.patients[0], &s.index[0], cci, cDate)
	if errCode != 0 {
		errMsg := C.GoString(C.ErrorDescription(errCode))
		return fmt.Errorf("error scheduling appointment: %s", errMsg)
	}

	s.patients = [C.MAX_PATIENTS]C.Patient{}
	s.count_patients = 0
	if err := s.LoadPatients(); err != nil {
		return fmt.Errorf("failed to reload patients after adding: %w", err)
	}

	return nil
}

func (s *PatientService) DeletePatient(ci string) error {
	cci := C.CString(ci)
	defer C.free(unsafe.Pointer(cci))
	errCode := C.DeletePatient(&s.patients[0], &s.index, cci)
	if errCode != 0 {
		errMsg := C.GoString(C.ErrorDescription(errCode))
		return fmt.Errorf("error deleting patient: %s", errMsg)
	}

	s.patients = [C.MAX_PATIENTS]C.Patient{}
	s.count_patients = 0
	if err := s.LoadPatients(); err != nil {
		return fmt.Errorf("failed to reload patients after adding: %w", err)
	}

	return nil
}

func (s *PatientService) LoadPatients() error {
	errorCode := C.LoadPatients(&s.patients[0], &s.count_patients)
	if errorCode != 0 {
		return fmt.Errorf("Error loading patients: %d", errorCode)
	}
	return nil
}

func (s *PatientService) LoadIndex() error {
	errorCode := C.LoadIndex(&s.index)
	if errorCode != 0 {
		errMsg := C.GoString(C.ErrorDescription(errorCode))
		return fmt.Errorf("Error loading index: %s", errMsg)
	}
	return nil
}

func (s *PatientService) Load() error {
	if err := s.LoadPatients(); err != nil {
		return fmt.Errorf("failed to load patients: %w", err)
	}

	if err := s.LoadIndex(); err != nil {
		return fmt.Errorf("failed to load index: %w", err)
	}

	return nil
}

func (s *PatientService) Save() error {
	errorCode := C.SavePatients(&s.patients[0], s.count_patients)
	if errorCode != 0 {
		errMsg := C.GoString(C.ErrorDescription(errorCode))
		return fmt.Errorf("error saving patients: %s", errMsg)
	}

	// Save the index after saving patients
	errorCode = C.SaveIndex(&s.index)
	if errorCode != 0 {
		errMsg := C.GoString(C.ErrorDescription(errorCode))
		return fmt.Errorf("error saving index: %s", errMsg)
	}

	return nil
}

func (s *PatientService) CreateIndex() error {
	for i := C.size_t(0); i < s.count_patients; i++ {
		if s.patients[i].age == 0 {
			return fmt.Errorf("patient at index %d is uninitialized", i)
		}

		// Create a new index entry
		errCode := C.NewPatientIndex(&s.index, &s.patients[i].ci[0], i)
		if errCode != 0 {
			errMsg := C.GoString(C.ErrorDescription(errCode))
			return fmt.Errorf("error creating index for patient %d: %s", i, errMsg)
		}
	}

	fmt.Printf("Index created with %d entries.\n", s.count_patients)
	return nil
}

func (s *PatientService) SaveIndex() error {
	errorCode := C.SaveIndex(&s.index)
	if errorCode != 0 {
		errMsg := C.GoString(C.ErrorDescription(errorCode))
		return fmt.Errorf("error saving index: %s", errMsg)
	}
	return nil
}

func (s *PatientService) WriteIndexToFile() error {
	file, err := os.Create("index_log.txt")
	if err != nil {
		return fmt.Errorf("failed to create index file: %w", err)
	}
	defer file.Close()

	for i := 0; i < C.MAX_INDEX; i++ {
		if s.index[i].ci[0] == 0 {
			continue // Skip uninitialized index entries
		}
		_, err := fmt.Fprintf(file, "%d -> |%s|%d|\n", i, C.GoString(&s.index[i].ci[0]), s.index[i].position)
		if err != nil {
			return fmt.Errorf("failed to write index entry: %w", err)
		}
	}

	return nil
}

func (s *PatientService) Run() error {
	err := s.Load()
	if err != nil {
		return fmt.Errorf("failed to load data: %w", err)
	}

	fmt.Printf("Loaded %d patients.\n", s.count_patients)

	for _, cPatient := range s.patients {
		//patient := ParseCPatient(&cPatient)
		if cPatient.age == 0 {
			continue // Skip uninitialized patients
		}
		C.ShowPatient(&cPatient)
	}

	err = s.CreateIndex()
	if err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}

	fmt.Printf("\n\n\n\n")

	// // Add a patient for testing
	// testPatient := Patient{
	// 	ID:              "31668000",
	// 	Name:            "Test Patient",
	// 	Age:             30,
	// 	Diagnosis:       "Test Diagnosis",
	// 	Gender:          'M',
	// 	Disability:      false,
	// 	DocSpecialty:    "General",
	// 	AppointmentDate: "2023-10-01",
	// }

	// err = s.AddPatient(testPatient)
	// if err != nil {
	// 	return fmt.Errorf("failed to add test patient: %w", err)
	// }
	// fmt.Printf("Test patient added successfully.\n")

	// // Update a patient for testing
	// testPatient := Patient{
	// 	ID:              "31668000",
	// 	Name:            "Ricardo Martinez",
	// 	Age:             31,
	// 	Diagnosis:       "Updated Diagnosis",
	// 	Gender:          'M',
	// 	Disability:      false,
	// 	DocSpecialty:    "General",
	// 	AppointmentDate: "2023-10-01",
	// }

	// err = s.UpdatePatient(testPatient)
	// if err != nil {
	// 	return fmt.Errorf("failed to update test patient: %w", err)
	// }
	// fmt.Printf("Test patient updated successfully.\n")

	// // Delete a patient for testing
	// err = s.DeletePatient("31668000")
	// if err != nil {
	// 	return fmt.Errorf("failed to delete test patient: %w", err)
	// }
	// fmt.Printf("Test patient deleted successfully.\n")

	// // Test scheduling an appointment
	// testCI := "77889900"
	// testDate := "2025-12-24"
	// err = s.ScheduleAppointment(testCI, testDate)
	// if err != nil {
	// 	return fmt.Errorf("failed to schedule appointment for CI %s on date %s: %w", testCI, testDate, err)
	// }
	// fmt.Printf("Appointment scheduled successfully for CI %s on date %s.\n", testCI, testDate)

	// Save the patients after loading them
	errorCode := C.SavePatients(&s.patients[0], s.count_patients)
	if errorCode != 0 {
		errMsg := C.GoString(C.ErrorDescription(errorCode))
		return fmt.Errorf("error saving patients: %s", errMsg)
	}

	// Save the index after creating it
	err = s.SaveIndex()
	if err != nil {
		return fmt.Errorf("failed to save index: %w", err)
	}
	fmt.Printf("Index saved successfully.\n")

	for _, index := range s.index {
		if index.ci[0] == 0 {
			continue // Skip uninitialized index entries
		}
		fmt.Printf("Index CI: %s, Position: %d\n", C.GoString(&index.ci[0]), index.position)
		GetPatientResponse, err := s.GetPatient(C.GoString(&index.ci[0]))
		if err != nil {
			return fmt.Errorf("failed to get patient: %w", err)
		}
		litter.Dump(GetPatientResponse.Patient)
		fmt.Printf("\n")
	}

	return nil
}

// Wrapper for ErrorDescription function from errors.h
func ErrorDescription(code C.int) string {
	return C.GoString(C.ErrorDescription(code))
}

func ParseCPatient(cp *C.Patient) Patient {
	// fmt.Printf("Gender: %d\n", cp.gender)
	return Patient{
		ID:              C.GoString(&cp.ci[0]),
		Name:            C.GoString(&cp.name[0]),
		Age:             int(cp.age),
		Diagnosis:       C.GoString(&cp.diagnosis[0]),
		Gender:          byte(cp.gender),
		Disability:      cp.disability != 0,
		DocSpecialty:    C.GoString(&cp.doc_specialty[0]),
		AppointmentDate: C.GoString(&cp.appointment_date[0]),
	}
}
