package models

// ServiceError is a general Service error response struct
type ServiceError struct {
	Code      int
	ErrorCode string
	Message   string
	Status    string
}
