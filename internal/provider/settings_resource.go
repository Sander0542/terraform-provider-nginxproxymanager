// Copyright (c) Sander Jochems
// SPDX-License-Identifier: MIT

package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/sander0542/nginxproxymanager-go"
	"github.com/sander0542/terraform-provider-nginxproxymanager/internal/provider/models"
)

var _ resource.Resource = &SettingsResource{}

func NewSettingsResource() resource.Resource {
	return &SettingsResource{}
}

type SettingsResource struct {
	client *nginxproxymanager.APIClient
	auth   context.Context
}

type updateRequest struct {
	Id      string
	Request *nginxproxymanager.UpdateSettingRequest
}

func (r *SettingsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_settings"
}

func (r *SettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Settings --- This resource can be used to manage settings.",
		Attributes: map[string]schema.Attribute{
			"default_site": schema.SingleNestedAttribute{
				MarkdownDescription: "What to show when Nginx is hit with an unknown Host.",
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"page": schema.StringAttribute{
						MarkdownDescription: "What to show when Nginx is hit with an unknown Host.",
						Required:            true,
						Validators: []validator.String{
							stringvalidator.OneOf("congratulations", "404", "444", "redirect", "html"),
						},
					},
					"redirect": schema.StringAttribute{
						MarkdownDescription: "Redirect to.",
						Optional:            true,
						Computed:            true,
					},
					"html": schema.StringAttribute{
						MarkdownDescription: "HTML Content.",
						Optional:            true,
						Computed:            true,
					},
				},
			},
		},
	}
}

func (r *SettingsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if data := resourceConfigure(ctx, req, resp); data != nil {
		r.client = data.Client
		r.auth = data.Auth
	}
}

func (r *SettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *models.Settings

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	r.updateSettings(ctx, data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	settings, _, err := r.client.SettingsAPI.GetSettings(r.auth).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read settings, got error: %s", err))
		return
	}

	data.Write(ctx, settings, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *models.Settings

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	settings, _, err := r.client.SettingsAPI.GetSettings(r.auth).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read settings, got error: %s", err))
		return
	}

	data.Write(ctx, settings, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *models.Settings

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	r.updateSettings(ctx, data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	settings, _, err := r.client.SettingsAPI.GetSettings(r.auth).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read settings, got error: %s", err))
		return
	}

	data.Write(ctx, settings, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddWarning("Settings not changed", "The settings resource has been removed, but the settings have not been changed.")
}

func (r *SettingsResource) updateSettings(ctx context.Context, data *models.Settings, diags *diag.Diagnostics) {
	requests := map[string]updateRequest{}

	if !data.DefaultSite.IsUnknown() {
		defaultSite, tmpDiags := models.SettingDefaultSiteAs(ctx, data.DefaultSite)
		diags.Append(tmpDiags...)

		requests["default_site"] = updateRequest{
			Id:      "default-site",
			Request: defaultSite.ToRequest(ctx, diags),
		}
	}

	if diags.HasError() {
		return
	}

	for attributeName, request := range requests {
		_, _, err := r.client.SettingsAPI.UpdateSetting(r.auth, request.Id).UpdateSettingRequest(*request.Request).Execute()
		if err != nil {
			diags.AddAttributeError(path.Root(attributeName), "Client Error", fmt.Sprintf("Unable to update setting, got error: %s", err))
		}
	}
}
