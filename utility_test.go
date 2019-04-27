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
		n := randRangeFloat(min, max, 3)
		t.Log(n)
		if n < min || n > max {
			t.FailNow()
		}
	}
}

func TestNormalizeDistance(t *testing.T) {
	test := []struct {
		value    float64
		expected float64
	}{
		{4.44444, 4.444},
		{4.9999, 5.000},
		{4.5999, 4.600},
	}
	for _, c := range test {
		actual := NormalizeDistance(c.value)
		t.Log(c.value, c.expected, actual)
		if actual != c.expected {
			t.Fail()
		}
	}

}
