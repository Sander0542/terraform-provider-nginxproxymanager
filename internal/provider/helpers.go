// Copyright (c) Sander Jochems
// SPDX-License-Identifier: MIT

package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/sander0542/nginxproxymanager-go"
)

func resourceConfigure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) (context.Context, *nginxproxymanager.APIClient) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return nil, nil
	}

	providerData, ok := req.ProviderData.(*NginxProxyManagerProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *NginxProxyManagerProviderData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return nil, nil
	}

	return providerData.Auth, providerData.Client
}

func dataSourceConfigure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) (context.Context, *nginxproxymanager.APIClient) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return nil, nil
	}

	providerData, ok := req.ProviderData.(*NginxProxyManagerProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *NginxProxyManagerProviderData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return nil, nil
	}

	return providerData.Auth, providerData.Client
}
