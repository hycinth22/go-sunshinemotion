package sunshinemotion

import (
	"testing"
	"time"
)

func TestToken_ValidFormat(t *testing.T) {
	tests := []struct {
		name  string
		token Token
		valid bool
	}{
		{"empty token",
			Token{},
			false,
		},
		{"valid token",
			Token{
				TokenID:    "tokenID",
				UserID:     1111,
				SchoolID:   60,
				ExpireTime: time.Now().AddDate(1, 0, 0),
			},
			true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.token.ValidFormat() != test.valid {
				t.Log(test)
				t.Fail()
			}
		})
	}
}

func TestToken_Expired(t *testing.T) {
	tests := []struct {
		name  string
		token Token
		valid bool
	}{
		{"expired token",
			Token{
				TokenID:    "tokenID",
				UserID:     1111,
				SchoolID:   60,
				ExpireTime: time.Now().AddDate(-1, 0, 0),
			},
			true,
		},
		{"valid token",
			Token{
				TokenID:    "tokenID",
				UserID:     1111,
				SchoolID:   60,
				ExpireTime: time.Now().AddDate(1, 0, 0),
			},
			false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.token.Expired() != test.valid {
				t.Log(test)
				t.Fail()
			}
		})
	}
}
