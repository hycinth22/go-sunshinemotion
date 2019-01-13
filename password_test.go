package sunshinemotion

import "testing"

func TestPasswordHash(t *testing.T) {
	tests := []struct {
		name     string
		password string
		hash     string
	}{
		{"Empty String", "", "d41d8cd98f00b204e9800998ecf8427e"},
		{"Digit String", "123456", "e10adc3949ba59abbe56e057f20f883e"},
		{"Alpha-Symbol String", "Message-Digest", "703a404c0e706b05e970cc3b1d137cb7"},
		{
			"Large String",
			"Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
			"fa5c89f3c88b81bfd5e821b0316569af"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := PasswordHash(test.password)
			expect := test.hash
			if actual != expect {
				t.Fail()
				t.Logf(`"%s"(length %d) hash should be %s but %s`, expect, len(expect), expect, actual)
			}
		})
	}
}
