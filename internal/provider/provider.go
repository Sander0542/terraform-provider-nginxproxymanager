// Copyright (c) Sander Jochems
// SPDX-License-Identifier: MIT

package provider

import (
	"context"
	"net/url"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/sander0542/nginxproxymanager-go"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure NginxProxyManagerProvider satisfies various provider interfaces.
var _ provider.Provider = &NginxProxyManagerProvider{}
var _ provider.ProviderWithFunctions = &NginxProxyManagerProvider{}
var _ provider.ProviderWithEphemeralResources = &NginxProxyManagerProvider{}

// NginxProxyManagerProvider defines the provider implementation.
type NginxProxyManagerProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// NginxProxyManagerProviderModel describes the provider data model.
type NginxProxyManagerProviderModel struct {
	Url      types.String `tfsdk:"url"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

type NginxProxyManagerProviderData struct {
	Client *nginxproxymanager.APIClient
	Auth   context.Context
}

func (p *NginxProxyManagerProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "nginxproxymanager"
	resp.Version = p.version
}

func (p *NginxProxyManagerProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use the Nginx Proxy Manager (NPM) provider to interact with resources from Nginx Proxy Manager.",
		Attributes: map[string]schema.Attribute{
			"url": schema.StringAttribute{
				MarkdownDescription: "Full Nginx Proxy Manager URL with protocol and port (e.g. `http://localhost:81`). You should **NOT** supply any path (`/api`), the SDK will use the appropriate paths. Can be specified via the `NGINXPROXYMANAGER_URL` environment variable.",
				Optional:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Username for Nginx Proxy Manager authentication. Can be specified via the `NGINXPROXYMANAGER_USERNAME` environment variable.",
				Optional:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "Password for Nginx Proxy Manager authentication. Can be specified via the `NGINXPROXYMANAGER_PASSWORD` environment variable.",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *NginxProxyManagerProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {

	var data NginxProxyManagerProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		tflog.Trace(ctx, "Failed to load provider configuration")
		return
	}

	// Url
	apiUrl := data.Url.ValueString()
	if apiUrl == "" {
		tflog.Trace(ctx, "Url is not set in configuration, checking environment variables")
		apiUrl = os.Getenv("NGINXPROXYMANAGER_URL")
	}

	parsedUrl, err := url.Parse(apiUrl)
	if err != nil {
		resp.Diagnostics.AddAttributeError(
			path.Root("url"),
			"Url is required",
			"Please provide a valid url value",
		)

		return
	}
	parsedUrl = parsedUrl.JoinPath("/api")

	// Username
	username := data.Username.ValueString()
	if username == "" {
		tflog.Trace(ctx, "Username is not set in configuration, checking environment variables")
		username = os.Getenv("NGINXPROXYMANAGER_USERNAME")
	}

	if username == "" {
		tflog.Debug(ctx, "Username is not set in configuration or environment variables")
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Username is required",
			"Please provide a username value",
		)
	}

	// Password
	password := data.Password.ValueString()
	if password == "" {
		tflog.Trace(ctx, "Password is not set in configuration, checking environment variables")
		password = os.Getenv("NGINXPROXYMANAGER_PASSWORD")
	}

	if password == "" {
		tflog.Debug(ctx, "Password is not set in configuration or environment variables")
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Password is required",
			"Please provide a password value",
		)
	}

	if resp.Diagnostics.HasError() {
		tflog.Trace(ctx, "Failed to load provider configuration")
		return
	}

	tflog.MaskMessageStrings(ctx, username, password)
	tflog.Info(ctx, "Initializing the Nginx Proxy Manager API client")

	config := nginxproxymanager.NewConfiguration()
	config.Servers[0].URL = parsedUrl.String()
	client := nginxproxymanager.NewAPIClient(config)

	auth := context.Background()

	tflog.Info(ctx, "Authenticating with the Nginx Proxy Manager API")

	tokenRequest := nginxproxymanager.RequestTokenRequest{
		Identity: username,
		Secret:   password,
	}

	tokenResponse, _, err := client.TokensAPI.RequestToken(auth).RequestTokenRequest(tokenRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Failed to authenticate with the Nginx Proxy Manager API", err.Error())
		return
	}

	tflog.Info(ctx, "Successfully authenticated with the Nginx Proxy Manager API")

	auth = context.WithValue(auth, nginxproxymanager.ContextAccessToken, tokenResponse.GetToken())

	providerData := NginxProxyManagerProviderData{
		Auth:   auth,
		Client: client,
	}

	resp.DataSourceData = &providerData
	resp.ResourceData = &providerData
	resp.EphemeralResourceData = &providerData

	tflog.Info(ctx, "Successfully initialized the Nginx Proxy Manager API client")
}

func (p *NginxProxyManagerProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewAccessListResource,
		NewCertificateCustomResource,
		NewCertificateLetsencryptResource,
		NewDeadHostResource,
		NewProxyHostResource,
		NewRedirectionHostResource,
		NewSettingsResource,
		NewStreamResource,
	}
}

func (p *NginxProxyManagerProvider) EphemeralResources(ctx context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{
		NewUserTokenEphemeralResource,
	}
}

func (p *NginxProxyManagerProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewAccessListDataSource,
		NewAccessListsDataSource,
		NewCertificateDataSource,
		NewCertificatesDataSource,
		NewDeadHostDataSource,
		NewDeadHostsDataSource,
		NewProxyHostDataSource,
		NewProxyHostsDataSource,
		NewRedirectionHostDataSource,
		NewRedirectionHostsDataSource,
		NewSettingsDataSource,
		NewStreamDataSource,
		NewStreamsDataSource,
		NewUserDataSource,
		NewUserMeDataSource,
		NewUsersDataSource,
		NewVersionDataSource,
	}
}

func (p *NginxProxyManagerProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &NginxProxyManagerProvider{
			version: version,
		}
	}
}
