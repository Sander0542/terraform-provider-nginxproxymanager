// Copyright (c) Sander Jochems
// SPDX-License-Identifier: MIT

package provider

import (
	"context"
	"github.com/sander0542/nginxproxymanager-go"
	"github.com/sander0542/terraform-provider-nginxproxymanager/internal/provider/models"

	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ ephemeral.EphemeralResource = &UserTokenEphemeralResource{}
var _ ephemeral.EphemeralResourceWithConfigure = &UserTokenEphemeralResource{}

func NewUserTokenEphemeralResource() ephemeral.EphemeralResource {
	return &UserTokenEphemeralResource{}
}

type UserTokenEphemeralResource struct {
	client *nginxproxymanager.APIClient
	auth   context.Context
}

func (r *UserTokenEphemeralResource) Metadata(_ context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_token"
}

func (r *UserTokenEphemeralResource) Schema(ctx context.Context, _ ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Users --- This ephemeral resource can be used to retrieve an ephemeral token for the current user.",
		Attributes: map[string]schema.Attribute{
			"token": schema.StringAttribute{
				Computed:            true,
				Sensitive:           true,
				MarkdownDescription: "The token for the current user.",
			},
		},
	}
}

func (r *UserTokenEphemeralResource) Configure(ctx context.Context, req ephemeral.ConfigureRequest, resp *ephemeral.ConfigureResponse) {
	if data := ephemeralResourceConfigure(ctx, req, resp); data != nil {
		r.client = data.Client
		r.auth = data.Auth
	}
}

func (r *UserTokenEphemeralResource) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	var data *models.UserToken

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	token, ok := r.auth.Value(nginxproxymanager.ContextAccessToken).(string)

	if !ok {
		resp.Diagnostics.AddError("Client Error", "Unable to read token")
	}

	data.Token = types.StringValue(token)

	resp.Diagnostics.Append(resp.Result.Set(ctx, &data)...)
}
