package sunshinemotion

import (
	"testing"
	"time"
)

func TestToken_Valid(t *testing.T) {
	tests := []struct {
		name  string
		token Token
		valid bool
	}{
		{"empty token",
			Token{},
			false,
		},
		{"expired token",
			Token{
				TokenID:    "tokenID",
				UserID:     1111,
				SchoolID:   60,
				ExpireTime: time.Now().AddDate(-1, 0, 0),
			},
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
			if test.token.Valid() != test.valid {
				t.Fail()
			}
		})
	}
}
