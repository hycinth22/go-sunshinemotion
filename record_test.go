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
	const timeLayout = "2006-01-02 15:04:05"
	type testCase struct {
		name         string
		remark       Remark
		remarkString string
	}
	tests := []struct {
		name         string
		remark       Remark
		remarkString string
	}{
		{
			"Basic Remark",
			Remark{
				Time:       time.Unix(0, 0),
				DeviceName: "TestDevice",
				DeviceIMEI: "263004834925257",
				DeviceIMSI: "1234567890",
			},
			"[0, TestDevice, 263004834925257, 1234567890]",
		},
		{
			"Cheated Remark",
			Remark{
				Time:         time.Unix(1547705586, 0),
				DeviceName:   "TestDevice",
				DeviceIMEI:   "263004834925257",
				DeviceIMSI:   "1234567890",
				Xposed:       true,
				MockLocation: true,
				Root:         true,
			},
			"[1547705586, TestDevice, 263004834925257, 1234567890, xposed, mocklocation, root]",
		},
		{
			"Custom Remark",
			Remark{
				Time:       time.Unix(1547705599, 0),
				DeviceName: "TestDevice",
				DeviceIMEI: "263004834925257",
				DeviceIMSI: "1234567890",
				Custom: []string{
					"CustomField1",
					"MyField2",
				},
			},
			"[1547705599, TestDevice, 263004834925257, 1234567890, CustomField1, MyField2]",
		},
		{
			"Cheated Custom Remark",
			Remark{
				Time:         time.Unix(1547705606, 0),
				DeviceName:   "TestDevice",
				DeviceIMEI:   "263004834925257",
				DeviceIMSI:   "1234567890",
				Xposed:       true,
				MockLocation: true,
				Root:         true,
				Custom: []string{
					"CustomField1",
					"MyField2",
				},
			},
			"[1547705606, TestDevice, 263004834925257, 1234567890, xposed, mocklocation, root, CustomField1, MyField2]",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expect := test.remarkString
			actual := test.remark.String()
			if actual != expect {
				t.Fatal("incorrect result: ", actual)
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
