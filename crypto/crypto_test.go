package crypto

import (
	"testing"
)

func TestCalcXTcode(t *testing.T) {
	tests := []struct {
		name      string
		userID    int64
		beginTime string
		distance  string
		xtcode    string
	}{
		{"empty case", 0, "", "", "41B79A796B6B16A7701F70684D0A3A6E"},
		{"normal case 1", 1234, "2006-01-02 15:04:05", "3.410", "ED9335F08B363EF8C9954302EE84547F"},
		{"normal case 2", 5678, "2017-02-03 16:05:06", "4.520", "38A2FDC3A89C6200CE41BF0639F145ED"},
		{"normal case 3", 9012, "2028-03-04 17:06:07", "5.630", "1630BC477C9BB0150773434F9A07C0DD"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expect := test.xtcode
			actual := CalcXTcode(test.userID, test.beginTime, test.distance)
			if actual != expect {
				t.Errorf(`test case failed. %v xtcode should be %s but %s`, test, expect, actual)
				t.Fail()
			}
		})
	}
}

func TestCalcLi(t *testing.T) {
	tests := []struct {
		name string
		p0   string
		p1   string
		li   string
	}{
		{"empty string case", "", "", "C383EF76A3AE3C080066F1D11DA9104B"},
		{"normal case",
			"",
			"Lorem ipsum dolor sit amet, consectetur adipisicing elit",
			"39B3827C087FBB73A18074493038FCB3"},
		{"normal case",
			"Lorem ipsum dolor sit amet, consectetur adipisicing elit",
			"",
			"39B3827C087FBB73A18074493038FCB3"},
		{"normal case",
			"Lorem ipsum dolor sit amet, consectetur adipisicing elit",
			"sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
			"4B9E5FAC650CA0F00C0AAF1B2E477C61"},
		{"normal case",
			"sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
			"Lorem ipsum dolor sit amet, consectetur adipisicing elit",
			"A2BAA4DCD5791013D1A92CA371479D52"},
		{"large string case ",
			"Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
			"Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
			"0707B9E88FB988E68F3C4656F99E52F4"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expect := test.li
			actual := CalcLi(test.p0, test.p1)
			if actual != expect {
				t.Errorf(`test case failed. %v li should be %s but %s`, test, expect, actual)
				t.Fail()
			}
		})
	}
}