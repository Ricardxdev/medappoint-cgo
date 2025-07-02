#ifndef ERRORS_H
#define ERRORS_H

typedef enum {
    ERR_NULL_PTR = 100,                     // Null pointer argument
    ERR_INVALID_ARG = 101,                  // Invalid argument (generic)
    ERR_OUT_OF_RANGE = 102,                 // Value out of allowed range
    ERR_ALLOC = 103,                        // Memory allocation failed
    ERR_IO = 104,                           // File I/O error
    ERR_DUPLICATE = 105,                    // Duplicate entry
    ERR_NOT_FOUND = 106,                    // Entry not found
    ERR_ASSIGN = 107,                       // Assignment to destination pointer failed

    // Field-specific validation errors (grouped for each field)
    ERR_FIELD_CI_NULL = 200,                // CI is NULL
    ERR_FIELD_CI_FORMAT = 201,              // CI must be exactly 8 digits
    ERR_FIELD_NAME_NULL = 202,              // Name is NULL
    ERR_FIELD_NAME_TOO_LONG = 203,          // Name is too long
    ERR_FIELD_AGE_INVALID = 204,            // Invalid age (must be >= 0)
    ERR_FIELD_GENDER_INVALID = 205,         // Invalid gender (must be 'M' or 'F')
    ERR_FIELD_DIAGNOSIS_NULL = 206,         // Diagnosis is NULL
    ERR_FIELD_DIAGNOSIS_TOO_LONG = 207,     // Diagnosis is too long
    ERR_FIELD_SPECIALTY_NULL = 208,         // Specialty is NULL
    ERR_FIELD_SPECIALTY_TOO_LONG = 209,     // Specialty is too long
    ERR_FIELD_APPOINTMENT_DATE_NULL = 210,  // Appointment date is NULL
    ERR_FIELD_APPOINTMENT_DATE_FORMAT = 211,// Appointment date must be YYYY-MM-DD (10 chars)

    // Additional context-specific error codes
    ERR_PARSE_LINE = 300,                   // Malformed or unreadable line in file
    ERR_INDEX_RANGE = 301                   // Hash/index out of allowed range
} ErrorCodes;

static inline const char* ErrorDescription(int code) {
    switch (code) {
        case ERR_NULL_PTR: return "Null pointer argument";
        case ERR_INVALID_ARG: return "Invalid argument";
        case ERR_OUT_OF_RANGE: return "Value out of allowed range";
        case ERR_ALLOC: return "Memory allocation failed";
        case ERR_IO: return "File I/O error";
        case ERR_DUPLICATE: return "Duplicate entry";
        case ERR_NOT_FOUND: return "Entry not found";
        case ERR_ASSIGN: return "Assignment to destination pointer failed";
        case ERR_FIELD_CI_NULL: return "CI is NULL";
        case ERR_FIELD_CI_FORMAT: return "CI must be exactly 8 digits";
        case ERR_FIELD_NAME_NULL: return "Name is NULL";
        case ERR_FIELD_NAME_TOO_LONG: return "Name is too long";
        case ERR_FIELD_AGE_INVALID: return "Invalid age (must be >= 0)";
        case ERR_FIELD_GENDER_INVALID: return "Invalid gender (must be 'M' or 'F')";
        case ERR_FIELD_DIAGNOSIS_NULL: return "Diagnosis is NULL";
        case ERR_FIELD_DIAGNOSIS_TOO_LONG: return "Diagnosis is too long";
        case ERR_FIELD_SPECIALTY_NULL: return "Specialty is NULL";
        case ERR_FIELD_SPECIALTY_TOO_LONG: return "Specialty is too long";
        case ERR_FIELD_APPOINTMENT_DATE_NULL: return "Appointment date is NULL";
        case ERR_FIELD_APPOINTMENT_DATE_FORMAT: return "Appointment date must be YYYY-MM-DD (10 chars)";
        case ERR_PARSE_LINE: return "Malformed or unreadable line in file";
        case ERR_INDEX_RANGE: return "Hash/index out of allowed range";
        default: return "Unknown error code";
    }
}

#endif // ERRORS_H