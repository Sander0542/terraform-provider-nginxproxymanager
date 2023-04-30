package nginxproxymanager

import (
	"context"
	"fmt"
	"net/http"

	"github.com/getsentry/sentry-go"
)

type tracingTransport struct {
	inner http.RoundTripper
	ctx   context.Context
}

func newTracingTransport(ctx context.Context, inner http.RoundTripper) *tracingTransport {
	return &tracingTransport{inner, ctx}
}

func (t *tracingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	span := sentry.StartSpan(t.ctx, "nginxproxymanager.go.http_request")
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
