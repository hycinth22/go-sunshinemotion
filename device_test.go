package sunshinemotion

import (
	"strconv"
	"testing"
)

func TestGenerateDevice(t *testing.T) {
	const testTimes = 100
	for i := 0; i < testTimes; i++ {
		t.Run("Test "+strconv.Itoa(i), func(t *testing.T) {
			device := GenerateDevice()
			t.Log(*device)
			if device == nil {
				t.FailNow()
			}
		})
	}
}
