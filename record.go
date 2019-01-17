package sunshinemotion

import (
	"strconv"
	"strings"
	"time"

	"github.com/inkedawn/go-sunshinemotion/crypto"
)

// an record which represents a sport progress.
type Record struct {
	UserID    uint
	Distance  float64
	BeginTime time.Time
	EndTime   time.Time
	Remark    Remark
}

// Remark need to meet some requirements, otherwise illegal.
// see MakeBasicRemark() for details.
type Remark struct {
	Time                               time.Time
	DeviceName, DeviceIMEI, DeviceIMSI string

	Xposed, MockLocation, Root bool
	Custom                     []string
}

const (
	remark_xposed       = "xposed"
	remark_mockLocation = "mocklocation"
	remark_root         = "root"
)

// The Remark is an Ordered String Group,
// its First element MUST be the Unix Timestamp (seconds),
// Second element MUST be the Device Name,
// Third element MUST be the Device IMEI,
// Forth element MUST be the Device IMSI,
// and the rest elements are other flags.
//
// Finally it will be transformed into a string in a specific format.
// you can call remark.String() to perform the transformation, see it for transformation details.
func MakeBasicRemark(time time.Time, device *Device) Remark {
	return Remark{
		time,
		device.Name, device.IMEI, device.IMSI,
		false, false, false,
		nil,
	}
}

func (r *Remark) fieldsStringGroup() []string {
	fields := make([]string, 0, 4)
	timestamp := strconv.FormatInt(r.Time.Unix(), 10)
	fields = append(fields, timestamp)
	fields = append(fields, r.DeviceName)
	fields = append(fields, r.DeviceIMEI)
	fields = append(fields, r.DeviceIMSI)
	// the order is not changeable
	if r.Xposed {
		fields = append(fields, remark_xposed)
	}
	if r.MockLocation {
		fields = append(fields, remark_mockLocation)
	}
	if r.Root {
		fields = append(fields, remark_root)
	}
	fields = append(fields, r.Custom...)
	return fields
}

// the transformation result is in the format of java language method java.util.ArrayList.toString()
// for example: [1546140832, Android,25,7.1.2, 263004834925257, 1234567890]
func (r *Remark) String() string {
	fileds := r.fieldsStringGroup()
	var b strings.Builder
	b.WriteString("[")
	for i, s := range fileds {
		b.WriteString(s)
		if i != len(fileds)-1 {
			b.WriteString(", ")
		}
	}
	b.WriteString("]")
	return b.String()
}

// a check code for Record
func (record *Record) XTcode() string {
	return crypto.CalcXTcode(record.UserID, toRPCTimeStr(record.BeginTime), toRPCDistanceStr(record.Distance))
}
