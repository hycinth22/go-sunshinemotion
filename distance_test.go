package ssmt

import (
	"fmt"
	"testing"
)

func TestNormalizeDistance(t *testing.T) {
	type args struct {
		distance float64
	}
	tests := []struct {
		name                  string
		args                  args
		wantNormalizeDistance float64
	}{
		{"", args{4.44444}, 4.444},
		{"", args{4.9999}, 5.000},
		{"", args{4.5999}, 4.600},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotNormalizeDistance := NormalizeDistance(tt.args.distance); gotNormalizeDistance != tt.wantNormalizeDistance {
				t.Errorf("NormalizeDistance() = %v, want %v", gotNormalizeDistance, tt.wantNormalizeDistance)
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

func Test_randRangeDistance(t *testing.T) {
	type args struct {
		min float64
		max float64
	}
	tests := []struct {
		name string
		args args
	}{
		{"", args{4.995, 5}},
		{"", args{2.995, 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := 0; i < 1000; i++ {
				t.Run(fmt.Sprintf("%s_%d", tt.name, i), func(t *testing.T) {
					t.Parallel()
					got := randRangeDistance(tt.args.min, tt.args.max)
					roundGot, _ := fromServiceStdDistance(toServiceStdDistance(got))
					if roundGot < tt.args.min || roundGot >= tt.args.max {
						t.Errorf("randRangeDistance() = %v, roundGot = %v, want [%.3f, %.3f)", got, roundGot, tt.args.min, tt.args.max)
						t.Fail()
					}
				})
			}
		})
	}
}

func TestDistanceRangeAround(t *testing.T) {
	type args struct {
		x     float64
		width int
	}
	tests := []struct {
		name    string
		args    args
		wantMin float64
		wantMax float64
	}{
		{"", args{x: 1.000124521313423, width: 2}, 0.995, 1.005},
		{"", args{x: 3.000124521313423, width: 2}, 2.995, 3.005},
		{"", args{x: 5.000124521313423, width: 2}, 4.995, 5.005},
		{"", args{x: 1.000124521313423, width: 2}, 0.995, 1.005},
		{"", args{x: 3.4000124521313423, width: 2}, 3.395, 3.405},
		{"", args{x: 3.5000124521313423, width: 2}, 3.495, 3.505},
		{"", args{x: 3.6000124521313423, width: 2}, 3.595, 3.605},
		{"", args{x: 4.5000124521313423, width: 2}, 4.495, 4.505},
		{"", args{x: 100.0000124521313423, width: 2}, 99.995, 100.005},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMin, gotMax := DistanceRangeAround(tt.args.x, tt.args.width)
			if gotMin != tt.wantMin {
				t.Errorf("DistanceRangeAround() gotMin = %v, want %v", gotMin, tt.wantMin)
			}
			if gotMax != tt.wantMax {
				t.Errorf("DistanceRangeAround() gotMax = %v, want %v", gotMax, tt.wantMax)
			}
		})
	}
}
