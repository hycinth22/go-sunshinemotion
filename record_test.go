package ssmt

import (
	"testing"
	"time"
)

func checkRecordSpeed(t *testing.T, r Record, MinutePerKM Float64Range) {
	dur := r.EndTime.Sub(r.BeginTime)
	v := dur.Minutes() / r.Distance
	if v < MinutePerKM.Min {
		t.Logf("too small speed %v min/km", v)
		t.Fail()
	}
	if v > MinutePerKM.Max {
		t.Logf("too large speed %v min/km", v)
		t.Fail()
	}
}

func TestSmartCreateRecordsBefore(t *testing.T) {
	tests := []struct {
		sex string
	}{
		{"F"}, {"M"},
	}
	for _, test := range tests {
		t.Run(test.sex, func(t *testing.T) {
			t.Logf(test.sex)
			for i := 0; i < 1000; i++ {
				t.Run(t.Name(), func(t *testing.T) {
					t.Parallel()
					limit := GetDefaultLimitParams(test.sex)
					timePoint := time.Now()
					records := SmartCreateRecordsBefore(0, 0, limit, limit.LimitTotalMaxDistance, timePoint)
					if len(records) == 0 {
						t.Fail()
					}
					for _, r := range records {
						checkRecordSpeed(t, r, limit.MinutePerKM)
						if r.BeginTime.After(timePoint) || r.EndTime.After(timePoint) {
							t.Logf("invalid time range %v %v", r.BeginTime, r.EndTime)
							t.Fail()
						}
						t.Logf("%+v", r)
					}
				})
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
			timePoint := time.Now()
			limit := GetDefaultLimitParams(test.sex)
			records := SmartCreateRecordsAfter(0, 0, limit, limit.LimitTotalMaxDistance, timePoint)
			if len(records) == 0 {
				t.Fail()
			}
			for _, r := range records {
				checkRecordSpeed(t, r, limit.MinutePerKM)
				if r.BeginTime.Before(timePoint) || r.EndTime.Before(timePoint) {
					t.Logf("invalid time range %v %v", r.BeginTime, r.EndTime)
					t.Fail()
				}
				t.Logf("%+v", r)
			}

		})
	}
}

func Test_generateRandomTimeDuration(t *testing.T) {
	type args struct {
		mpkmLimit Float64Range
		distance  float64
	}
	flimit := GetDefaultLimitParams("F")
	mlimit := GetDefaultLimitParams("M")
	tests := []struct {
		name string
		args args
	}{
		{"f", args{flimit.MinutePerKM, 2.5}},
		{"m", args{mlimit.MinutePerKM, 4.5}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := 0; i < 1000; i++ {
				t.Logf("generateRandomTimeDuration limit %v %v", tt.args.mpkmLimit.Min, tt.args.mpkmLimit.Max)
				got := generateRandomTimeDuration(tt.args.mpkmLimit, tt.args.distance)
				t.Logf("generateRandomTimeDuration() = %v", got)
				v := got.Minutes() / tt.args.distance
				t.Logf("v = %v", v)
				if v < tt.args.mpkmLimit.Min {
					t.Logf("too small speed %v min/km", v)
					t.Fail()
				}
				if v > tt.args.mpkmLimit.Max {
					t.Logf("too large speed %v min/km", v)
					t.Fail()
				}
			}
		})
	}
}
