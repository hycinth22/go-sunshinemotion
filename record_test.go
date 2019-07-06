package ssmt

import (
	"testing"
	"time"
)

func TestSmartCreateDistance(t *testing.T) {
	flimit := GetDefaultLimitParams("F")
	mlimit := GetDefaultLimitParams("M")
	tests := []struct {
		name   string
		limit  LimitParams
		remain float64
	}{
		{"p1 f", flimit, 2 * flimit.RandDistance.Max},
		{"p1 m", mlimit, 2 * mlimit.RandDistance.Max},
		{"p2 f", flimit, flimit.LimitSingleDistance.Min},
		{"p2 m", mlimit, mlimit.LimitSingleDistance.Min},
		{"p2 m-3.0", mlimit, 3},
		{"p2 m-4.0", mlimit, 4},
		{"p2 m-4.5", mlimit, 4.5},
		{"std f", flimit, flimit.LimitTotalMaxDistance - MinDistanceAccurency},
		{"std m", mlimit, mlimit.LimitTotalMaxDistance - MinDistanceAccurency},
	}
	strict := false
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for i := 0; i < 1000; i++ {
				remain := test.remain
				for remain > test.limit.LimitSingleDistance.Min {
					println("remain", remain)
					distance := smartCreateDistance(test.limit, remain)
					t.Log("smartCreateDistance ", distance)
					if distance > remain {
						t.Log("too far distance", distance)
						t.FailNow()
					}
					if distance < test.limit.LimitSingleDistance.Min-EPSILON_Distance || distance >= test.limit.LimitSingleDistance.Max-EPSILON_Distance {
						t.Log("fail", distance)
						t.FailNow()
					}
					if strict && distance < test.limit.RandDistance.Min-EPSILON_Distance || distance >= test.limit.RandDistance.Max-EPSILON_Distance {
						t.Log("doesn't qualified RandDistance limitation", distance)
						t.FailNow()
					}
					remain -= distance
				}
			}
		})
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
			limit := GetDefaultLimitParams(test.sex)
			records := SmartCreateRecordsBefore(0, 0, limit, limit.LimitTotalMaxDistance-0.1, time.Now())
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
			records := SmartCreateRecordsAfter(0, 0, limit, limit.LimitTotalMaxDistance-0.1, time.Now())
			for _, r := range records {
				t.Logf("%+v", r)
			}

		})
	}
}
