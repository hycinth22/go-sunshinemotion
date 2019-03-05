package lib

import (
	"testing"
	"time"
)

func TestSmartCreateRecords(t *testing.T) {
	records := SmartCreateRecords(0, 0, &LimitParams{
		RandDistance:        Float64Range{2.6, 4.0},
		LimitSingleDistance: Float64Range{2.0, 4.0},
		LimitTotalDistance:  Float64Range{2.0, 5.0},
		MinuteDuration:      IntRange{11, 20},
	}, 5, time.Now())
	for _, r := range records {
		t.Logf("%+v", r)
	}
}
