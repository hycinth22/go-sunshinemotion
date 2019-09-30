package ssmt

import (
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
}
