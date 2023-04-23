package nginxproxymanager

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sander0542/terraform-provider-nginxproxymanager/client"
)

var (
	_ datasource.DataSource              = &proxyHostsDataSource{}
	_ datasource.DataSourceWithConfigure = &proxyHostsDataSource{}
)

func NewProxyHostsDataSource() datasource.DataSource {
	return &proxyHostsDataSource{}
}

type proxyHostsDataSource struct {
	client *client.Client
}

type proxyHostsDataSourceModel struct {
	ProxyHosts []proxyHost `tfsdk:"proxy_hosts"`
}

type proxyHost struct {
	ID                    types.Int64  `tfsdk:"id"`
	CreatedAt             types.String `tfsdk:"created_at"`
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
}

func (d *proxyHostsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_proxy_hosts"
}

func (d *proxyHostsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Proxy Hosts data source",
		Attributes: map[string]schema.Attribute{
			"proxy_hosts": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Computed: true,
						},
						"created_at": schema.StringAttribute{
							Computed: true,
						},
						"modified_on": schema.StringAttribute{
							Computed: true,
						},
						"domain_names": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
						},
						"forward_scheme": schema.StringAttribute{
							Computed: true,
						},
						"forward_host": schema.StringAttribute{
							Computed: true,
						},
						"forward_port": schema.Int64Attribute{
							Computed: true,
						},
						"certificate_id": schema.StringAttribute{
							Computed: true,
						},
						"ssl_forced": schema.BoolAttribute{
							Computed: true,
						},
						"hsts_enabled": schema.BoolAttribute{
							Computed: true,
						},
						"hsts_subdomains": schema.BoolAttribute{
							Computed: true,
						},
						"http2_support": schema.BoolAttribute{
							Computed: true,
						},
						"block_exploits": schema.BoolAttribute{
							Computed: true,
						},
						"caching_enabled": schema.BoolAttribute{
							Computed: true,
						},
						"allow_websocket_upgrade": schema.BoolAttribute{
							Computed: true,
						},
						"access_list_id": schema.Int64Attribute{
							Computed: true,
						},
						"advanced_config": schema.StringAttribute{
							Computed: true,
						},
						"enabled": schema.BoolAttribute{
							Computed: true,
						},
						"meta": schema.MapAttribute{
							ElementType: types.StringType,
							Computed:    true,
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

func (d *proxyHostsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data proxyHostsDataSourceModel

	proxyHosts, err := d.client.GetProxyHosts()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Nginx Proxy Manager Proxy Host",
			err.Error(),
		)
		return
	}

	for _, value := range *proxyHosts {
		domainNames, diags := types.ListValueFrom(ctx, types.StringType, value.DomainNames)
		resp.Diagnostics.Append(diags...)
		meta, diags := types.MapValueFrom(ctx, types.StringType, value.Meta.Map())
		resp.Diagnostics.Append(diags...)

		proxyHost := proxyHost{
			ID:                    types.Int64Value(value.ID),
			CreatedAt:             types.StringValue(value.CreatedAt),
			ModifiedOn:            types.StringValue(value.ModifiedOn),
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
