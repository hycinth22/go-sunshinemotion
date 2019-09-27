package ssmt

import (
	"regexp"
	"testing"
)

func TestGenerateIMEI(t *testing.T) {
	for i := 0; i < 100; i++ {
		imei := GenerateIMEI()
		t.Log("GenerateIMEI ", imei)
		if imei == "" || len(imei) != 15 {
			t.FailNow()
		}
	}
}

func TestRandModel(t *testing.T) {
	for i := 0; i < 100; i++ {
		model := RandModel()
		t.Log("RandModel ", model)
		if model == "" || len(model) < 2 || len(model) > 7 {
			t.FailNow()
		}
	}
}

func TestRandScreen(t *testing.T) {
	re := regexp.MustCompile(`\d+x\d+`)
	for i := 0; i < 100; i++ {
		screen := RandScreen()
		t.Log("RandScreen ", screen)
		if screen == "" || len(screen) < 2 || !re.MatchString(screen) {
			t.FailNow()
		}
	}
}

func TestRandRangeFloat(t *testing.T) {
	for i := 0; i < 100; i++ {
		min := -1231231.1
		max := 32221232132.20123
		n := randRangeFloat(min, max)
		t.Log(n)
		if n < min || n > max {
			t.FailNow()
		}
	}
}
