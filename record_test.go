package ssmt

import (
	"fmt"
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

func Test_smartCreateDistance(t *testing.T) {
	type args struct {
		limitParams LimitParams
		remain      float64
	}
	flimit := GetDefaultLimitParams("F")
	mlimit := GetDefaultLimitParams("M")
	flimitReduce := LimitParams{
		RandDistance:          Float64Range{flimit.LimitSingleDistance.Min, flimit.LimitSingleDistance.Min + (flimit.RandDistance.Max-flimit.LimitSingleDistance.Min)*0.87},
		LimitSingleDistance:   flimit.LimitSingleDistance,
		LimitTotalMaxDistance: flimit.LimitTotalMaxDistance,
		MinutePerKM:           flimit.MinutePerKM,
	}
	mlimitReduce := LimitParams{
		RandDistance:          Float64Range{mlimit.LimitSingleDistance.Min, mlimit.LimitSingleDistance.Min + (mlimit.RandDistance.Max-mlimit.LimitSingleDistance.Min)*0.87},
		LimitSingleDistance:   mlimit.LimitSingleDistance,
		LimitTotalMaxDistance: mlimit.LimitTotalMaxDistance,
		MinutePerKM:           mlimit.MinutePerKM,
	}
	tests := []struct {
		name   string
		args   args
		inRand bool // strict mode, must be meet the randDistance limitation
	}{
		{"p1 f", args{flimit, 2 * flimit.RandDistance.Max}, true},
		{"p1 m", args{mlimit, 2 * mlimit.RandDistance.Max}, true},
		{"p2 f", args{flimit, flimit.LimitSingleDistance.Min}, true},
		{"p2 m", args{mlimit, mlimit.LimitSingleDistance.Min}, true},
		{"p2 m-3.0", args{mlimit, 3}, false},
		{"p2 m-4.0", args{mlimit, 4}, false},
		{"p2 m-4.5", args{mlimit, 4.5}, false},
		{"std f", args{flimit, flimit.LimitTotalMaxDistance}, false},
		{"std m", args{mlimit, mlimit.LimitTotalMaxDistance}, false},
		{"reduce f", args{flimitReduce, flimit.LimitTotalMaxDistance}, true},
		{"reduce m", args{mlimitReduce, mlimit.LimitTotalMaxDistance}, true},
	}
	validate := func(arg args, singleDistance float64, remain float64, inRand bool) error {
		if singleDistance > remain {
			return fmt.Errorf("singleDistance(%v) exceed remain. remain = %v", singleDistance, remain)
		}
		if !arg.limitParams.LimitSingleDistance.In(singleDistance, EpsilonDistance) {
			return fmt.Errorf("singleDistance(%v) unqualified for LimitSingleDistance limitation [%v, %v)", singleDistance, arg.limitParams.LimitSingleDistance.Min, arg.limitParams.LimitSingleDistance.Max)
		}
		if inRand && !arg.limitParams.RandDistance.In(singleDistance, EpsilonDistance) {
			return fmt.Errorf("singleDistance(%v) unqualified for RandDistance limitation [%v, %v)", singleDistance, arg.limitParams.RandDistance.Min, arg.limitParams.RandDistance.Max)
		}
		return nil
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := 0; i < 1000; i++ {
				remain := tt.args.remain
				t.Run(fmt.Sprintf("%s_%d", tt.name, i), func(t *testing.T) {
					for remain > tt.args.limitParams.LimitSingleDistance.Min {
						t.Log("remain", remain)
						gotSingleDistance := smartCreateDistance(tt.args.limitParams, remain)
						t.Log("smartCreateDistance() = ", gotSingleDistance)
						if err := validate(tt.args, gotSingleDistance, remain, tt.inRand); err != nil {
							t.Errorf("error: %v. \nargs: %+v", err, tt.args)
							t.FailNow()
						}
						remain -= gotSingleDistance
					}
				})
			}
		})
	}
}
