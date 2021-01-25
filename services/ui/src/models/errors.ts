

enum ErrorType {
    BadRequest = "ERROR_BAD_REQUEST",
    InternalError = "ERROR_INTERNAL",
    NotFound = "ERROR_NOT_FOUND",
    ValidationError = "ERROR_VALIDATION",
}

enum ValidationErrors {
    Required = "property is required",
}

// Error when calling the api.
export interface IApiError {
	// The type of error that happened.
    type?: ErrorType;
	// The description of the error.
    error?: string;
	// A map of data attached to the error. In case of a validation errors this contains the errors per validated property.
    data?: Map<string, string>;
}