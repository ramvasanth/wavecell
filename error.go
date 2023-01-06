package wavecell

import (
	"fmt"
	"github.com/fairyhunter13/iso8601/v2"
	"github.com/pkg/errors"
)

var (
	// ErrEmptyAPIKEY is the error for empty API key.
	ErrEmptyAPIKEY = errors.New("API key is empty")
	// ErrEmptySubAccountID is the error for empty sub account ID.
	ErrEmptySubAccountID = errors.New("Sub account ID is empty")
)

// ResponseError is the standard response struct for error.
type ResponseError struct {
	Code      int          `json:"code"`
	Message   string       `json:"message,omitempty"`
	ErrorID   string       `json:"errorId"`
	Timestamp iso8601.Time `json:"timestamp"`
}

// Error returns the error message.
func (r *ResponseError) Error() (res string) {
	res = fmt.Sprintf(
		"error from the Wavecell platform, code: %d, message: %s, error ID: %s, timestamp: %s",
		r.Code,
		r.Message,
		r.ErrorID,
		r.Timestamp,
	)
	return
}
