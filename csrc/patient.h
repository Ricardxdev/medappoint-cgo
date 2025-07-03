#ifndef PATIENT_H
#define PATIENT_H

#include <stddef.h>

// ——————————————————————————————————————————————————————————————————————————————
// Constants & File names
// ——————————————————————————————————————————————————————————————————————————————
#define MAX_PATIENTS      100      // maximum number of patients
#define MAX_INDEX         1000     // size of index hash table
#define NAME_LEN          25      // max name length
#define DIAG_LEN          50      // max diagnosis length
#define SPEC_LEN          50      // max specialty length

#define PATIENT_FILE      "data/patients.bin"
#define INDEX_FILE        "data/index.dat"

// ——————————————————————————————————————————————————————————————————————————————
// Data Structures
// ——————————————————————————————————————————————————————————————————————————————
typedef struct {
    char ci[9];                    // 8 chars + NUL
    char name[NAME_LEN];
    int  age;
    char diagnosis[DIAG_LEN];
    char gender;                   // 'M' or 'F'
    int  disability;               // 0 or 1
    char doc_specialty[SPEC_LEN];
    char appointment_date[11];     // "YYYY-MM-DD"+NUL
} Patient;

typedef struct PatientIndex {
    char    ci[9];
    size_t  position;              // position in patients array/file
    int next;     // for collision chains (optional)
} PatientIndex;

typedef PatientIndex Index[MAX_INDEX];

// ——————————————————————————————————————————————————————————————————————————————
// Creation & Parsing
// ——————————————————————————————————————————————————————————————————————————————
// Initialize a Patient record.
//   dest:           output pointer to Patient
//   ci, name, …:    input fields
// returns 0 on success, error code otherwise
int NewPatient(
    Patient*       dest,
    const char*    ci,
    const char*    name,
    int            age,
    const char*    diagnosis,
    char           gender,
    int            disability,
    const char*    doc_specialty,
    const char*    appointment_date
);

// Parse an 8-digit CI string to a size_t.
//   dest: output pointer to parsed integer
// returns 0 on success, error code otherwise
int ParseCI(size_t* dest, const char* ci);

// Compute a hash from a CI string for the index.
//   dest: output pointer to hash (0..MAX_INDEX-1)
// returns 0 on success, error code otherwise
int Hash(size_t* dest, const char* str);

// ——————————————————————————————————————————————————————————————————————————————
// Index Management
// ——————————————————————————————————————————————————————————————————————————————
// Insert a new index entry.
//   index:    pointer to Index array
//   ci:       patient CI
//   position: patient’s position in file/array
// returns 0 on success, error code otherwise
int NewPatientIndex(Index* index, const char* ci, size_t position);

// Add, update or delete Patient in the in-memory array + index
int AddPatient(
    size_t*    count,
    Index*     index,
    Patient*   new_patient
);

int UpdatePatient(
    Index*     index,
    const char* ci,
    Patient*   updated_patient
);

int DeletePatient(
    Patient*    patients,
    Index*      index,
    const char* ci
);

// ——————————————————————————————————————————————————————————————————————————————
// Persistence
// ——————————————————————————————————————————————————————————————————————————————
// Save/load patients array to/from binary file
int SavePatients(Patient patients[], size_t patientsCount);
int LoadPatients(Patient* dest, size_t* dest_size);

// Save/load index to/from text file
int SaveIndex(Index* index);
int LoadIndex(Index* dest);

// Sync both files in one call
int SyncFiles(Patient patients[], size_t count, Index* index);

// ——————————————————————————————————————————————————————————————————————————————
// Queries & Display
// ——————————————————————————————————————————————————————————————————————————————
// Retrieve a single Patient by CI via index
int GetPatient(
    Patient*            p_dest,
    size_t*             i_dest,
    Index         index,
    const char*         ci
);

// Print helpers
void ShowPatient(const Patient* p);
int  ShowPatients(const Patient* patients, size_t count);

// ——————————————————————————————————————————————————————————————————————————————
// Appointments
// ——————————————————————————————————————————————————————————————————————————————
// Update a patient’s appointment date
int ScheduleAppointment(
    Patient*     patients,
    Index  index,
    const char*  ci,
    const char*  date
);

#endif // PATIENT_H