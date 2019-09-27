package ssmt

import "testing"

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
