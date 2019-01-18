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
	tests := []struct {
		name   string
		record Record
		xtcode string
	}{
		{"empty case",
			Record{
				UserID:    0,
				BeginTime: time.Unix(0, 0),
				Distance:  0.0},
			"146618DF0D8DF1DC698E0D70FF5F7F4B",
		},
		{"normal case 1",
			Record{
				UserID:    1234,
				BeginTime: time.Date(2006, 1, 2, 15, 4, 5, 0, time.Local),
				Distance:  3.410,
			},
			"ED9335F08B363EF8C9954302EE84547F",
		},
		{"normal case 2",
			Record{
				UserID:    5678,
				BeginTime: time.Date(2017, 2, 3, 16, 5, 6, 0, time.Local),
				Distance:  4.520,
			},
			"38A2FDC3A89C6200CE41BF0639F145ED",
		},
		{"normal case 3",
			Record{
				UserID:    9012,
				BeginTime: time.Date(2028, 3, 4, 17, 6, 7, 0, time.Local),
				Distance:  5.630,
			},
			"1630BC477C9BB0150773434F9A07C0DD",
		},
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
