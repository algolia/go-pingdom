package solarwinds

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClientErrors(t *testing.T) {
	user := RandString(10) + "@algolia.com"
	errs := []error{
		NewErrorAttemptDeleteActiveUser(user),
		NewNetworkError(errors.New("underlying network error")),
	}
	expectedErrMsg := []string{
		fmt.Sprintf("status: %d, err: deleting active user %v is not supported", ErrCodeDeleteActiveUserException, user),
		fmt.Sprintf("status: %d, err: underlying network error", ErrCodeNetworkException),
	}
	for i, err := range errs {
		if err != nil {
			errMsg := err.Error()
			assert.Equal(t, errMsg, expectedErrMsg[i])
		} else {
			t.Error("should not reach here")
		}
	}
}
