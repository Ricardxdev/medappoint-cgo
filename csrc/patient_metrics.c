#include "patient_metrics.h"
#include <string.h>

// Returns a list of patients with disabilities.
int ListDisabledPatients(const Patient* patients, size_t count, Patient* dest, size_t* result_count) {
    if (patients == NULL || count == 0 || dest == NULL || result_count == NULL) return -1;
    *result_count = 0;
    for (size_t i = 0; i < count; i++) {
        if (patients[i].disability == 1) {
            dest[(*result_count)++] = patients[i];
        }
    }
    return 0;
}

// Returns a list of patients with appointments on a given date (YYYY-MM-DD).
int ListPatientsByAppointmentDate(const Patient* patients, size_t count, const char* date, Patient* dest, size_t* result_count) {
    if (patients == NULL || count == 0 || date == NULL || dest == NULL || result_count == NULL) return -1;
    *result_count = 0;
    for (size_t i = 0; i < count; i++) {
        if (strcmp(patients[i].appointment_date, date) == 0) {
            dest[(*result_count)++] = patients[i];
        }
    }
    return 0;
}

// Returns a list of patients by doctor specialty.
int ListPatientsBySpecialty(const Patient* patients, size_t count, const char* specialty, Patient* dest, size_t* result_count) {
    if (patients == NULL || count == 0 || specialty == NULL || dest == NULL || result_count == NULL) return -1;
    *result_count = 0;
    for (size_t i = 0; i < count; i++) {
        if (strcmp(patients[i].doc_specialty, specialty) == 0) {
            dest[(*result_count)++] = patients[i];
        }
    }
    return 0;
}

// Returns a list of female patients.
int ListFemalePatients(const Patient* patients, size_t count, Patient* dest, size_t* result_count) {
    if (patients == NULL || count == 0 || dest == NULL || result_count == NULL) return -1;
    *result_count = 0;
    for (size_t i = 0; i < count; i++) {
        if (patients[i].gender == 'F') {
            dest[(*result_count)++] = patients[i];
        }
    }
    return 0;
}

// Returns a list of male patients.
int ListMalePatients(const Patient* patients, size_t count, Patient* dest, size_t* result_count) {
    if (patients == NULL || count == 0 || dest == NULL || result_count == NULL) return -1;
    *result_count = 0;
    for (size_t i = 0; i < count; i++) {
        if (patients[i].gender == 'M') {
            dest[(*result_count)++] = patients[i];
        }
    }
    return 0;
}

// Returns a list of patients under a certain age.
int ListPatientsUnderAge(const Patient* patients, size_t count, int age_limit, Patient* dest, size_t* result_count) {
    if (patients == NULL || count == 0 || dest == NULL || result_count == NULL) return -1;
    *result_count = 0;
    for (size_t i = 0; i < count; i++) {
        if (patients[i].age <= 0) continue; // Skip invalid ages
        if (patients[i].age < age_limit) {
            dest[(*result_count)++] = patients[i];
        }
    }
    return 0;
}
