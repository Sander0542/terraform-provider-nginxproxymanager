package nginxproxymanager

import (
	"fmt"
	"net/http"

	"github.com/getsentry/sentry-go"
)

type tracingTransport struct {
	inner http.RoundTripper
}

func newTracingTransport(inner http.RoundTripper) *tracingTransport {
	if inner == nil {
		inner = http.DefaultTransport
	}
	return &tracingTransport{inner}
}

func (t *tracingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	span := sentry.StartSpan(req.Context(), "nginxproxymanager.go.http_request")
	req.Header.Set("sentry-trace", span.ToSentryTrace())
	span.Description = fmt.Sprintf("%s %s", req.Method, req.URL.Path)
	span.SetTag("http.method", req.Method)
	span.SetTag("http.scheme", req.URL.Scheme)
	span.SetTag("http.path", req.URL.Path)
	span.SetTag("http.query", req.URL.RawQuery)
	span.SetTag("http.fragment", req.URL.RawFragment)
	defer span.Finish()

	res, err := t.inner.RoundTrip(req.WithContext(span.Context()))
	return res, err
}
