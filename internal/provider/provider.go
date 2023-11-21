// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	PPSClient "github.com/theochita/go-pleasant-password"
)

// Ensure ScaffoldingProvider satisfies various provider interfaces.
var _ provider.Provider = &PleasantpasswordProvider{}

// ScaffoldingProvider defines the provider implementation.
type PleasantpasswordProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// ScaffoldingProviderModel describes the provider data model.
type PleasantpasswordProviderModel struct {
	ServerURL      types.String `tfsdk:"server_url"`
	Username       types.String `tfsdk:"username"`
	Password       types.String `tfsdk:"password"`
	Allow_insecure types.Bool   `tfsdk:"allow_insecure"`
	OTP_CODE       types.String `tfsdk:"opt_code"`
	OTP_PROVIDER   types.String `tfsdk:"opt_provider"`
}

type ProviderClient struct {
	Client PPSClient.APIClient
	Ctx    context.Context
}

func (p *PleasantpasswordProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "pleasantpassword"
	resp.Version = p.version
}

func (p *PleasantpasswordProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: `PPS API : Programatically control Pleasant Password Server.`,
		Attributes: map[string]schema.Attribute{
			"server_url": schema.StringAttribute{
				MarkdownDescription: "Server URL (defaults to https://api.airbyte.com/v1)",
				Required:            true,
			},
			"password": schema.StringAttribute{
				Sensitive: true,
				Required:  true,
			},
			"username": schema.StringAttribute{
				Required: true,
			},
			"opt_code": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
			"opt_provider": schema.StringAttribute{
				Optional: true,
			},
			"allow_insecure": schema.BoolAttribute{
				Optional: true,
			},
		},
	}
}

func (p *PleasantpasswordProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data PleasantpasswordProviderModel

	if data.Allow_insecure.IsNull() {
		data.Allow_insecure = types.BoolValue(false)
	}

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	// if data.Endpoint.IsNull() { /* ... */ }

	// Example client configuration for data sources and resources

	if data.Allow_insecure.ValueBool() {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	cfg := PPSClient.NewConfiguration()
	cfg.Host = data.ServerURL.ValueString()
	cfg.Scheme = "https"
	client := PPSClient.NewAPIClient(cfg)

	clientctx := context.Background()
	//param := PPSClient.NewOauth2TokenInputWithDefaults()

	grantType := "password"
	password := data.Password.ValueString()
	username := data.Username.ValueString()

	opt_code := data.OTP_CODE.ValueString()
	opt_provider := data.OTP_PROVIDER.ValueString()

	var override_bearertoken string = ""

	override_bearertoken = ""

	if override_bearertoken == "" {

		res, httperr, err := client.AuthenticationAPI.PostOauthToken(clientctx).GrantType(grantType).Username(username).Password(password).XPleasantOTP(opt_code).XPleasantOTPProvider(opt_provider).Execute()

		if err != nil {
			fmt.Printf("error in auth request {%s}  \n", err)
			if httperr != nil {
				resp.Diagnostics.AddError(fmt.Sprintf("Auth Request Failed with %s", err.Error()), fmt.Sprintf("%s", httperr.Body)) // ADD REAL ERROR MODES TO THE OPENAPI SPECC
			} else {
				resp.Diagnostics.AddError("Auth Request Failed", err.Error())
			}

			return
		} else {

			fmt.Printf("Authenticated Successful\n")
			fmt.Printf("Bearer token: %s\n", *res.AccessToken)
			clientctx = context.WithValue(context.Background(), PPSClient.ContextAccessToken, *res.AccessToken)

		}

	} else {
		clientctx = context.WithValue(context.Background(), PPSClient.ContextAccessToken, override_bearertoken)
	}

	resp.DataSourceData = ProviderClient{Client: *client, Ctx: clientctx}
	resp.ResourceData = ProviderClient{Client: *client, Ctx: clientctx}

}

func (p *PleasantpasswordProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewFolderResource,
		NewCredentialResource,
	}
}

func (p *PleasantpasswordProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewFolderDataSource,
		NewCredentialDataSource,
		NewSearchDataSource,
		NewFolderRootDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &PleasantpasswordProvider{
			version: version,
		}
	}
}
