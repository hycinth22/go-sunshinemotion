package sunshinemotion

import (
	"testing"
	"time"
)

func TestMakeBasicRemark(t *testing.T) {
	remarkTime := time.Now()
	device := GenerateDevice()
	remark := MakeBasicRemark(remarkTime, device)
	if len(remark.Custom) != 0 {
		t.FailNow()
	}
	if remark.Time != remarkTime {
		t.FailNow()
	}
	if remark.DeviceName != device.Name {
		t.FailNow()
	}
	if remark.DeviceIMEI != device.IMEI {
		t.FailNow()
	}
	if remark.DeviceIMSI != device.IMSI {
		t.FailNow()
	}
}

func TestRemark_String(t *testing.T) {
	t.SkipNow() // TODO: Implement me
	const timeLayout = "2006-01-02 15:04:05"
	tests := []struct {
		name         string
		remark       Remark
		remarkString string
	}{
		{"Basic Remark",
			Remark{},
			""},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expect := test.remarkString
			actual := test.remark.String()
			if actual != expect {
				t.Log("incorrect result: ", actual)
				t.Fail()
			}
		})
	}
}

func TestRecord_XTcode(t *testing.T) {
	t.SkipNow() // TODO: Implement me
	tests := []struct {
		name   string
		record Record
		xtcode string
	}{
		// TODO: test cases
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expect := test.xtcode
			actual := test.record.XTcode()
			if actual != expect {
				t.Log("incorrect result: ", actual)
				t.Fail()
			}
		})
	}
}
