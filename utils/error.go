package utils

import "jaeger-tenant/pkg/structures"

// CreateErrorResponse returns Response containing error
func CreateErrorResponse(err error) *structures.Response {
	return &structures.Response{
		Error: err.Error(),
	}
}

// CreateMessageResponse returns Response containing msg
func CreateMessageResponse(msg string) *structures.Response {
	return &structures.Response{
		Message: msg,
	}
}
