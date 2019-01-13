package sunshinemotion

import "net/http"

type sessionTransport struct {
	device *Device
}

func (t *sessionTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", t.device.UserAgent)
	return http.DefaultTransport.RoundTrip(req)
}
