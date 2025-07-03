package models

import (
	"fmt"
	"unsafe"
)

/*
#cgo CFLAGS: -I${SRCDIR}/../../csrc -DCGO_BUILD
#include "patient.h"
#include "errors.h"
#include "patient_metrics.h"
#include <stdlib.h>
*/
import "C"

// Wrapper for ListDisabledPatients function from patient_metrics.h
func (s *PatientService) ListDisabledPatients() ([]Patient, error) {
	var resultCount C.size_t
	var dest [C.MAX_PATIENTS]C.Patient
	errCode := C.ListDisabledPatients(&s.patients[0], s.max_patients, &dest[0], &resultCount)
	if errCode != 0 {
		return nil, fmt.Errorf("error listing disabled patients: %s", ErrorDescription(errCode))
	}

	result := make([]Patient, resultCount)
	for i := 0; i < int(resultCount); i++ {
		result[i] = ParseCPatient(&dest[i])
	}

	return result, nil
}

// Wrapper for ListPatientsByAppointmentDate function from patient_metrics.h
func (s *PatientService) ListPatientsByAppointmentDate(date string) ([]Patient, error) {
	cDate := C.CString(date)
	defer C.free(unsafe.Pointer(cDate))

	var resultCount C.size_t
	var dest [C.MAX_PATIENTS]C.Patient
	errCode := C.ListPatientsByAppointmentDate(&s.patients[0], s.max_patients, cDate, &dest[0], &resultCount)
	if errCode != 0 {
		return nil, fmt.Errorf("error listing patients by appointment date: %s", ErrorDescription(errCode))
	}

	result := make([]Patient, resultCount)
	for i := 0; i < int(resultCount); i++ {
		result[i] = ParseCPatient(&dest[i])
	}

	return result, nil
}

// Wrapper for ListPatientsBySpecialty function from patient_metrics.h
func (s *PatientService) ListPatientsBySpecialty(specialty string) ([]Patient, error) {
	cSpecialty := C.CString(specialty)
	defer C.free(unsafe.Pointer(cSpecialty))

	var resultCount C.size_t
	var dest [C.MAX_PATIENTS]C.Patient
	errCode := C.ListPatientsBySpecialty(&s.patients[0], s.max_patients, cSpecialty, &dest[0], &resultCount)
	if errCode != 0 {
		return nil, fmt.Errorf("error listing patients by specialty: %s", ErrorDescription(errCode))
	}

	result := make([]Patient, resultCount)
	for i := 0; i < int(resultCount); i++ {
		result[i] = ParseCPatient(&dest[i])
	}

	return result, nil
}

// Wrapper for ListFemalePatients function from patient_metrics.h
func (s *PatientService) ListFemalePatients() ([]Patient, error) {
	var resultCount C.size_t
	var dest [C.MAX_PATIENTS]C.Patient
	errCode := C.ListFemalePatients(&s.patients[0], s.max_patients, &dest[0], &resultCount)
	if errCode != 0 {
		return nil, fmt.Errorf("error listing female patients: %s", ErrorDescription(errCode))
	}

	result := make([]Patient, resultCount)
	for i := 0; i < int(resultCount); i++ {
		result[i] = ParseCPatient(&dest[i])
	}

	return result, nil
}

// Wrapper for ListMalePatients function from patient_metrics.h
func (s *PatientService) ListMalePatients() ([]Patient, error) {
	var resultCount C.size_t
	var dest [C.MAX_PATIENTS]C.Patient
	errCode := C.ListMalePatients(&s.patients[0], s.max_patients, &dest[0], &resultCount)
	if errCode != 0 {
		return nil, fmt.Errorf("error listing male patients: %s", ErrorDescription(errCode))
	}

	result := make([]Patient, resultCount)
	for i := 0; i < int(resultCount); i++ {
		result[i] = ParseCPatient(&dest[i])
	}

	return result, nil
}

// Wrapper for ListPatientsUnderAge function from patient_metrics.h
func (s *PatientService) ListPatientsUnderAge(ageLimit int) ([]Patient, error) {
	var resultCount C.size_t
	var dest [C.MAX_PATIENTS]C.Patient
	errCode := C.ListPatientsUnderAge(&s.patients[0], s.max_patients, C.int(ageLimit), &dest[0], &resultCount)
	if errCode != 0 {
		return nil, fmt.Errorf("error listing patients under age %d: %s", ageLimit, ErrorDescription(errCode))
	}

	result := make([]Patient, resultCount)
	for i := 0; i < int(resultCount); i++ {
		result[i] = ParseCPatient(&dest[i])
	}

	return result, nil
}
