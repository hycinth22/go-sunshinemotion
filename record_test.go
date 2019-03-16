package lib

import (
	"testing"
	"time"
)

func TestSmartCreateRecords(t *testing.T) {
	s := CreateSession()
	s.UserInfo.Sex = "F"; s.UpdateLimitParams()
	records_f := SmartCreateRecords(0, 0, s.LimitParams, s.LimitParams.LimitTotalDistance.Max - 0.1, time.Now())
	for _, r := range records_f {
		t.Logf("%+v", r)
	}
	s.UserInfo.Sex = "M"; s.UpdateLimitParams()
	records_m := SmartCreateRecords(0, 0, s.LimitParams, s.LimitParams.LimitTotalDistance.Max - 0.1, time.Now())
	for _, r := range records_m {
		t.Logf("%+v", r)
	}
}
