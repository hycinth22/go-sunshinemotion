package lib

import (
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
		if model == "" || len(model) < 2 || len(model) > 4 {
			t.FailNow()
		}
	}
}
