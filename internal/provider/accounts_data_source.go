package provider

import (
	"context"
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zesty-co/terraform-provider-zesty/internal/client"
	"github.com/zesty-co/terraform-provider-zesty/internal/models"
)

type AccountsDataSource struct {
	client *client.Client
}

var (
	_ datasource.DataSource              = &AccountsDataSource{}
	_ datasource.DataSourceWithConfigure = &AccountsDataSource{}
)

func NewAccountsDataSource() datasource.DataSource {
	return &AccountsDataSource{}
}

func (d *AccountsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_accounts"
}

type accountsDataSourceModel struct {
	Accounts []accountModel `tfsdk:"accounts"`
}

type accountModel struct {
	ID            types.String   `tfsdk:"id"`
	CloudProvider types.String   `tfsdk:"cloud_provider"`
	AWSRegion     types.String   `tfsdk:"aws_region"`
	RoleARN       types.String   `tfsdk:"role_arn"`
	ExternalID    types.String   `tfsdk:"external_id"`
	Products      []productModel `tfsdk:"products"`
}

type productModel struct {
	Name   types.String `tfsdk:"name"`
	Active types.Bool   `tfsdk:"active"`
	Values types.String `tfsdk:"values"`
}

// Schema defines the schema for the data source.
func (d *AccountsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the list of accounts.",
		Attributes: map[string]schema.Attribute{
			"accounts": schema.ListNestedAttribute{
				Description: "List of accounts.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Account ID",
							Computed:    true,
						},
						"cloud_provider": schema.StringAttribute{
							Description: "Name of cloud provider (e.g. AWS, GCP, Azure)",
							Computed:    true,
						},
						"role_arn": schema.StringAttribute{
							Description: "Role ARN generated on the cloud provider",
							Computed:    true,
						},
						"external_id": schema.StringAttribute{
							Description: "External ID (UUID)",
							Computed:    true,
						},
						"aws_region": schema.StringAttribute{
							Optional: true,
							Computed: false,
						},
						"products": schema.ListNestedAttribute{
							Description: "List of products activated on the account",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										Description: "Name of product (e.g. Kompass)",
										Computed:    true,
									},
									"active": schema.BoolAttribute{
										Description: "Status of product",
										Computed:    true,
									},
									"values": schema.StringAttribute{
										Description: "Key-value pairs of product-specific values",
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

func (d *AccountsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state accountsDataSourceModel

	accounts, err := d.client.GetAccounts()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Zesty Onboarded Accounts",
			err.Error(),
		)
		return
	}

	tflog.Info(ctx, "Received accounts", map[string]any{"count": len(*accounts)})

	for _, account := range *accounts {
		roleARN, exists := account.AdditionalData["roleARN"]
		if !exists {
			resp.Diagnostics.AddError(
				"Missing role ARN for account",
				account.AccountID,
			)
			return
		}
		roleARNString, ok := roleARN.(string)
		if !ok {
			resp.Diagnostics.AddError(
				"Erroneous role ARN for account",
				account.AccountID,
			)
			return
		}

		externalID, exists := account.AdditionalData["externalID"]
		if !exists {
			resp.Diagnostics.AddError(
				"Missing external ID for account",
				account.AccountID,
			)
			return
		}
		externalIDString, ok := externalID.(string)
		if !ok {
			resp.Diagnostics.AddError(
				"Erroneous external ID for account",
				account.AccountID,
			)
			return
		}
		accountState := accountModel{
			ID:            types.StringValue(account.AccountID),
			CloudProvider: types.StringValue(string(account.CloudProvider)),
			RoleARN:       types.StringValue(roleARNString),
			ExternalID:    types.StringValue(externalIDString),
		}

		var productNames []string
		for name := range account.Products {
			productNames = append(productNames, string(name))
		}
		sort.Strings(productNames)

		for _, name := range productNames {
			details := account.Products[models.Product(name)]
			accountState.Products = append(accountState.Products, productModel{
				Name:   types.StringValue(name),
				Active: types.BoolValue(details.Active),
			})
		}

		tflog.Info(ctx, "Adding account to state", map[string]any{"account": accountState})

		state.Accounts = append(state.Accounts, accountState)
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *AccountsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected: *client.Client, got: %T.\nPlease report this issue to Zesty Support.", req.ProviderData),
		)

		return
	}

	d.client = client
}
