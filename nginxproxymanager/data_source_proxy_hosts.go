package nginxproxymanager

import (
	"context"
	"github.com/getsentry/sentry-go"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sander0542/terraform-provider-nginxproxymanager/client"
	"github.com/sander0542/terraform-provider-nginxproxymanager/nginxproxymanager/common"
)

var (
	_ common.IDataSource                 = &proxyHostsDataSource{}
	_ datasource.DataSourceWithConfigure = &proxyHostsDataSource{}
)

func NewProxyHostsDataSource() datasource.DataSource {
	b := &common.DataSource{Name: "proxy_hosts"}
	d := &proxyHostsDataSource{b, nil}
	b.IDataSource = d
	return d
}

type proxyHostsDataSource struct {
	*common.DataSource
	client *client.Client
}

type proxyHostsDataSourceModel struct {
	ProxyHosts []proxyHostItem `tfsdk:"proxy_hosts"`
}

type proxyHostItem struct {
	ID                    types.Int64  `tfsdk:"id"`
	CreatedOn             types.String `tfsdk:"created_on"`
	ModifiedOn            types.String `tfsdk:"modified_on"`
	OwnerUserID           types.Int64  `tfsdk:"owner_user_id"`
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

func (d *proxyHostsDataSource) SchemaImpl(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Proxy Hosts data source",
		Attributes: map[string]schema.Attribute{
			"proxy_hosts": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Description: "The ID of the proxy host.",
							Computed:    true,
						},
						"created_on": schema.StringAttribute{
							Description: "The date and time the proxy host was created.",
							Computed:    true,
						},
						"modified_on": schema.StringAttribute{
							Description: "The date and time the proxy host was last modified.",
							Computed:    true,
						},
						"owner_user_id": schema.Int64Attribute{
							Description: "The ID of the user that owns the proxy host.",
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
				},
			},
		},
	}
}

func (d *proxyHostsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*client.Client)
}

func (d *proxyHostsDataSource) ReadImpl(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data proxyHostsDataSourceModel

	proxyHosts, err := d.client.GetProxyHosts(ctx)
	if err != nil {
		sentry.CaptureException(err)
		resp.Diagnostics.AddError(
			"Error reading proxy hosts",
			"Could not read proxy hosts, unexpected error: "+err.Error())
		return
	}

	for _, value := range *proxyHosts {
		domainNames, diags := types.ListValueFrom(ctx, types.StringType, value.DomainNames)
		resp.Diagnostics.Append(diags...)
		meta, diags := types.MapValueFrom(ctx, types.StringType, value.Meta.Map())
		resp.Diagnostics.Append(diags...)

		proxyHost := proxyHostItem{
			ID:                    types.Int64Value(value.ID),
			CreatedOn:             types.StringValue(value.CreatedOn),
			ModifiedOn:            types.StringValue(value.ModifiedOn),
			OwnerUserID:           types.Int64Value(value.OwnerUserID),
			DomainNames:           domainNames,
			ForwardScheme:         types.StringValue(value.ForwardScheme),
			ForwardHost:           types.StringValue(value.ForwardHost),
			ForwardPort:           types.Int64Value(int64(value.ForwardPort)),
			CertificateID:         types.StringValue(string(value.CertificateID)),
			SSLForced:             types.BoolValue(value.SSLForced.Bool()),
			HSTSEnabled:           types.BoolValue(value.HSTSEnabled.Bool()),
			HSTSSubdomains:        types.BoolValue(value.HSTSSubdomains.Bool()),
			HTTP2Support:          types.BoolValue(value.HTTP2Support.Bool()),
			BlockExploits:         types.BoolValue(value.BlockExploits.Bool()),
			CachingEnabled:        types.BoolValue(value.CachingEnabled.Bool()),
			AllowWebsocketUpgrade: types.BoolValue(value.AllowWebsocketUpgrade.Bool()),
			AccessListID:          types.Int64Value(value.AccessListID),
			AdvancedConfig:        types.StringValue(value.AdvancedConfig),
			Enabled:               types.BoolValue(value.Enabled.Bool()),
			Meta:                  meta,
		}
		data.ProxyHosts = append(data.ProxyHosts, proxyHost)
	}

	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
