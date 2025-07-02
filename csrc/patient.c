#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <errno.h>
#include <ctype.h>
#include "patient.h"
#include "errors.h"

int NewPatient(
    Patient* dest,
    const char* ci,
    const char* name,
    int age,
    const char* diagnosis,
    char gender,
    int disability,
    const char* doc_specialty,
    const char* appointment_date
) {
    // Check if destination pointer is not null
    if (dest == NULL) return ERR_NULL_PTR;
    if (ci == NULL) return ERR_FIELD_CI_NULL;
    if (strlen(ci) != 8) return ERR_FIELD_CI_FORMAT;
    if (strspn(ci, "0123456789") != 8) return ERR_FIELD_CI_FORMAT;
    if (name == NULL) return ERR_FIELD_NAME_NULL;
    if (strlen(name) == 0) return ERR_FIELD_NAME_NULL;
    if (strlen(name) > NAME_LEN) return ERR_FIELD_NAME_TOO_LONG;
    if (age < 0) return ERR_FIELD_AGE_INVALID;
    if (diagnosis == NULL) return ERR_FIELD_DIAGNOSIS_NULL;
    if (strlen(diagnosis) == 0) return ERR_FIELD_DIAGNOSIS_NULL;
    if (strlen(diagnosis) > DIAG_LEN) return ERR_FIELD_DIAGNOSIS_TOO_LONG;
    if (gender != 'M' && gender != 'F') return ERR_FIELD_GENDER_INVALID;
    if (disability != 0 && disability != 1) return ERR_INVALID_ARG;
    if (doc_specialty == NULL) return ERR_FIELD_SPECIALTY_NULL;
    if (strlen(doc_specialty) == 0) return ERR_FIELD_SPECIALTY_NULL;
    if (strlen(doc_specialty) > SPEC_LEN) return ERR_FIELD_SPECIALTY_TOO_LONG;
    if (appointment_date == NULL) return ERR_FIELD_APPOINTMENT_DATE_NULL;
    if (strlen(appointment_date) != 10) return ERR_FIELD_APPOINTMENT_DATE_FORMAT;

    Patient p;
    memset(&p, 0, sizeof(Patient)); // Initialize all fields to zero
    strcpy(p.ci, ci);
    strcpy(p.name, name);
    p.age = age;
    strcpy(p.diagnosis, diagnosis);
    p.gender = gender;
    p.disability = disability;
    strcpy(p.doc_specialty, doc_specialty);
    strcpy(p.appointment_date, appointment_date);

    *dest = p;

    return 0;
}

int ParseCI(size_t* dest, const char* ci) {
    if (ci == NULL) return ERR_FIELD_CI_NULL;
    if (strlen(ci) != 8) return ERR_FIELD_CI_FORMAT;
    for (int i = 0; i < 8; i++) {
        if (!isdigit(ci[i])) return ERR_FIELD_CI_FORMAT;
    }
    *dest = atoi(ci);
    if (*dest < 0) return ERR_FIELD_CI_FORMAT;
    return 0;
}

int Hash(size_t* dest, const char* str) {
    size_t hash;
    int err = ParseCI(&hash, str);
    if (err != 0) return err; // Error in parsing CI
    *dest = hash % MAX_INDEX; // Example: limit to 1000 buckets
    return 0;
}

int NewPatientIndex(
    Index* index,
    const char* ci,
    size_t position
) {
    // Check if destination pointer is not null
    if (index == NULL) return ERR_NULL_PTR;
    if (ci == NULL) return ERR_FIELD_CI_NULL;
    if (strlen(ci) != 8) return ERR_FIELD_CI_FORMAT;

    PatientIndex p;
    strcpy(p.ci, ci);
    if (!p.ci) {
        return ERR_ALLOC;
    }
    p.position = position;

    size_t hash;
    int error = Hash(&hash, p.ci);
    if (error != 0) {
        return error;
    }

    (*index)[hash] = p;
    if ((*index)[hash].ci[0] == '\0') return ERR_ASSIGN;

    return 0;
}

void FreePatient(Patient* p) {
    if (p) {
        free(p);
    }
}
    
void SortPatients(Patient* arr, int left, int right) {
    if (right <= left) return; // nada que ordenar
    // quicksort clásico solo sobre [left, right]
    int i = left, j = right;
    Patient pivot = arr[(left + right) / 2];
    while (i <= j) {
        while (strcmp(arr[i].ci, pivot.ci) < 0) i++;
        while (strcmp(arr[j].ci, pivot.ci) > 0) j--;
        if (i <= j) {
            Patient tmp = arr[i];
            arr[i] = arr[j];
            arr[j] = tmp;
            i++;
            j--;
        }
    }
    if (left < j) SortPatients(arr, left, j);
    if (i < right) SortPatients(arr, i, right);
}

int SavePatients(Patient patients[], size_t patientsCount) {
    if (patients == NULL) return ERR_NULL_PTR;
    SortPatients(patients, 0, patientsCount - 1);
    FILE* file = fopen(PATIENT_FILE, "wb");
    if (file == NULL) return ERR_IO;
    printf("Saving %zu patients to %s\n", patientsCount, PATIENT_FILE);
    for (size_t i = 0; i < patientsCount; i++) {
        if (patients[i].age == 0) {
            continue;
        }
        printf("Saving patient %zu: CI=%s, Name=%s, Age=%d\n", i, patients[i].ci, patients[i].name, patients[i].age);
        if (fwrite(&patients[i], sizeof(Patient), 1, file) != 1) {
            fclose(file);
            return ERR_IO;
        }
    }
    fclose(file);
    return 0;
}

int SaveIndex(Index* index) {
    FILE *file = fopen(INDEX_FILE, "w");
    if (!file) return ERR_IO;
    for (size_t i = 0; i < MAX_INDEX; i++) {
        if ((*index)[i].ci[0] == '\0') continue;
        if (fprintf(file, "|%s|%zu|\n", (*index)[i].ci, (*index)[i].position) < 0) {
            fclose(file);
            return ERR_IO;
        }
    }
    fclose(file);
    return 0;
}

int LoadIndex(
    Index* dest
) {
    if (dest == NULL) return ERR_NULL_PTR;
    FILE *file = fopen(INDEX_FILE, "r");
    if (file == NULL) return ERR_IO;
    #define LINE_BUF_SIZE 256
    char line[LINE_BUF_SIZE];
    while (fgets(line, sizeof(line), file)) {
        PatientIndex idx;
        if (line[0] == '\0' || line[0] == '\n') continue;
        if (sscanf(line, "|%8s|%zu|", idx.ci, &idx.position) != 2) {
            fclose(file);
            return ERR_PARSE_LINE;
        }
        size_t error = NewPatientIndex(dest, idx.ci, idx.position);
        if (error != 0) {
            fclose(file);
            return error;
        }
    }
    if (ferror(file)) {
        fclose(file);
        return ERR_PARSE_LINE;
    }
    fclose(file);
    return 0;
}

int GetPatient(Patient* p_dest, size_t* i_dest, const Index index, const char* ci) {
    if (p_dest == NULL || i_dest == NULL || ci == NULL) return ERR_NULL_PTR;
    if (index == NULL) return ERR_NULL_PTR;
    size_t hash = 0;
    int err = Hash(&hash, ci);
    if (err != 0) return ERR_INVALID_ARG;
    size_t position = index[hash].position;
    printf("Hash position for CI %s: %zu, File position: %zu\n", ci, hash, position);
    FILE* file = fopen(PATIENT_FILE, "rb");
    if (file == NULL) return ERR_IO;
    fseek(file, position * sizeof(Patient), SEEK_SET);
    if (fread(p_dest, sizeof(Patient), 1, file) != 1) {
        fclose(file);
        return ERR_IO;
    }
    *i_dest = hash;
    fclose(file);
    return 0;
}

int AddPatient(
    Patient* patients,
    size_t* count,
    Index* index,
    Patient* new_patient
) {
    if (patients == NULL || count == NULL || new_patient == NULL) {
        return ERR_NULL_PTR;
    }
    if (*count >= MAX_PATIENTS) {
        return ERR_OUT_OF_RANGE;
    }
    patients[*count] = *new_patient;
    size_t hash;
    int error = Hash(&hash, new_patient->ci);
    if (error != 0) {
        return ERR_INVALID_ARG;
    }
    (*index)[hash].position = *count;
    strcpy((*index)[hash].ci, new_patient->ci);
    (*count)++;
    return 0;
}

int UpdatePatient(
    Patient* patients,
    Index* index,
    const char* ci,
    Patient* updated_patient
) {
    if (patients == NULL || index == NULL || ci == NULL || updated_patient == NULL) {
        return ERR_NULL_PTR;
    }
    size_t hash;
    int error = Hash(&hash, ci);
    if (error != 0) {
        return ERR_INVALID_ARG;
    }
    if ((*index)[hash].ci[0] == '\0' || strcmp((*index)[hash].ci, ci) != 0) {
        return ERR_NOT_FOUND;
    }
    size_t position = (*index)[hash].position;
    patients[position] = *updated_patient;
    return 0;
}

int SyncFiles(
    Patient patients[],
    size_t count,
    Index* index
) {
    if (patients == NULL || index == NULL) {
        return ERR_NULL_PTR;
    }
    int error = SavePatients(patients, count);
    if (error != 0) {
        return error;
    }
    error = SaveIndex(index);
    if (error != 0) {
        return error;
    }
    return 0;
}

int LoadPatients(Patient* dest, size_t* dest_size) {
    if (dest == NULL) return ERR_NULL_PTR;
    FILE* file = fopen(PATIENT_FILE, "rb");
    if (file == NULL) {
        int errnum = errno;                              // capture errno
        return errnum;    // or return errnum if you want to propagate the raw errno
    }
    size_t count = 0;
    while (fread(&dest[count], sizeof(Patient), 1, file) == 1) {
        count++;
    }
    if (ferror(file)) {
        int errnum = errno;                              // capture errno
        fclose(file);
        return errnum;
    }
    fclose(file);
    *dest_size = count;
    return 0;
}

void ShowPatient(const Patient* p) {
    if (p == NULL) return; // Error: null pointer
    if (p->age == 0) {
        printf("No patient data available.\n");
        return;
    }
    printf("CI: %s\n", p->ci);
    printf("Name: %s\n", p->name);
    printf("Age: %d\n", p->age);
    printf("Diagnosis: %s\n", p->diagnosis);
    printf("Gender: %c\n", p->gender);
    printf("Disability: %d\n", p->disability);
    printf("Doctor Specialty: %s\n", p->doc_specialty);
    printf("Appointment Date: %s\n", p->appointment_date);
    return;
}

int ShowPatients(const Patient* patients, size_t count) {
    printf("Showing %zu patients:\n", count);
    if (patients == NULL) return ERR_NULL_PTR;
    if (count == 0) return ERR_OUT_OF_RANGE;
    for (size_t i = 0; i < count; i++) {
        if (patients[i].ci[0] == '\0' || !isdigit(patients[i].ci[0])) continue; // Skip empty entries
        printf("=====================\n");
        printf("Patient %zu:\n", i + 1);
        ShowPatient(&patients[i]);
        printf("\n");
    }
    return 0;
}

int DeletePatient(Patient* patients, Index* index, const char* ci) {
    if (patients == NULL) return ERR_NULL_PTR;
    if (ci == NULL) return ERR_FIELD_CI_NULL;
    Patient empty_patient;
    memset(&empty_patient, 0, sizeof(Patient));
    int error = UpdatePatient(patients, index, ci, &empty_patient);
    if (error != 0) return error;
    size_t hash;
    error = Hash(&hash, ci);
    if (error != 0) {
        return ERR_INVALID_ARG;
    }
    PatientIndex empty_index;
    memset(&empty_index, 0, sizeof(PatientIndex));
    (*index)[hash] = empty_index;
    return 0;
}

int ScheduleAppointment(Patient* patients, Index index, const char* ci, const char* date) {
    if (patients == NULL) return ERR_NULL_PTR;
    if (ci == NULL) return ERR_FIELD_CI_NULL;
    if (date == NULL) return ERR_FIELD_APPOINTMENT_DATE_NULL;
    Patient patient;
    size_t index_position;
    int error = GetPatient(&patient, &index_position, index, ci);
    if (error != 0) return error;
    strcpy(patient.appointment_date, date);
    error = UpdatePatient(patients, (Index*)index, ci, &patient);
    if (error != 0) return error;
    return 0;
}


#ifndef CGO_BUILD
int main() {
    // load the index
    Index index;
    memset(index, 0, sizeof(index)); // Initialize the index array
    int error = 0;
    // Load Patients from the binary file
    size_t patient_count = 0;
    Patient patients[MAX_PATIENTS];
    memset(patients, 0, sizeof(patients)); // Initialize the patients array
    error = LoadPatients(patients, &patient_count);
    if (error != 0) {
        printf("Error loading patients: %d\n", error);
        return error;
    }
    printf("Loaded %zu patients from binary file.\n", patient_count);
    // Create index entries for each patient
    for (size_t i = 0; i < patient_count; i++) {
        if (patients[i].age == 0) continue; // Skip empty entries
        error = NewPatientIndex(&index, patients[i].ci, i);
        if (error != 0) {
            printf("Error creating index for patient %zu: %d\n", i, error);
            printf("Skipping patient %s.\n", patients[i].ci);
            continue; // Skip this patient and continue with the next
        }
    }
    printf("Index created with %zu entries.\n", patient_count);

    // // Save the index to a file
    // error = SaveIndex(&index);
    // if (error != 0) {
    //     printf("Error saving index: %d\n", error);
    //     return error;
    // }
    // printf("Index saved to %s.\n", INDEX_FILE);

    // memset(index, 0, sizeof(index)); // Initialize the index array
    // // Load the index from the file
    // error = LoadIndex(&index);
    // if (error != 0) {
    //     printf("Error loading index: %d\n", error);
    //     return error;
    // }
    // printf("Index loaded with %u entries.\n", MAX_INDEX);

    // // Print the index entries
    // printf("Index entries:\n");
    // for (size_t i = 0; i < MAX_INDEX; i++) {
    //     if (index[i].ci[0] == '\0') continue; // Skip empty entries
    //     printf("Index %zu: CI=%s, Position=%zu\n", i, index[i].ci, index[i].position);
    // }

    // // Example of retrieving a patient using the index
    // for (size_t i = 0; i < MAX_INDEX; i++) {
    //     if (index[i].ci[0] == '\0') continue; // Skip empty entries
    //     printf("Index %zu: CI=%s, Position=%zu\n", i, index[i].ci, index[i].position);
    //     Patient p_temp;
    //     size_t index_position;
    //     int error = GetPatient(&p_temp, &index_position, index, index[i].ci);
    //     if (error == 0) {
    //         printf("Patient found: %s, %s, %d\n", p_temp.ci, p_temp.name, p_temp.age);
    //     } else {
    //         printf("Error retrieving patient for CI: %s, error: %d\n", index[i].ci, error);
    //     }
    // }

    // // Patient* p = NULL;
    // // size_t result = NewPatient(&p, "12345678", "John Doe", 30, "Flu", 'M', 0, "Cardiology", "2023-10-01");
    // // if (result != NEWPATIENT_SUCCESS) {
    // //     printf("Error creating patient: %s\n", NewPatientErrorDescription(result));
    // //     return result;
    // // }

    // // error = AddPatient(patients, &patient_count, &index, p);
    // // if (error != 0) {
    // //     printf("Error adding patient: %d\n", error);
    // //     FreePatient(p);
    // //     return error;
    // // }

    // Show all patients
    error = ShowPatients(patients, patient_count);
    if (error != 0) {
        printf("Error showing patients: %d\n", error);
        return error;
    }

    error = SyncFiles(patients, patient_count, &index);
    if (error != 0) {
        printf("Error syncing files: %d\n", error);
        return error;
    }

    return 0;
}
#endif

int GeneratePatients() {
    // load the index
    Index index;
    memset(index, 0, sizeof(index)); // Initialize the index array
    int error = 0;
    // Load Patients from the binary file
    size_t patient_count = 15;
    Patient patients[] = {
    {"12345678", "Alice Johnson", 34, "Hypertension", 'F', 0, "Cardiology", "2023-02-15"},
    {"87654321", "Bob Smith", 47, "Diabetes Type 2", 'M', 1, "Endocrinology", "2023-03-10"},
    {"11223344", "Carla Gomez", 29, "Asthma", 'F', 0, "Pulmonology", "2023-04-22"},
    {"44332211", "Daniel Lee", 52, "Coronary Artery Disease", 'M', 1, "Cardiology", "2023-05-05"},
    {"55667788", "Emily Chen", 41, "Hypothyroidism", 'F', 0, "Endocrinology", "2023-06-18"},
    {"88776655", "Frank Miller", 65, "COPD", 'M', 1, "Pulmonology", "2023-07-12"},
    {"33445566", "Grace Kim", 23, "Migraine", 'F', 0, "Neurology", "2023-08-03"},
    {"66554433", "Henry Patel", 38, "Epilepsy", 'M', 0, "Neurology", "2023-09-27"},
    {"77889900", "Isabella Rossi", 56, "Osteoarthritis", 'F', 1, "Rheumatology", "2023-10-14"},
    {"00998877", "Jack Wilson", 44, "Chronic Kidney Disease", 'M', 0, "Nephrology", "2023-11-21"},
    {"22334455", "Karen Davis", 31, "Depression", 'F', 0, "Psychiatry", "2023-12-09"},
    {"55443322", "Luis Martinez", 27, "Ulcerative Colitis", 'M', 1, "Gastroenterology", "2024-01-16"},
    {"66778899", "Maria Silva", 49, "Breast Cancer", 'F', 0, "Oncology", "2024-02-28"},
    {"99887766", "Noah Brown", 36, "Multiple Sclerosis", 'M', 1, "Neurology", "2024-03-19"},
    {"13572468", "Olivia Clark", 58, "Glaucoma", 'F', 0, "Ophthalmology", "2024-04-07"}
};
    //Patient patients[MAX_PATIENTS];
    // memset(patients, 0, sizeof(patients)); // Initialize the patients array
    // error = LoadPatients(patients, &patient_count);
    // if (error != 0) {
    //     printf("Error loading patients: %d\n", error);
    //     return error;
    // }
    printf("Loaded %zu patients from binary file.\n", patient_count);
    // Create index entries for each patient
    for (size_t i = 0; i < patient_count; i++) {
        if (patients[i].age == 0) continue; // Skip empty entries
        error = NewPatientIndex(&index, patients[i].ci, i);
        if (error != 0) {
            printf("Error creating index for patient %zu: %d\n", i, error);
            printf("Skipping patient %s.\n", patients[i].ci);
            continue; // Skip this patient and continue with the next
        }
    }
    printf("Index created with %zu entries.\n", patient_count);

    // // Save the index to a file
    // error = SaveIndex(&index);
    // if (error != 0) {
    //     printf("Error saving index: %d\n", error);
    //     return error;
    // }
    // printf("Index saved to %s.\n", INDEX_FILE);

    // memset(index, 0, sizeof(index)); // Initialize the index array
    // // Load the index from the file
    // error = LoadIndex(&index);
    // if (error != 0) {
    //     printf("Error loading index: %d\n", error);
    //     return error;
    // }
    // printf("Index loaded with %u entries.\n", MAX_INDEX);

    // // Print the index entries
    // printf("Index entries:\n");
    // for (size_t i = 0; i < MAX_INDEX; i++) {
    //     if (index[i].ci[0] == '\0') continue; // Skip empty entries
    //     printf("Index %zu: CI=%s, Position=%zu\n", i, index[i].ci, index[i].position);
    // }

    // // Example of retrieving a patient using the index
    // for (size_t i = 0; i < MAX_INDEX; i++) {
    //     if (index[i].ci[0] == '\0') continue; // Skip empty entries
    //     printf("Index %zu: CI=%s, Position=%zu\n", i, index[i].ci, index[i].position);
    //     Patient p_temp;
    //     size_t index_position;
    //     int error = GetPatient(&p_temp, &index_position, index, index[i].ci);
    //     if (error == 0) {
    //         printf("Patient found: %s, %s, %d\n", p_temp.ci, p_temp.name, p_temp.age);
    //     } else {
    //         printf("Error retrieving patient for CI: %s, error: %d\n", index[i].ci, error);
    //     }
    // }

    // // Patient* p = NULL;
    // // size_t result = NewPatient(&p, "12345678", "John Doe", 30, "Flu", 'M', 0, "Cardiology", "2023-10-01");
    // // if (result != NEWPATIENT_SUCCESS) {
    // //     printf("Error creating patient: %s\n", NewPatientErrorDescription(result));
    // //     return result;
    // // }

    // // error = AddPatient(patients, &patient_count, &index, p);
    // // if (error != 0) {
    // //     printf("Error adding patient: %d\n", error);
    // //     FreePatient(p);
    // //     return error;
    // // }

    // Show all patients
    error = ShowPatients(patients, patient_count);
    if (error != 0) {
        printf("Error showing patients: %d\n", error);
        return error;
    }

    error = SyncFiles(patients, patient_count, &index);
    if (error != 0) {
        printf("Error syncing files: %d\n", error);
        return error;
    }

    return 0;
}