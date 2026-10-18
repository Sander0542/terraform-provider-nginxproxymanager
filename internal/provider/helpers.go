// Copyright (c) Sander Jochems
// SPDX-License-Identifier: MIT

package provider

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func errorWithResponseBody(err error) string {
	var responseError interface {
		Body() []byte
	}
	if errors.As(err, &responseError) {
		if body := strings.TrimSpace(string(responseError.Body())); body != "" {
			return fmt.Sprintf("%s: %s", err, body)
		}
	}

	return err.Error()
}

func resourceConfigure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) *NginxProxyManagerProviderData {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return nil
	}

	providerData, ok := req.ProviderData.(*NginxProxyManagerProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *NginxProxyManagerProviderData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return nil
	}

	return providerData
}

func dataSourceConfigure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) *NginxProxyManagerProviderData {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return nil
	}

	providerData, ok := req.ProviderData.(*NginxProxyManagerProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *NginxProxyManagerProviderData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return nil
	}

	return providerData
}

func ephemeralResourceConfigure(_ context.Context, req ephemeral.ConfigureRequest, resp *ephemeral.ConfigureResponse) *NginxProxyManagerProviderData {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return nil
	}

	providerData, ok := req.ProviderData.(*NginxProxyManagerProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Ephemeral Resource Configure Type",
			fmt.Sprintf("Expected *NginxProxyManagerProviderData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return nil
	}

	return providerData
}
