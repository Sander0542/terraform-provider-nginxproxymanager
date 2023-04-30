package nginxproxymanager

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sander0542/terraform-provider-nginxproxymanager/client"
	"github.com/sander0542/terraform-provider-nginxproxymanager/nginxproxymanager/common"
)

var (
	_ common.IDataSource                 = &proxyHostDataSource{}
	_ datasource.DataSourceWithConfigure = &proxyHostDataSource{}
)

func NewProxyHostDataSource() datasource.DataSource {
	b := &common.DataSource{}
	d := &proxyHostDataSource{b, nil}
	b.IDataSource = d
	return d
}

type proxyHostDataSource struct {
	*common.DataSource
	client *client.Client
}

type proxyHostDataSourceModel struct {
	ID                    types.Int64  `tfsdk:"id"`
	CreatedOn             types.String `tfsdk:"created_on"`
	ModifiedOn            types.String `tfsdk:"modified_on"`
	DomainNames           types.List   `tfsdk:"domain_names"`
	ForwardScheme         types.String `tfsdk:"forward_scheme"`
	ForwardHost           types.String `tfsdk:"forward_host"`
	ForwardPort           types.Int64  `tfsdk:"forward_port"`
	CertificateID         types.String `tfsdk:"certificate_id"`
	SSLForced             types.Bool   `tfsdk:"ssl_forced"`
	HSTSEnabled           types.Bool   `tfsdk:"hsts_enabled"`
	HSTSSubdomains        types.Bool   `tfsdk:"hsts_subdomains"`
	HTTP2Support          types.Bool   `tfsdk:"http2_support"`
	BlockExploits         types.Bool   `tfsdk:"block_exploits"`
	CachingEnabled        types.Bool   `tfsdk:"caching_enabled"`
	AllowWebsocketUpgrade types.Bool   `tfsdk:"allow_websocket_upgrade"`
	AccessListID          types.Int64  `tfsdk:"access_list_id"`
	AdvancedConfig        types.String `tfsdk:"advanced_config"`
	Enabled               types.Bool   `tfsdk:"enabled"`
	Meta                  types.Map    `tfsdk:"meta"`
	Locations             types.List   `tfsdk:"locations"`
}

func (d *proxyHostDataSource) MetadataImpl(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_proxy_host"
}

func (d *proxyHostDataSource) SchemaImpl(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a proxy host by ID.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "The ID of the proxy host.",
				Required:    true,
			},
			"created_on": schema.StringAttribute{
				Description: "The date and time the proxy host was created.",
				Computed:    true,
			},
			"modified_on": schema.StringAttribute{
				Description: "The date and time the proxy host was last modified.",
				Computed:    true,
			},
			"domain_names": schema.ListAttribute{
				Description: "The domain names associated with the proxy host.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"forward_scheme": schema.StringAttribute{
				Description: "The scheme used to forward requests to the proxy host. Can be either `http` or `https`.",
				Computed:    true,
			},
			"forward_host": schema.StringAttribute{
				Description: "The host used to forward requests to the proxy host.",
				Computed:    true,
			},
			"forward_port": schema.Int64Attribute{
				Description: "The port used to forward requests to the proxy host.",
				Computed:    true,
			},
			"certificate_id": schema.StringAttribute{
				Description: "The ID of the certificate used by the proxy host.",
				Computed:    true,
			},
			"ssl_forced": schema.BoolAttribute{
				Description: "Whether SSL is forced for the proxy host.",
				Computed:    true,
			},
			"hsts_enabled": schema.BoolAttribute{
				Description: "Whether HSTS is enabled for the proxy host.",
				Computed:    true,
			},
			"hsts_subdomains": schema.BoolAttribute{
				Description: "Whether HSTS is enabled for subdomains of the proxy host.",
				Computed:    true,
			},
			"http2_support": schema.BoolAttribute{
				Description: "Whether HTTP/2 is supported for the proxy host.",
				Computed:    true,
			},
			"block_exploits": schema.BoolAttribute{
				Description: "Whether exploits are blocked for the proxy host.",
				Computed:    true,
			},
			"caching_enabled": schema.BoolAttribute{
				Description: "Whether caching is enabled for the proxy host.",
				Computed:    true,
			},
			"allow_websocket_upgrade": schema.BoolAttribute{
				Description: "Whether websocket upgrades are allowed for the proxy host.",
				Computed:    true,
			},
			"access_list_id": schema.Int64Attribute{
				Description: "The ID of the access list used by the proxy host.",
				Computed:    true,
			},
			"advanced_config": schema.StringAttribute{
				Description: "The advanced configuration used by the proxy host.",
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether the proxy host is enabled.",
				Computed:    true,
			},
			"meta": schema.MapAttribute{
				Description: "The meta data associated with the proxy host.",
				ElementType: types.StringType,
				Computed:    true,
			},
			"locations": schema.ListNestedAttribute{
				Description: "The locations associated with the proxy host.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"path": schema.StringAttribute{
							Description: "The path associated with the location.",
							Computed:    true,
						},
						"forward_scheme": schema.StringAttribute{
							Description: "The scheme used to forward requests to the location. Can be either `http` or `https`.",
							Computed:    true,
						},
						"forward_host": schema.StringAttribute{
							Description: "The host used to forward requests to the location.",
							Computed:    true,
						},
						"forward_port": schema.Int64Attribute{
							Description: "The port used to forward requests to the location.",
							Computed:    true,
						},
						"advanced_config": schema.StringAttribute{
							Description: "The advanced configuration used by the location.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *proxyHostDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*client.Client)
}

func (d *proxyHostDataSource) ReadImpl(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data proxyHostDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	proxyHost, err := d.client.GetProxyHost(data.ID.ValueInt64Pointer())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading proxy host",
			"Could not read proxy host, unexpected error: "+err.Error())
		return
	}
	if proxyHost == nil {
		resp.Diagnostics.AddError(
			"Error reading proxy host",
			fmt.Sprintf("No proxy host found with ID: %d", data.ID))
		return
	}

	domainNames, diags := types.ListValueFrom(ctx, types.StringType, proxyHost.DomainNames)
	resp.Diagnostics.Append(diags...)
	meta, diags := types.MapValueFrom(ctx, types.StringType, proxyHost.Meta.Map())
	resp.Diagnostics.Append(diags...)

	data.ID = types.Int64Value(proxyHost.ID)
	data.CreatedOn = types.StringValue(proxyHost.CreatedOn)
	data.ModifiedOn = types.StringValue(proxyHost.ModifiedOn)
	data.DomainNames = domainNames
	data.ForwardScheme = types.StringValue(proxyHost.ForwardScheme)
	data.ForwardHost = types.StringValue(proxyHost.ForwardHost)
	data.ForwardPort = types.Int64Value(int64(proxyHost.ForwardPort))
	data.CertificateID = types.StringValue(string(proxyHost.CertificateID))
	data.SSLForced = types.BoolValue(proxyHost.SSLForced.Bool())
	data.HSTSEnabled = types.BoolValue(proxyHost.HSTSEnabled.Bool())
	data.HSTSSubdomains = types.BoolValue(proxyHost.HSTSSubdomains.Bool())
	data.HTTP2Support = types.BoolValue(proxyHost.HTTP2Support.Bool())
	data.BlockExploits = types.BoolValue(proxyHost.BlockExploits.Bool())
	data.CachingEnabled = types.BoolValue(proxyHost.CachingEnabled.Bool())
	data.AllowWebsocketUpgrade = types.BoolValue(proxyHost.AllowWebsocketUpgrade.Bool())
	data.AccessListID = types.Int64Value(proxyHost.AccessListID)
	data.AdvancedConfig = types.StringValue(proxyHost.AdvancedConfig)
	data.Enabled = types.BoolValue(proxyHost.Enabled.Bool())
	data.Meta = meta

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
