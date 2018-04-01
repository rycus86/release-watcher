package transport

import "net/http"

const DefaultUserAgent = "Release Watcher (https://github.com/rycus86/release-watcher)"

type HttpTransportWithUserAgent struct {
	UserAgent string
}

func (t *HttpTransportWithUserAgent) RoundTrip(request *http.Request) (*http.Response, error) {
	if t.UserAgent != "" {
		request.Header.Set("User-Agent", t.UserAgent)
	} else {
		request.Header.Set("User-Agent", DefaultUserAgent)
	}

	return http.DefaultTransport.RoundTrip(request)
}
