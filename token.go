package sunshinemotion

import (
	"errors"
	"time"
)

// Use only UserID enough to distinguish different users,
// but SchoolID is sometimes needed as a extra information, this's why it's in Token struct.
type Token struct {
	TokenID    string
	UserID     uint
	SchoolID   uint64
	ExpireTime time.Time
}

var ErrTokenExpired = errors.New("token has expired")

// only check if token is valid in format.
// it does not send any network request to verify.
func (token *Token) Valid() bool {
	return token != nil &&
		len(token.TokenID) > 0 &&
		time.Now().Before(token.ExpireTime)
}
