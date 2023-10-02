package errors

import "net/http"

type DBError struct{}

func (e *DBError) Error() string {
	return "DB Error"
}

func (e *DBError) StatusCode() int {
	return http.StatusInternalServerError
}

type BadTimestampError struct {
}

func (e *BadTimestampError) Error() string {
	return "timestamp format: \"2022-10-31T10:00:00Z\""
}

func (e *BadTimestampError) StatusCode() int {
	return http.StatusBadRequest
}
