package nginxproxymanager

import (
	"context"
	"fmt"
	"os"

	"github.com/getsentry/sentry-go"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/sander0542/terraform-provider-nginxproxymanager/client"
)

var (
	_ provider.Provider = &nginxproxymanagerProvider{}
)

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &nginxproxymanagerProvider{
			Version: version,
		}
	}
}

type nginxproxymanagerProvider struct {
	Version string
}

type nginxproxymanagerProviderModel struct {
	Host     types.String `tfsdk:"host"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

func (p *nginxproxymanagerProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "nginxproxymanager"
}

func (p *nginxproxymanagerProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with Nginx Proxy Manager API.",
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Description: "URI for Nginx Proxy Manager API. May also be provided via NGINX_PROXY_MANAGER_HOST environment variable.",
				Optional:    true,
			},
			"username": schema.StringAttribute{
				Description: "Username for Nginx Proxy Manager API. May also be provided via NGINX_PROXY_MANAGER_USERNAME environment variable.",
				Optional:    true,
			},
			"password": schema.StringAttribute{
				Description: "Password for Nginx Proxy Manager API. May also be provided via NGINX_PROXY_MANAGER_PASSWORD environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

func (p *nginxproxymanagerProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Nginx Proxy Manager provider")

	var config nginxproxymanagerProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown Nginx Proxy Manager API Host",
			"The provider cannot create the Nginx Proxy Manager API client as there is an unknown configuration value for the Nginx Proxy Manager API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the NGINX_PROXY_MANAGER_HOST environment variable.",
		)
	}

	if config.Username.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Unknown Nginx Proxy Manager API Username",
			"The provider cannot create the Nginx Proxy Manager API client as there is an unknown configuration value for the Nginx Proxy Manager API username. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the NGINX_PROXY_MANAGER_USERNAME environment variable.",
		)
	}

	if config.Password.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Unknown Nginx Proxy Manager API Password",
			"The provider cannot create the Nginx Proxy Manager API client as there is an unknown configuration value for the Nginx Proxy Manager API password. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the NGINX_PROXY_MANAGER_PASSWORD environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	host := os.Getenv("NGINX_PROXY_MANAGER_HOST")
	username := os.Getenv("NGINX_PROXY_MANAGER_USERNAME")
	password := os.Getenv("NGINX_PROXY_MANAGER_PASSWORD")

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.Username.IsNull() {
		username = config.Username.ValueString()
	}

	if !config.Password.IsNull() {
		password = config.Password.ValueString()
	}

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing Nginx Proxy Manager API Host",
			"The provider cannot create the Nginx Proxy Manager API client as there is a missing or empty value for the Nginx Proxy Manager API host. "+
				"Set the host value in the configuration or use the NGINX_PROXY_MANAGER_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if username == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing Nginx Proxy Manager API Username",
			"The provider cannot create the Nginx Proxy Manager API client as there is a missing or empty value for the Nginx Proxy Manager API username. "+
				"Set the username value in the configuration or use the NGINX_PROXY_MANAGER_USERNAME environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if password == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Missing Nginx Proxy Manager API Password",
			"The provider cannot create the Nginx Proxy Manager API client as there is a missing or empty value for the Nginx Proxy Manager API password. "+
				"Set the password value in the configuration or use the NGINX_PROXY_MANAGER_PASSWORD environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Initializing Sentry")

	environment := "production"
	if p.Version == "dev" || p.Version == "test" {
		environment = p.Version
	}

	err := sentry.Init(sentry.ClientOptions{
		Dsn:              "https://2ec435e840424aeb8c40b56dea37e4dd@o476647.ingest.sentry.io/4505102669447168",
		EnableTracing:    true,
		Environment:      environment,
		Release:          fmt.Sprintf("terraform-provider-nginxproxymanager@%s", p.Version),
		TracesSampleRate: 1.0,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Initialize Sentry",
			"An unexpected error occurred when initializing Sentry. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Sentry Error: "+err.Error(),
		)
		return
	}
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetContext("terraform", map[string]interface{}{
			"version": req.TerraformVersion,
		})
	})

	tflog.Debug(ctx, "Creating Nginx Proxy Manager client")

	npmClient, err := client.NewClient(&host, &username, &password)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Nginx Proxy Manager API Client",
			"An unexpected error occurred when creating the Nginx Proxy Manager API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Nginx Proxy Manager Client Error: "+err.Error(),
		)
		return
	}
	npmClient.HTTPClient.Transport = newTracingTransport(npmClient.HTTPClient.Transport)

	ctx = tflog.SetField(ctx, "nginxproxymanager_host", host)
	ctx = tflog.SetField(ctx, "nginxproxymanager_username", username)
	ctx = tflog.SetField(ctx, "nginxproxymanager_password", password)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "nginxproxymanager_password")

	resp.DataSourceData = npmClient
	resp.ResourceData = npmClient

	tflog.Info(ctx, "Configured Nginx Proxy Manager client", map[string]any{"success": true})
}

func (p *nginxproxymanagerProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewCertificatesDataSource,
		NewCertificateDataSource,
		NewProxyHostDataSource,
		NewProxyHostsDataSource,
	}
}

func (p *nginxproxymanagerProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewProxyHostResource,
	}
}
