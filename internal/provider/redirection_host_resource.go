// Copyright (c) Sander Jochems
// SPDX-License-Identifier: MIT

package provider

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sander0542/nginxproxymanager-go"
	"github.com/sander0542/terraform-provider-nginxproxymanager/internal/provider/models"
	"strconv"
)

var _ resource.Resource = &RedirectionHostResource{}
var _ resource.ResourceWithImportState = &RedirectionHostResource{}

func NewRedirectionHostResource() resource.Resource {
	return &RedirectionHostResource{}
}

type RedirectionHostResource struct {
	client *nginxproxymanager.APIClient
	auth   context.Context
}

func (r *RedirectionHostResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_redirection_host"
}

func (r *RedirectionHostResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Hosts --- This resource can be used to manage a redirection host.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "The Id of the redirection host.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"created_on": schema.StringAttribute{
				MarkdownDescription: "The date and time the redirection host was created.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"modified_on": schema.StringAttribute{
				MarkdownDescription: "The date and time the redirection host was last modified.",
				Computed:            true,
			},
			"owner_user_id": schema.Int64Attribute{
				MarkdownDescription: "The Id of the user that owns the redirection host.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"domain_names": schema.ListAttribute{
				MarkdownDescription: "The domain names associated with the redirection host.",
				Required:            true,
				ElementType:         types.StringType,
			},
			"forward_http_code": schema.Int64Attribute{
				MarkdownDescription: "The HTTP code used to forward requests to the redirection host. Must be one of `300`, `301`, `302`, `303`, `307` or `308`",
				Computed:            true,
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.OneOf(300, 301, 302, 303, 307, 308),
				},
				Default: int64default.StaticInt64(300),
			},
			"forward_scheme": schema.StringAttribute{
				MarkdownDescription: "The scheme used to forward requests to the redirection host. Must be one of `auto`, `http` or `https`.",
				Computed:            true,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("auto", "http", "https"),
				},
				Default: stringdefault.StaticString("auto"),
			},
			"forward_domain_name": schema.StringAttribute{
				MarkdownDescription: "The domain name used to forward requests to the redirection host.",
				Required:            true,
			},
			"preserve_path": schema.BoolAttribute{
				MarkdownDescription: "Whether the path is preserved for the redirection host.",
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(true),
			},
			"block_exploits": schema.BoolAttribute{
				MarkdownDescription: "Whether exploits are blocked for the redirection host.",
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(false),
			},
			"certificate_id": schema.Int64Attribute{
				MarkdownDescription: "The Id of the certificate used by the redirection host.",
				Optional:            true,
			},
			"ssl_forced": schema.BoolAttribute{
				MarkdownDescription: "Whether SSL is forced for the redirection host.",
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(false),
			},
			"hsts_enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether HSTS is enabled for the redirection host.",
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(false),
			},
			"hsts_subdomains": schema.BoolAttribute{
				MarkdownDescription: "Whether HSTS is enabled for subdomains of the redirection host.",
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(false),
			},
			"http2_support": schema.BoolAttribute{
				MarkdownDescription: "Whether HTTP/2 is supported for the redirection host.",
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(false),
			},
			"advanced_config": schema.StringAttribute{
				MarkdownDescription: "The advanced configuration used by the redirection host.",
				Computed:            true,
				Optional:            true,
				Default:             stringdefault.StaticString(""),
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether the redirection host is enabled.",
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(true),
			},
			"meta": schema.MapAttribute{
				MarkdownDescription: "The meta data associated with the redirection host.",
				ElementType:         types.StringType,
				Computed:            true,
			},
		},
	}
}

func (r *RedirectionHostResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if auth, client := resourceConfigure(ctx, req, resp); client != nil {
		r.client = client
		r.auth = auth
	}
}

func (r *RedirectionHostResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *models.RedirectionHost

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	hostEnabled := data.Enabled.ValueBool()

	request := data.ToCreateRequest(ctx, &resp.Diagnostics)
	redirectionHost, _, err := r.client.RedirectionHostsAPI.CreateRedirectionHost(r.auth).CreateRedirectionHostRequest(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create redirection host, got error: %s", err))
		return
	}

	data.Write(ctx, redirectionHost, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), data.Id)...)

	err = r.toggleHost(redirectionHost.GetId(), redirectionHost.GetEnabled(), hostEnabled)
	if err != nil {
		resp.Diagnostics.AddAttributeError(path.Root("enabled"), "Client Error", fmt.Sprintf("Unable to update redirection host, got err: %s", err))
		return
	}

	data.Enabled = types.BoolValue(hostEnabled)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RedirectionHostResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *models.RedirectionHost

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	redirectionHost, _, err := r.client.RedirectionHostsAPI.GetRedirectionHost(r.auth, data.Id.ValueInt64()).Execute()
	if err != nil {
		if err.Error() == "404 Not Found" {
			resp.State.RemoveResource(ctx)
			return
		} else {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read redirection host, got error: %s", err))
			return
		}
	}

	data.Write(ctx, redirectionHost, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RedirectionHostResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *models.RedirectionHost

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	hostEnabled := data.Enabled.ValueBool()

	request := data.ToUpdateRequest(ctx, &resp.Diagnostics)
	redirectionHost, _, err := r.client.RedirectionHostsAPI.UpdateRedirectionHost(r.auth, data.Id.ValueInt64()).UpdateRedirectionHostRequest(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update redirection host, got error: %s", err))
		return
	}

	data.Write(ctx, redirectionHost, &resp.Diagnostics)

	err = r.toggleHost(redirectionHost.GetId(), redirectionHost.GetEnabled(), hostEnabled)
	if err != nil {
		resp.Diagnostics.AddAttributeError(path.Root("enabled"), "Client Error", fmt.Sprintf("Unable to update redirection host, got err: %s", err))
		return
	}

	data.Enabled = types.BoolValue(hostEnabled)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (r *RedirectionHostResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *models.RedirectionHost

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	success, _, err := r.client.RedirectionHostsAPI.DeleteRedirectionHost(r.auth, data.Id.ValueInt64()).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete redirection host, got error: %s", err))
		return
	}

	if !success {
		resp.Diagnostics.AddError("Server Error", "Unable to delete redirection host.")
		return
	}
}

func (r *RedirectionHostResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Could not convert id to number, got error: %s", err))
		return
	}

	redirectionHost, _, err := r.client.RedirectionHostsAPI.GetRedirectionHost(r.auth, id).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read redirection host, got error: %s", err))
		return
	}

	diags := resp.State.SetAttribute(ctx, path.Root("id"), types.Int64Value(redirectionHost.GetId()))
	resp.Diagnostics.Append(diags...)
}

func (r *RedirectionHostResource) toggleHost(hostId int64, current bool, desired bool) error {
	if desired && !current {
		enableResponse, _, err := r.client.RedirectionHostsAPI.EnableRedirectionHost(r.auth, hostId).Execute()
		if err != nil {
			return err
		} else if !enableResponse {
			return errors.New("unable to enable redirection host")
		}
	} else if !desired && current {
		disableResponse, _, err := r.client.RedirectionHostsAPI.DisableRedirectionHost(r.auth, hostId).Execute()
		if err != nil {
			return err
		} else if !disableResponse {
			return errors.New("unable to disable redirection host")
		}
	}

	return nil
}
