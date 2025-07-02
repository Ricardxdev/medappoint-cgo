#ifndef PATIENT_METRICS_H
#define PATIENT_METRICS_H

#include <stddef.h>
#include "patient.h"

// Returns a list of patients with disabilities.
int ListDisabledPatients(const Patient* patients, size_t count, Patient* dest, size_t* result_count);

// Returns a list of patients with appointments on a given date (YYYY-MM-DD).
int ListPatientsByAppointmentDate(const Patient* patients, size_t count, const char* date, Patient* dest, size_t* result_count);

// Returns a list of patients by doctor specialty.
int ListPatientsBySpecialty(const Patient* patients, size_t count, const char* specialty, Patient* dest, size_t* result_count);

// Returns a list of female patients.
int ListFemalePatients(const Patient* patients, size_t count, Patient* dest, size_t* result_count);

// Returns a list of male patients.
int ListMalePatients(const Patient* patients, size_t count, Patient* dest, size_t* result_count);

// Returns a list of patients under a certain age.
int ListPatientsUnderAge(const Patient* patients, size_t count, int age_limit, Patient* dest, size_t* result_count);

#endif // PATIENT_METRICS_H
