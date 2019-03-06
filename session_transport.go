package sunshinemotion

import (
	"net/http"
	"strconv"
)

type sessionTransport struct {
	device *Device
	s      *Session
}

func (t *sessionTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	token := t.s.token
	req.Header.Set("User-Agent", t.device.UserAgent)
	req.Header["UserID"] = []string{strconv.FormatUint(uint64(t.s.token.UserID), 64)}
	if token.UserID != 0 {
		req.Header["TokenID"] = []string{token.TokenID}
	}
	req.Header["crack"] = []string{"0"}
	return http.DefaultTransport.RoundTrip(req)
}
