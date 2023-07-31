package sentry

import (
	"context"
	"errors"
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func CaptureDiagnostics(diags diag.Diagnostics) {
	for _, err := range diags.Errors() {
		sentry.CaptureException(errors.New(err.Detail()))
	}
	for _, warn := range diags.Warnings() {
		sentry.CaptureMessage(warn.Detail())
	}
}

func StartResource(ctx context.Context, operation string, name string) *sentry.Span {
	return sentry.StartSpan(ctx, fmt.Sprintf("terraform.resource.%s", operation), sentry.WithTransactionName(fmt.Sprintf("resource.%s.%s", name, operation)))
}

func StartDataSource(ctx context.Context, operation string, name string) *sentry.Span {
	return sentry.StartSpan(ctx, fmt.Sprintf("terraform.data_source.%s", operation), sentry.WithTransactionName(fmt.Sprintf("data.%s.%s", name, operation)))
}
