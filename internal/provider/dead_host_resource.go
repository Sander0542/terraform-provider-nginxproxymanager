// Copyright (c) Sander Jochems
// SPDX-License-Identifier: MIT

package provider

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sander0542/nginxproxymanager-go"
	"github.com/sander0542/terraform-provider-nginxproxymanager/internal/provider/models"
	"strconv"
)

var _ resource.Resource = &DeadHostResource{}
var _ resource.ResourceWithImportState = &DeadHostResource{}

func NewDeadHostResource() resource.Resource {
	return &DeadHostResource{}
}

type DeadHostResource struct {
	client *nginxproxymanager.APIClient
	auth   context.Context
}

func (r *DeadHostResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dead_host"
}

func (r *DeadHostResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Hosts --- This resource can be used to manage a dead host.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "The Id of the dead host.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"created_on": schema.StringAttribute{
				MarkdownDescription: "The date and time the dead host was created.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"modified_on": schema.StringAttribute{
				MarkdownDescription: "The date and time the dead host was last modified.",
				Computed:            true,
			},
			"owner_user_id": schema.Int64Attribute{
				MarkdownDescription: "The Id of the user that owns the dead host.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"domain_names": schema.ListAttribute{
				MarkdownDescription: "The domain names associated with the dead host.",
				Required:            true,
				ElementType:         types.StringType,
			},
			"certificate_id": schema.Int64Attribute{
				MarkdownDescription: "The Id of the certificate used by the dead host.",
				Optional:            true,
			},
			"ssl_forced": schema.BoolAttribute{
				MarkdownDescription: "Whether SSL is forced for the dead host.",
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(false),
			},
			"hsts_enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether HSTS is enabled for the dead host.",
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(false),
			},
			"hsts_subdomains": schema.BoolAttribute{
				MarkdownDescription: "Whether HSTS is enabled for subdomains of the dead host.",
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(false),
			},
			"http2_support": schema.BoolAttribute{
				MarkdownDescription: "Whether HTTP/2 is supported for the dead host.",
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(false),
			},
			"advanced_config": schema.StringAttribute{
				MarkdownDescription: "The advanced configuration used by the dead host.",
				Computed:            true,
				Optional:            true,
				Default:             stringdefault.StaticString(""),
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether the dead host is enabled.",
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(true),
			},
			"meta": schema.MapAttribute{
				MarkdownDescription: "The meta data associated with the dead host.",
				ElementType:         types.StringType,
				Computed:            true,
			},
		},
	}
}

func (r *DeadHostResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if auth, client := resourceConfigure(ctx, req, resp); client != nil {
		r.client = client
		r.auth = auth
	}
}

func (r *DeadHostResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *models.DeadHost

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	hostEnabled := data.Enabled.ValueBool()

	request := data.ToCreateRequest(ctx, &resp.Diagnostics)
	deadHost, _, err := r.client.Class404HostsAPI.Create404Host(r.auth).Create404HostRequest(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create dead host, got error: %s", err))
		return
	}

	data.Write(ctx, deadHost, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), data.Id)...)

	err = r.toggleHost(deadHost.GetId(), deadHost.GetEnabled(), hostEnabled)
	if err != nil {
		resp.Diagnostics.AddAttributeError(path.Root("enabled"), "Client Error", fmt.Sprintf("Unable to update dead host, got err: %s", err))
		return
	}

	data.Enabled = types.BoolValue(hostEnabled)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DeadHostResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *models.DeadHost

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	deadHost, _, err := r.client.Class404HostsAPI.GetDeadHost(r.auth, data.Id.ValueInt64()).Execute()
	if err != nil {
		if err.Error() == "404 Not Found" {
			resp.State.RemoveResource(ctx)
			return
		} else {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read dead host, got error: %s", err))
			return
		}
	}

	data.Write(ctx, deadHost, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DeadHostResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *models.DeadHost

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	hostEnabled := data.Enabled.ValueBool()

	request := data.ToUpdateRequest(ctx, &resp.Diagnostics)
	deadHost, _, err := r.client.Class404HostsAPI.UpdateDeadHost(r.auth, data.Id.ValueInt64()).UpdateDeadHostRequest(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update dead host, got error: %s", err))
		return
	}

	data.Write(ctx, deadHost, &resp.Diagnostics)

	err = r.toggleHost(deadHost.GetId(), deadHost.GetEnabled(), hostEnabled)
	if err != nil {
		resp.Diagnostics.AddAttributeError(path.Root("enabled"), "Client Error", fmt.Sprintf("Unable to update dead host, got err: %s", err))
		return
	}

	data.Enabled = types.BoolValue(hostEnabled)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (r *DeadHostResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *models.DeadHost

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	success, _, err := r.client.Class404HostsAPI.DeleteDeadHost(r.auth, data.Id.ValueInt64()).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete dead host, got error: %s", err))
		return
	}

	if !success {
		resp.Diagnostics.AddError("Server Error", "Unable to delete dead host.")
		return
	}
}

func (r *DeadHostResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Could not convert id to number, got error: %s", err))
		return
	}

	deadHost, _, err := r.client.Class404HostsAPI.GetDeadHost(r.auth, id).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read dead host, got error: %s", err))
		return
	}

	diags := resp.State.SetAttribute(ctx, path.Root("id"), types.Int64Value(deadHost.GetId()))
	resp.Diagnostics.Append(diags...)
}

func (r *DeadHostResource) toggleHost(hostId int64, current bool, desired bool) error {
	if desired && !current {
		enableResponse, _, err := r.client.Class404HostsAPI.EnableDeadHost(r.auth, hostId).Execute()
		if err != nil {
			return err
		} else if !enableResponse {
			return errors.New("unable to enable dead host")
		}
	} else if !desired && current {
		disableResponse, _, err := r.client.Class404HostsAPI.DisableDeadHost(r.auth, hostId).Execute()
		if err != nil {
			return err
		} else if !disableResponse {
			return errors.New("unable to disable dead host")
		}
	}

	return nil
}
