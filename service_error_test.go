package ssmt

import (
	"errors"
	"testing"
)

func TestTranslateServiceError(t *testing.T) {
	var err IServiceError
	err = translateServiceError(1, "")
	t.Logf("%+v", err)
	if err != nil {
		t.FailNow()
	}
	err = translateServiceError(0, "")
	t.Logf("%+v", err)
	if err == nil || err.GetCode() != 0 {
		t.FailNow()
	}
	err = translateServiceError(-1, "")
	t.Logf("%+v", err)
	if err == nil || err.GetCode() != -1 {
		t.FailNow()
	}
	err = translateServiceError(5, "testMsg")
	t.Logf("%+v", err)
	if err == nil || err.GetMsg() != "testMsg" {
		t.FailNow()
	}
	err2 := translateServiceError(5, "testMsg2")
	t.Logf("%+v", err)
	if err2 == nil || !err2.Equal(err) || err2.Equal(errors.New("")) {
		t.FailNow()
	}
}
