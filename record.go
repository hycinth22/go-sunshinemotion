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
	fields := r.fieldsStringGroup()
	timestamp := strconv.FormatInt(r.Time.Unix(), 10)
	fields = append(fields, timestamp)
	fields = append(fields, r.DeviceName)
	fields = append(fields, r.DeviceIMEI)
	fields = append(fields, r.DeviceIMSI)
	if r.Xposed {
		fields = append(fields, "xposed")
	}
	if r.MockLocation {
		fields = append(fields, "mocklocation")
	}
	if r.Root {
		fields = append(fields, "root")
	}
	fields = append(fields, r.Custom...)
	return fields
}

// the transformation result is in the format of java language method java.util.ArrayList.toString()
// for example: [1546140832, Android,25,7.1.2, 263004834925257, 1234567890]
func (r *Remark) String() string {
	var b strings.Builder
	b.WriteString("[")
	for _, s := range r.Custom {
		b.WriteString(s)
	}
	b.WriteString("]")
	return b.String()
}

// a check code for Record
func (record *Record) XTcode() string {
	return crypto.CalcXTcode(record.UserID, toRPCTimeStr(record.BeginTime), toRPCDistanceStr(record.Distance))
}
