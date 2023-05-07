package nginxproxymanager

import (
	"context"
	models2 "github.com/sander0542/terraform-provider-nginxproxymanager/client/models"
	attributes "github.com/sander0542/terraform-provider-nginxproxymanager/nginxproxymanager/attributes/resources"
	"github.com/sander0542/terraform-provider-nginxproxymanager/nginxproxymanager/models"
	"strconv"

	"github.com/getsentry/sentry-go"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sander0542/terraform-provider-nginxproxymanager/client"
	"github.com/sander0542/terraform-provider-nginxproxymanager/nginxproxymanager/common"
)

var (
	_ common.IResource                    = &proxyHostResource{}
	_ resource.ResourceWithConfigure      = &proxyHostResource{}
	_ resource.ResourceWithValidateConfig = &proxyHostResource{}
	_ resource.ResourceWithImportState    = &proxyHostResource{}
)

func NewProxyHostResource() resource.Resource {
	b := &common.Resource{Name: "proxy_host"}
	r := &proxyHostResource{b, nil}
	b.IResource = r
	return r
}

type proxyHostResource struct {
	*common.Resource
	client *client.Client
}

func (r *proxyHostResource) SchemaImpl(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage a proxy host.",
		Attributes:  attributes.ProxyHost,
		Blocks: map[string]schema.Block{
			"location": schema.ListNestedBlock{
				Description: "The locations associated with the proxy host.",
				NestedObject: schema.NestedBlockObject{
					Attributes: attributes.ProxyHostLocation,
				},
			},
		},
	}
}

func (r *proxyHostResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*client.Client)
}

func (r *proxyHostResource) CreateImpl(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan models.ProxyHost
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var item models2.ProxyHostCreate

	item.DomainNames = []string{}
	for _, domainName := range plan.DomainNames {
		item.DomainNames = append(item.DomainNames, domainName.ValueString())
	}
	item.ForwardScheme = plan.ForwardScheme.ValueString()
	item.ForwardHost = plan.ForwardHost.ValueString()
	item.ForwardPort = uint16(plan.ForwardPort.ValueInt64())
	item.CertificateID = plan.CertificateID.ValueString()
	item.SSLForced = plan.SSLForced.ValueBool()
	item.HSTSEnabled = plan.HSTSEnabled.ValueBool()
	item.HSTSSubdomains = plan.HSTSSubdomains.ValueBool()
	item.HTTP2Support = plan.HTTP2Support.ValueBool()
	item.BlockExploits = plan.BlockExploits.ValueBool()
	item.CachingEnabled = plan.CachingEnabled.ValueBool()
	item.AllowWebsocketUpgrade = plan.AllowWebsocketUpgrade.ValueBool()
	item.AccessListID = plan.AccessListID.ValueInt64()
	item.AdvancedConfig = plan.AdvancedConfig.ValueString()
	item.Meta = map[string]string{}
	item.Locations = []models2.ProxyHostLocationResource{}
	for _, location := range plan.Locations {
		item.Locations = append(item.Locations, models2.ProxyHostLocationResource{
			Path:           location.Path.ValueString(),
			ForwardScheme:  location.ForwardScheme.ValueString(),
			ForwardHost:    location.ForwardHost.ValueString(),
			ForwardPort:    uint16(location.ForwardPort.ValueInt64()),
			AdvancedConfig: location.AdvancedConfig.ValueString(),
		})
	}

	proxyHost, err := r.client.CreateProxyHost(ctx, &item)
	if err != nil {
		sentry.CaptureException(err)
		resp.Diagnostics.AddError("Error creating proxy host", "Could not create proxy host, unexpected error: "+err.Error())
		return
	}

	resp.Diagnostics.Append(plan.Load(ctx, proxyHost)...)

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *proxyHostResource) ReadImpl(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *models.ProxyHost
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	proxyHost, err := r.client.GetProxyHost(ctx, state.ID.ValueInt64Pointer())
	if err != nil {
		sentry.CaptureException(err)
		resp.Diagnostics.AddError("Error reading proxy host", "Could not read proxy host, unexpected error: "+err.Error())
		return
	}
	if proxyHost == nil {
		state = nil
	} else {
		resp.Diagnostics.Append(state.Load(ctx, proxyHost)...)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *proxyHostResource) UpdateImpl(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan models.ProxyHost
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var item models2.ProxyHostUpdate

	item.DomainNames = []string{}
	for _, domainName := range plan.DomainNames {
		item.DomainNames = append(item.DomainNames, domainName.ValueString())
	}
	item.ForwardScheme = plan.ForwardScheme.ValueString()
	item.ForwardHost = plan.ForwardHost.ValueString()
	item.ForwardPort = uint16(plan.ForwardPort.ValueInt64())
	item.CertificateID = plan.CertificateID.ValueString()
	item.SSLForced = plan.SSLForced.ValueBool()
	item.HSTSEnabled = plan.HSTSEnabled.ValueBool()
	item.HSTSSubdomains = plan.HSTSSubdomains.ValueBool()
	item.HTTP2Support = plan.HTTP2Support.ValueBool()
	item.BlockExploits = plan.BlockExploits.ValueBool()
	item.CachingEnabled = plan.CachingEnabled.ValueBool()
	item.AllowWebsocketUpgrade = plan.AllowWebsocketUpgrade.ValueBool()
	item.AccessListID = plan.AccessListID.ValueInt64()
	item.AdvancedConfig = plan.AdvancedConfig.ValueString()
	item.Meta = map[string]string{}
	item.Locations = []models2.ProxyHostLocationResource{}
	for _, location := range plan.Locations {
		item.Locations = append(item.Locations, models2.ProxyHostLocationResource{
			Path:           location.Path.ValueString(),
			ForwardScheme:  location.ForwardScheme.ValueString(),
			ForwardHost:    location.ForwardHost.ValueString(),
			ForwardPort:    uint16(location.ForwardPort.ValueInt64()),
			AdvancedConfig: location.AdvancedConfig.ValueString(),
		})
	}

	proxyHost, err := r.client.UpdateProxyHost(ctx, plan.ID.ValueInt64Pointer(), &item)
	if err != nil {
		sentry.CaptureException(err)
		resp.Diagnostics.AddError("Error updating proxy host", "Could not update proxy host, unexpected error: "+err.Error())
		return
	}

	resp.Diagnostics.Append(plan.Load(ctx, proxyHost)...)

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *proxyHostResource) DeleteImpl(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.ProxyHost
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteProxyHost(ctx, state.ID.ValueInt64Pointer())
	if err != nil {
		sentry.CaptureException(err)
		resp.Diagnostics.AddError("Error deleting proxy host", "Could not delete proxy host, unexpected error: "+err.Error())
		return
	}
}

func (r *proxyHostResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data models.ProxyHost

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.SSLForced.ValueBool() && data.CertificateID.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("ssl_forced"),
			"Certificate ID is required when SSL is forced",
			"Certificate ID is required when SSL is forced")
	}

	if data.HTTP2Support.ValueBool() && data.CertificateID.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("http2_support"),
			"Certificate ID is required when HTTP/2 Support is enabled",
			"Certificate ID is required when HTTP/2 Support is enabled")
	}

	if data.HSTSEnabled.ValueBool() && !data.SSLForced.ValueBool() {
		resp.Diagnostics.AddAttributeError(
			path.Root("hsts_enabled"),
			"SSL must be forced when HSTS is enabled",
			"SSL must be forced when HSTS is enabled")
	}

	if data.HSTSSubdomains.ValueBool() && !data.HSTSEnabled.ValueBool() {
		resp.Diagnostics.AddAttributeError(
			path.Root("hsts_subdomains"),
			"HSTS must be enabled when HSTS Subdomains is enabled",
			"HSTS must be enabled when HSTS Subdomains is enabled")
	}
}

func (r *proxyHostResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Error importing proxy host", "Could not import proxy host, unexpected error: "+err.Error())
		return
	}

	diags := resp.State.SetAttribute(ctx, path.Root("id"), types.Int64Value(id))
	resp.Diagnostics.Append(diags...)
}
