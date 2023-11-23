// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	PPSClient "github.com/theochita/go-pleasant-password"
)

var _ provider.Provider = &PleasantpasswordProvider{}

type PleasantpasswordProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

type PleasantpasswordProviderModel struct {
	ServerURL      types.String `tfsdk:"server_url"`
	Username       types.String `tfsdk:"username"`
	Password       types.String `tfsdk:"password"`
	Allow_insecure types.Bool   `tfsdk:"allow_insecure"`
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
				MarkdownDescription: "Required: The URL of the Pleasant Password Server, Can be specified with the `PPS_SERVER_URL` environment variable  ",
				Optional:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "Required: The password of the Pleasant Password Server, Can be specified with the `PPS_PASSWORD` environment variable",
				Sensitive:           true,
				Optional:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Required: The username of the Pleasant Password Server, Can be specified with the `PPS_USERNAME` environment variable",
				Optional:            true,
			},
			"allow_insecure": schema.BoolAttribute{
				MarkdownDescription: "Allow insecure connections to the Pleasant Password Server, Can be specified with the `PPS_ALLOW_INSECURE` environment variable",
				Optional:            true,
			},
		},
	}
}

func (p *PleasantpasswordProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data PleasantpasswordProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If any of the expected configurations are unknown, return
	if data.ServerURL.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("server_url"),
			"Unknown Server URL",
			"The provider cannot create the API client as there is unknown configuration value for the server_url",
		)
		return
	}

	if data.Username.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Unknown Username",
			"The provider cannot create the API client as there is unknown configuration value for the username",
		)
		return
	}

	if data.Password.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Unknown Password",
			"The provider cannot create the API client as there is unknown configuration value for the password",
		)
		return
	}

	if data.Allow_insecure.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("allow_insecure"),
			"Unknown Allow Insecure",
			"The provider cannot create the API client as there is unknown configuration value for the allow_insecure",
		)
		return
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if resp.Diagnostics.HasError() {
		return
	}

	env_server_url := os.Getenv("PPS_SERVER_URL")
	if env_server_url != "" && data.ServerURL.IsNull() {
		data.ServerURL = types.StringValue(env_server_url)
	}

	env_username := os.Getenv("PPS_USERNAME")
	if env_username != "" && data.Username.IsNull() {
		data.Username = types.StringValue(env_username)
	}

	env_password := os.Getenv("PPS_PASSWORD")
	if env_password != "" && data.Password.IsNull() {
		data.Password = types.StringValue(env_password)
	}

	env_ssl_insecure := os.Getenv("PPS_ALLOW_INSECURE")
	if env_ssl_insecure != "" && data.Allow_insecure.IsNull() {
		bool_env_ssl_insecure, err := strconv.ParseBool(env_ssl_insecure)
		if err != nil {
			resp.Diagnostics.AddError("Env ALLOW_INSECURE, Must be bool", err.Error())
		}
		data.Allow_insecure = types.BoolValue(bool_env_ssl_insecure)
	}

	// Last check for required values that has not been set from terraform or env
	if data.ServerURL.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("server_url"),
			"Missing API Server URL",
			"The provider cannot create the API client as there is a missing or empty value for the API URL. "+
				"Set the value in the configuration or use the PPS_SERVER_URL environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
		return
	}

	if data.Username.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing Username",
			"The provider cannot create the API client as there is a missing or empty value for the username. "+
				"Set the value in the configuration or use the PPS_USERNAME environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
		return
	}

	if data.Password.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Missing Password",
			"The provider cannot create the API client as there is a missing or empty value for the password. "+
				"Set the value in the configuration or use the PPS_PASSWORD environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
		return
	}

	if data.Allow_insecure.ValueBool() {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	} else {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: false}
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

	res, httperr, err := client.AuthenticationAPI.PostOauthToken(clientctx).GrantType(grantType).Username(username).Password(password).Execute()

	if err != nil {
		fmt.Printf("error in auth request {%s}  \n", err)
		if httperr != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Auth Request Failed with %s", err.Error()), fmt.Sprintf("%s", httperr.Body)) // ADD REAL ERROR MODES TO THE OPENAPI SPECC
		} else {
			resp.Diagnostics.AddError("Auth Request Failed", err.Error())
		}

		return
	} else {

		//fmt.Printf("Authenticated Successful\n")
		//fmt.Printf("Bearer token: %s\n", *res.AccessToken)
		clientctx = context.WithValue(context.Background(), PPSClient.ContextAccessToken, *res.AccessToken)
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
