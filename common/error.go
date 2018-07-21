package common

type Error interface {
	Error() string
	Status() string
}

// DetailedError contains an error with useful information, like the appropriate
// HTTP status code
type DetailedError struct {
	ID      string `json:"id"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func (d *DetailedError) Error() string {
	return d.Message
}

func NewDetailedError(id string, message string, status int) *DetailedError {
	return &DetailedError{
		ID:      id,
		Message: message,
		Status:  status,
	}
}
