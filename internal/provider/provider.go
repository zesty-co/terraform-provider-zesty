package provider

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zesty-co/terraform-provider-zesty/internal/client"
	"github.com/zesty-co/terraform-provider-zesty/internal/models"
)

type ZestyProvider struct {
	version string
}

type ZestyProviderModel struct {
	Host  types.String `tfsdk:"host"`
	Token types.String `tfsdk:"token"`
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &ZestyProvider{}
	}
}

// Metadata returns the provider type name.
func (p *ZestyProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "zesty"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *ZestyProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Description: "URI for Zesty API. May also be provided by the ZESTY_HOST environment variable.",
				Optional:    true,
			},
			"token": schema.StringAttribute{
				Description: "Token for Zesty API. May also be provided by the ZESTY_API_TOKEN environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

// Configure prepares a Zesty API client for data sources and resources.
func (p *ZestyProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Zesty API client")
	var config ZestyProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown Zesty API Host",
			"The provider cannot create the Zesty API client as there is an unknown configuration value for the Zesty API host.",
		)
	}

	if config.Token.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Unknown Zesty API Token",
			"The provider cannot create the Zesty API client as there is an unknown configuration value for the Zesty API username.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	host := os.Getenv("ZESTY_HOST")
	token := os.Getenv("ZESTY_API_TOKEN")

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.Token.IsNull() {
		token = config.Token.ValueString()
	}

	if host == "" {
		host = models.DefaultHostURL
	}

	if token == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Missing Zesty API Token",
			"The provider cannot create the Zesty API client as there is a missing or empty value for the Zesty API token.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "zesty_host", host)
	ctx = tflog.SetField(ctx, "zesty_api_token", token)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "zesty_api_token")
	tflog.Debug(ctx, "Creating Zesty API client")

	client, err := client.NewClient(&host, token)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Zesty API Client",
			fmt.Sprintf("An unexpected error occurred when creating the Zesty API client. Error: %s", err),
		)
		return
	}

	err = client.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Validate Zesty API Client",
			fmt.Sprintf("An unexpected error occurred when validating the Zesty API. Error: %s", err),
		)
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured Zesty API client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *ZestyProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewAccountsDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *ZestyProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewAccountResource,
	}
}
