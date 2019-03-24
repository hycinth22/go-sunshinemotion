package ssmt

import (
	"testing"
	"time"
)

func TestSmartCreateRecordsBefore(t *testing.T) {
	tests := []struct {
		sex string
	}{
		{"F"}, {"M"},
	}
	for _, test := range tests {
		t.Run(test.sex, func(t *testing.T) {
			t.Logf(test.sex)
			limit := GetDefaultLimitParams(test.sex)
			records := SmartCreateRecordsBefore(0, 0, limit, limit.LimitTotalDistance.Max-0.1, time.Now())
			for _, r := range records {
				t.Logf("%+v", r)
			}

		})
	}
}

func TestSmartCreateRecordsAfter(t *testing.T) {
	tests := []struct {
		sex string
	}{
		{"F"}, {"M"},
	}
	for _, test := range tests {
		t.Run(test.sex, func(t *testing.T) {
			t.Logf(test.sex)
			limit := GetDefaultLimitParams(test.sex)
			records := SmartCreateRecordsAfter(0, 0, limit, limit.LimitTotalDistance.Max-0.1, time.Now())
			for _, r := range records {
				t.Logf("%+v", r)
			}

		})
	}
}
