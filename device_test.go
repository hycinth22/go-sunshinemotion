package sunshinemotion

import (
	"testing"
)

func TestGenerateDevice(t *testing.T) {
	device := GenerateDevice()
	t.Log(*device)
	if device == nil {
		t.FailNow()
	}
	if len(device.Name) == 0 || len(device.Model) == 0 || len(device.Screen) == 0 {
		t.FailNow()
	}
	if len(device.IMEI) != 15 {
		t.FailNow()
	}
	// The length of the normal IMSI is 15. (except emptyIMSI, which is 1234567890)
	if len(device.IMSI) != 15 && device.IMSI != emptyIMSI {
		t.FailNow()
	}
	if len(device.UserAgent) == 0 {
		t.FailNow()
	}
}
