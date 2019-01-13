package crypto

import (
	"testing"
)

func TestEncodeBZ(t *testing.T) {
	tests := []struct {
		name   string
		bz     string
		secret string
	}{
		{"empty string", "", "3B6AD1E96C13F8C893CA9A649313DD13"},
		{"empty set", "[]", "CEC0CC216B9DB5288E5A9ED306BA2525"},
		{"normal set",
			"[1544345190, Android,24,7.0, 865964032623895, 460110101098930]",
			"4F9F4360340CB760BDA9751E5BC14BFB4E3FB012E80F310AF964BACA1C908E35450C22423DA0912F4F1C41188C88E96B552623AAE02C0860A08E01B6273E6447E976041C761485D22633BE9335C5E6951CC2B0FA5A7A02E49668B802656E957A",
		},
		{"cheated set",
			"[1544338246, xposed, root, Android,19,4.4.2, 354730010542250, 460075422517361]",
			"9456B81E71B9A52E39AECEF0AFD6A871EA8D9726E1121BC4633635D39958F6B31483F1DA46D1FF9314004637866754AD50AD1FB2911B66DB28B35306C208E683A5D4D250B1EFEE26409E35F4206AB9E067475254294B75F05171A095C3063647F679410A9EB6074D6495BB40EC634ED1",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expect := test.secret
			actual := EncryptBZ(test.bz)
			t.Log(test.bz, actual, len(actual))
			if actual != expect {
				t.Errorf(`%s's result should be "%s" but "%s"`, test.bz, expect, actual)
				t.Fail()
			}
		})
	}
}

func TestDecodeBZ(t *testing.T) {
	tests := []struct {
		name   string
		bz     string
		secret string
	}{
		{"empty string", "", "3B6AD1E96C13F8C893CA9A649313DD13"},
		{"empty set", "[]", "CEC0CC216B9DB5288E5A9ED306BA2525"},
		{"normal set",
			"[1544345190, Android,24,7.0, 865964032623895, 460110101098930]",
			"4F9F4360340CB760BDA9751E5BC14BFB4E3FB012E80F310AF964BACA1C908E35450C22423DA0912F4F1C41188C88E96B552623AAE02C0860A08E01B6273E6447E976041C761485D22633BE9335C5E6951CC2B0FA5A7A02E49668B802656E957A",
		},
		{"cheated set",
			"[1544338246, xposed, root, Android,19,4.4.2, 354730010542250, 460075422517361]",
			"9456B81E71B9A52E39AECEF0AFD6A871EA8D9726E1121BC4633635D39958F6B31483F1DA46D1FF9314004637866754AD50AD1FB2911B66DB28B35306C208E683A5D4D250B1EFEE26409E35F4206AB9E067475254294B75F05171A095C3063647F679410A9EB6074D6495BB40EC634ED1",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expect := test.bz
			actual, err := DecryptBZ(test.secret)
			if err != nil {
				t.Error(err)
				t.FailNow()
			}
			t.Log(test.secret, actual, len(actual))
			if actual != expect {
				t.Errorf(`%s's result should be "%s" but "%s"`, test.secret, expect, actual)
				t.Fail()
			}
		})
	}
}

func TestCalcXTcode(t *testing.T) {
	tests := []struct {
		name      string
		userID    uint
		beginTime string
		distance  string
		xtcode    string
	}{
		{"empty case", 0, "", "", "41B79A796B6B16A7701F70684D0A3A6E"},
		{"normal case 1", 1234, "2006-01-02 15:04:05", "3.41", "EB7895DF2007181B67C734AA4EE63998"},
		{"normal case 2", 5678, "2017-02-03 16:05:06", "4.52", "757F0B91B1117CF126B9A95B73A5EAB8"},
		{"normal case 3", 9012, "2028-03-04 17:06:07", "5.63", "4686B76C3E5D1255A346CAFD60439DB5"},
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
