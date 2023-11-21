// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	PPSClient "github.com/theochita/go-pleasant-password"
	"github.com/theochita/terraform-provider-pleasant-password-server/internal/provider/models"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &SearchDataSource{}

func NewSearchDataSource() datasource.DataSource {
	return &SearchDataSource{}
}

// ExampleDataSource defines the data source implementation.
type SearchDataSource struct {
	client *PPSClient.APIClient
	ctx    *context.Context
}

// ExampleDataSourceModel describes the data source data model.
type SearchDataSourceModel struct {
	Search      types.String                           `tfsdk:"search"`
	Credentials []models.V6CredentialSearchResult      `tfsdk:"credentials"`
	Folders     []models.V6CredentialGroupSearchResult `tfsdk:"folders"`
}

func (d SearchDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_search"
}

func (d *SearchDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Example data source",

		Attributes: map[string]schema.Attribute{
			"search": schema.StringAttribute{
				MarkdownDescription: "Example identifier",
				Required:            true,
			},
			"credentials": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "Example identifier",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Example identifier",
							Computed:            true,
						},
						"username": schema.StringAttribute{
							MarkdownDescription: "Example identifier",
							Computed:            true,
						},
						"url": schema.StringAttribute{
							MarkdownDescription: "Example identifier",
							Computed:            true,
						},
						"notes": schema.StringAttribute{
							MarkdownDescription: "Example identifier",
							Computed:            true,
						},
						"folderid": schema.StringAttribute{
							MarkdownDescription: "Example identifier",
							Computed:            true,
						},
						"path": schema.StringAttribute{
							MarkdownDescription: "Example identifier",
							Computed:            true,
						},
					},
				},
			},
			"folders": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "Example identifier",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Example identifier",
							Computed:            true,
						},
						"fullpath": schema.StringAttribute{
							MarkdownDescription: "Example identifier",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *SearchDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	providerclient, ok := req.ProviderData.(ProviderClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *PPSClient.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = &providerclient.Client
	d.ctx = &providerclient.Ctx

}

func (d *SearchDataSource) fetchCredentials(res []PPSClient.V6CredentialSearchResult) []models.V6CredentialSearchResult {
	var creds = []models.V6CredentialSearchResult{}
	for _, v := range res {
		cred := models.V6CredentialSearchResult{}
		cred.Id = types.StringValue(v.GetId())
		cred.Name = types.StringValue(v.GetName())
		cred.Username = types.StringValue(v.GetUsername())
		cred.Url = types.StringValue(v.GetUrl())
		cred.Notes = types.StringValue(v.GetNotes())
		cred.FolderId = types.StringValue(v.GetGroupId())
		cred.Path = types.StringValue(v.GetPath())
		creds = append(creds, cred)

	}
	return creds
}

func (d *SearchDataSource) fetchFolders(res []PPSClient.V6CredentialGroupSearchResult) []models.V6CredentialGroupSearchResult {
	var folders = []models.V6CredentialGroupSearchResult{}
	for _, v := range res {
		folder := models.V6CredentialGroupSearchResult{}
		folder.Id = types.StringValue(v.GetId())
		folder.Name = types.StringValue(v.GetName())
		folder.FullPath = types.StringValue(v.GetFullPath())
		folders = append(folders, folder)

	}
	return folders
}

func (d *SearchDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SearchDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := d.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read example, got error: %s", err))
	//     return
	// }

	// For the purposes of this example code, hardcoding a response value to
	// save into the Terraform state.
	//data.Id = types.StringValue("example-id")

	client := d.client

	params := PPSClient.NewV6SearchInputWithDefaults()
	params.Search = data.Search.ValueStringPointer()

	res, httpres, err := client.DefaultAPI.PostV6Search(*d.ctx).V6SearchInput(*params).Execute()
	if err != nil {
		resp.Diagnostics.AddError("failure to invoke API: ", err.Error())
		return
	}
	if httpres.StatusCode != 200 {
		resp.Diagnostics.AddError("Got an unexpected response code", fmt.Sprintf("Got an unexpected response code %v", httpres.StatusCode))
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	data.Credentials = d.fetchCredentials(res.Credentials)
	data.Folders = d.fetchFolders(res.Groups)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
