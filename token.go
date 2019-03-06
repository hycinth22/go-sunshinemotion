package sunshinemotion

import (
	"errors"
	"time"
)

type Token struct {
	TokenID    string
	ExpireTime time.Time

	// extra information about the token owner
	UserID   uint   // enough to distinguish different users,
	SchoolID uint64 // sometimes needed
}

var (
	ErrTokenInvalid = errors.New("登录状态无效")
	ErrTokenExpired = errors.New("登录已过期")
)

// only check if token is in valid format.
// it does not send any network request to verify.
func (token *Token) ValidFormat() bool {
	return token != nil &&
		len(token.TokenID) > 0
}

// only check if token exceed its ExpireTime.
// it does not send any network request to verify.
func (token *Token) Expired() bool {
	return time.Now().After(token.ExpireTime)
}
