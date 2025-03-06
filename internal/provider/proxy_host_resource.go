// Copyright (c) Sander Jochems
// SPDX-License-Identifier: MIT

package provider

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sander0542/nginxproxymanager-go"
	"github.com/sander0542/terraform-provider-nginxproxymanager/internal/provider/models"
	"strconv"
)

var _ resource.Resource = &ProxyHostResource{}
var _ resource.ResourceWithImportState = &ProxyHostResource{}

func NewProxyHostResource() resource.Resource {
	return &ProxyHostResource{}
}

type ProxyHostResource struct {
	client *nginxproxymanager.APIClient
	auth   context.Context
}

func (r *ProxyHostResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_proxy_host"
}

func (r *ProxyHostResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Hosts --- This resource can be used to manage a proxy host.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "The Id of the proxy host.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"created_on": schema.StringAttribute{
				MarkdownDescription: "The date and time the proxy host was created.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"modified_on": schema.StringAttribute{
				MarkdownDescription: "The date and time the proxy host was last modified.",
				Computed:            true,
			},
			"owner_user_id": schema.Int64Attribute{
				MarkdownDescription: "The Id of the user that owns the proxy host.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"domain_names": schema.SetAttribute{
				MarkdownDescription: "The domain names associated with the proxy host.",
				Required:            true,
				ElementType:         types.StringType,
			},
			"forward_scheme": schema.StringAttribute{
				MarkdownDescription: "The scheme used to forward requests to the proxy host. Can be either `http` or `https`.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("http", "https"),
				},
			},
			"forward_host": schema.StringAttribute{
				MarkdownDescription: "The host used to forward requests to the proxy host.",
				Required:            true,
			},
			"forward_port": schema.Int64Attribute{
				MarkdownDescription: "The port used to forward requests to the proxy host.",
				Required:            true,
				Validators: []validator.Int64{
					int64validator.Between(1, 65535),
				},
			},
			"certificate_id": schema.Int64Attribute{
				MarkdownDescription: "The Id of the certificate used by the proxy host.",
				Optional:            true,
			},
			"ssl_forced": schema.BoolAttribute{
				MarkdownDescription: "Whether SSL is forced for the proxy host.",
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(false),
			},
			"hsts_enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether HSTS is enabled for the proxy host.",
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(false),
			},
			"hsts_subdomains": schema.BoolAttribute{
				MarkdownDescription: "Whether HSTS is enabled for subdomains of the proxy host.",
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(false),
			},
			"http2_support": schema.BoolAttribute{
				MarkdownDescription: "Whether HTTP/2 is supported for the proxy host.",
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(false),
			},
			"block_exploits": schema.BoolAttribute{
				MarkdownDescription: "Whether exploits are blocked for the proxy host.",
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(false),
			},
			"caching_enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether caching is enabled for the proxy host.",
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(false),
			},
			"allow_websocket_upgrade": schema.BoolAttribute{
				MarkdownDescription: "Whether websocket upgrades are allowed for the proxy host.",
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(false),
			},
			"access_list_id": schema.Int64Attribute{
				MarkdownDescription: "The Id of the access list used by the proxy host.",
				Optional:            true,
			},
			"advanced_config": schema.StringAttribute{
				MarkdownDescription: "The advanced configuration used by the proxy host.",
				Computed:            true,
				Optional:            true,
				Default:             stringdefault.StaticString(""),
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether the proxy host is enabled.",
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(true),
			},
			"meta": schema.MapAttribute{
				MarkdownDescription: "The meta data associated with the proxy host.",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"locations": schema.SetNestedAttribute{
				MarkdownDescription: "The locations associated with the proxy host.",
				Optional:            true,
				Computed:            true,
				Default:             setdefault.StaticValue(types.SetValueMust(models.ProxyHostLocation{}.GetType(), []attr.Value{})),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"path": schema.StringAttribute{
							MarkdownDescription: "The path associated with the location.",
							Required:            true,
						},
						"forward_scheme": schema.StringAttribute{
							MarkdownDescription: "The scheme used to forward requests to the location. Can be either `http` or `https`.",
							Required:            true,
							Validators: []validator.String{
								stringvalidator.OneOf("http", "https"),
							},
						},
						"forward_host": schema.StringAttribute{
							MarkdownDescription: "The host used to forward requests to the location.",
							Required:            true,
						},
						"forward_port": schema.Int64Attribute{
							MarkdownDescription: "The port used to forward requests to the location.",
							Required:            true,
							Validators: []validator.Int64{
								int64validator.Between(1, 65535),
							},
						},
						"advanced_config": schema.StringAttribute{
							MarkdownDescription: "The advanced configuration used by the location.",
							Computed:            true,
							Optional:            true,
							Default:             stringdefault.StaticString(""),
						},
					},
				},
			},
		},
	}
}

func (r *ProxyHostResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if auth, client := resourceConfigure(ctx, req, resp); client != nil {
		r.client = client
		r.auth = auth
	}
}

func (r *ProxyHostResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *models.ProxyHost

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	hostEnabled := data.Enabled.ValueBool()

	request := data.ToCreateRequest(ctx, &resp.Diagnostics)
	proxyHost, _, err := r.client.ProxyHostsAPI.CreateProxyHost(r.auth).CreateProxyHostRequest(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create proxy host, got error: %s", err))
		return
	}

	data.Write(ctx, proxyHost, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), data.Id)...)

	err = r.toggleHost(proxyHost.GetId(), proxyHost.GetEnabled(), hostEnabled)
	if err != nil {
		resp.Diagnostics.AddAttributeError(path.Root("enabled"), "Client Error", fmt.Sprintf("Unable to update proxy host, got err: %s", err))
		return
	}

	data.Enabled = types.BoolValue(hostEnabled)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProxyHostResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *models.ProxyHost

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	proxyHost, _, err := r.client.ProxyHostsAPI.GetProxyHost(r.auth, data.Id.ValueInt64()).Execute()
	if err != nil {
		if err.Error() == "404 Not Found" {
			resp.State.RemoveResource(ctx)
			return
		} else {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read proxy host, got error: %s", err))
			return
		}
	}

	data.Write(ctx, proxyHost, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProxyHostResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *models.ProxyHost

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	hostEnabled := data.Enabled.ValueBool()

	request := data.ToUpdateRequest(ctx, &resp.Diagnostics)
	proxyHost, _, err := r.client.ProxyHostsAPI.UpdateProxyHost(r.auth, data.Id.ValueInt64()).UpdateProxyHostRequest(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update proxy host, got error: %s", err))
		return
	}

	data.Write(ctx, proxyHost, &resp.Diagnostics)

	err = r.toggleHost(proxyHost.GetId(), proxyHost.GetEnabled(), hostEnabled)
	if err != nil {
		resp.Diagnostics.AddAttributeError(path.Root("enabled"), "Client Error", fmt.Sprintf("Unable to update proxy host, got err: %s", err))
		return
	}

	data.Enabled = types.BoolValue(hostEnabled)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (r *ProxyHostResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *models.ProxyHost

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	success, _, err := r.client.ProxyHostsAPI.DeleteProxyHost(r.auth, data.Id.ValueInt64()).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete proxy host, got error: %s", err))
		return
	}

	if !success {
		resp.Diagnostics.AddError("Server Error", "Unable to delete proxy host.")
		return
	}
}

func (r *ProxyHostResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Could not convert id to number, got error: %s", err))
		return
	}

	proxyHost, _, err := r.client.ProxyHostsAPI.GetProxyHost(r.auth, id).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read proxy host, got error: %s", err))
		return
	}

	diags := resp.State.SetAttribute(ctx, path.Root("id"), types.Int64Value(proxyHost.GetId()))
	resp.Diagnostics.Append(diags...)
}

func (r *ProxyHostResource) toggleHost(hostId int64, current bool, desired bool) error {
	if desired && !current {
		enableResponse, _, err := r.client.ProxyHostsAPI.EnableProxyHost(r.auth, hostId).Execute()
		if err != nil {
			return err
		} else if !enableResponse {
			return errors.New("unable to enable proxy host")
		}
	} else if !desired && current {
		disableResponse, _, err := r.client.ProxyHostsAPI.DisableProxyHost(r.auth, hostId).Execute()
		if err != nil {
			return err
		} else if !disableResponse {
			return errors.New("unable to disable proxy host")
		}
	}

	return nil
}
