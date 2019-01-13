package sunshinemotion

import (
	"errors"
	"testing"
)

func TestIsServiceError(t *testing.T) {
	if !IsServiceError(serviceError{}) {
		t.FailNow()
	}
	if IsServiceError(errors.New("0")) {
		t.FailNow()
	}
}

func TestIsNetworkErrorError(t *testing.T) {
	if !IsNetworkError(networkError{}) {
		t.FailNow()
	}
	if IsNetworkError(errors.New("0")) {
		t.FailNow()
	}
}
