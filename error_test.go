package wavecell

import (
	"fmt"
	"github.com/fairyhunter13/iso8601/v2"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestResponseError_Error(t *testing.T) {
	now := iso8601.Time{time.Now().UTC()}
	respErr := &ResponseError{
		Code:      1004,
		Message:   "Error from the Wavecell.",
		ErrorID:   "ERROR",
		Timestamp: now,
	}
	expected := fmt.Sprintf(
		"error from the Wavecell platform, code: %d, message: %s, error ID: %s, timestamp: %s",
		respErr.Code,
		respErr.Message,
		respErr.ErrorID,
		respErr.Timestamp,
	)
	assert.Equal(t, expected, respErr.Error())
}
