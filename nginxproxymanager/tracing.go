package nginxproxymanager

import (
	"errors"
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"net/http"
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

func Sentry(diags diag.Diagnostics) {
	for _, err := range diags.Errors() {
		sentry.CaptureException(errors.New(err.Detail()))
	}
	for _, warn := range diags.Warnings() {
		sentry.CaptureMessage(warn.Detail())
	}
}
