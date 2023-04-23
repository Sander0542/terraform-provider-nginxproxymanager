package nginxproxymanager

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sander0542/terraform-provider-nginxproxymanager/client"
)

var (
	_ datasource.DataSource              = &proxyHostDataSource{}
	_ datasource.DataSourceWithConfigure = &proxyHostDataSource{}
)

func NewProxyHostDataSource() datasource.DataSource {
	return &proxyHostDataSource{}
}

type proxyHostDataSource struct {
	client *client.Client
}

type proxyHostDataSourceModel struct {
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

func (d *proxyHostDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_proxy_host"
}

func (d *proxyHostDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Proxy Host data source",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Required: true,
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
	}
}

func (d *proxyHostDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*client.Client)
}

func (d *proxyHostDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data proxyHostDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	proxyHost, err := d.client.GetProxyHost(data.ID.ValueInt64Pointer())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Nginx Proxy Manager Proxy Host",
			err.Error(),
		)
		return
	}

	domainNames, diags := types.ListValueFrom(ctx, types.StringType, proxyHost.DomainNames)
	resp.Diagnostics.Append(diags...)
	meta, diags := types.MapValueFrom(ctx, types.StringType, proxyHost.Meta.Map())
	resp.Diagnostics.Append(diags...)

	data.ID = types.Int64Value(proxyHost.ID)
	data.CreatedAt = types.StringValue(proxyHost.CreatedAt)
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
