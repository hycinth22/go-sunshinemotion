package crypto

import (
	"testing"
)

func TestCalcDBXTcode(t *testing.T) {
	tests := []struct {
		name      string
		userID    int64
		beginTime string
		distance  string
		xtcode    string
	}{
		{"empty case", 0, "", "", "68AA3D4617390EDF63AFE905743789DF"},
		{"normal case 1", 8422, "2019-09-11 10:49:14", "2.992", "0C35DCDD29F069BB98D8D20A0C79BA99"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expect := test.xtcode
			actual := CalcDBXTcode(test.userID, test.beginTime, test.distance)
			if actual != expect {
				t.Errorf(`test case failed. %v xtcode should be %s but %s`, test, expect, actual)
				t.Fail()
			}
		})
	}
}

func TestCalcXTcode(t *testing.T) {
	tests := []struct {
		name      string
		userID    int64
		beginTime string
		distance  string
		xtcode    string
	}{
		{"empty case", 0, "", "", "7E9F5BF6DFB2153ACA648F768CA85EC2"},
		{"normal case 1", 1234, "2006-01-02 15:04:05", "3.410", "30F96005837052B043B2C28A61E05A6C"},
		{"normal case 2", 5678, "2017-02-03 16:05:06", "4.520", "51BC18ADF27E47D9DAD4B707FB91111C"},
		{"normal case 3", 9012, "2028-03-04 17:06:07", "5.630", "CD6402A5D6248237BECE0F9AF1936250"},
		{"normal case 4", 11732, "2019-03-19 14:38:02", "2.111", "3901495440D9F0FE2BCEC94325FF556E"},
		{"normal case 5", 11732, "2019-03-20 14:38:02", "2.333", "A4A2EDAF03911EB500A6A9D277BC5FEC"},
		{"normal case 6", 9593, "2019-04-27 21:31:58", "0.000", "9C66777808294D2DBC3B6EE625BA6EDB"},
		{"normal case 7", 7424, "2019-04-29 18:52:04", "4.900", "AB19DCBE2C7B92D65580DEFDAA817376"},
		{"normal case 8", 6418, "2019-05-16 07:53:33", "4.871", "222B3990F8A96653ED2A33DCE96A6FC7"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expect := test.xtcode
			actual := CalcXTcode(test.userID, test.beginTime, test.distance)
			t.Logf("%+v %s", test, actual)
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
		{"empty string case", "", "", "E629933040B660F4814DF4A144E74A24"},
		{"normal case",
			"",
			"Lorem ipsum dolor sit amet, consectetur adipisicing elit",
			"BB92F5237F15EEAF6D8D0033BB35C88E"},
		{"normal case",
			"Lorem ipsum dolor sit amet, consectetur adipisicing elit",
			"",
			"BB92F5237F15EEAF6D8D0033BB35C88E"},
		{"normal case",
			"Lorem ipsum dolor sit amet, consectetur adipisicing elit",
			"sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
			"0B68F974EE595C3EA522B3B9E7C5DC52"},
		{"normal case",
			"sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
			"Lorem ipsum dolor sit amet, consectetur adipisicing elit",
			"819E16D28A76AEBE5BB7BB56AAAAB903"},
		{"large string case ",
			"Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
			"Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
			"E2CC4D759AA44076B85A925C2FE638C9"},
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
