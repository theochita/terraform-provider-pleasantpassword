// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	PPSClient "github.com/theochita/go-pleasant-password"
	"github.com/theochita/terraform-provider-pleasant-password-server/internal/provider/models"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &CredentialDataSource{}

func NewCredentialDataSource() datasource.DataSource {
	return &CredentialDataSource{}
}

type CredentialDataSource struct {
	client *PPSClient.APIClient
	ctx    *context.Context
}

type CredentialDataSourceModel struct {
	Id           types.String `tfsdk:"id"`
	CredentialID types.String `tfsdk:"credential_id"`
	Tags         []models.Tag `tfsdk:"tags"`
	Name         types.String `tfsdk:"name"`
	Username     types.String `tfsdk:"username"`
	Password     types.String `tfsdk:"password"`
	Url          types.String `tfsdk:"url"`
	Notes        types.String `tfsdk:"notes"`
	FolderId     types.String `tfsdk:"folderid"`
	Created      types.String `tfsdk:"created"`
	Modified     types.String `tfsdk:"modified"`
	Expires      types.String `tfsdk:"expires"`
}

func (d CredentialDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_credential"
}

func (d *CredentialDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "The `credential` data source can be used to access information about a credential.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the credential",
				Computed:            true,
			},
			"credential_id": schema.StringAttribute{
				MarkdownDescription: "The identifier of the credential",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the credential",
				Computed:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "The username of the credential",
				Computed:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "The password of the credential",
				Computed:            true,
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "The URL of the credential",
				Computed:            true,
			},
			"notes": schema.StringAttribute{
				MarkdownDescription: "The notes of the credential",
				Computed:            true,
			},
			"folderid": schema.StringAttribute{
				MarkdownDescription: "The folder ID of the credential",
				Computed:            true,
			},
			"created": schema.StringAttribute{
				MarkdownDescription: "The creation date of the credential",
				Computed:            true,
			},
			"modified": schema.StringAttribute{
				MarkdownDescription: "The modification date of the credential",
				Computed:            true,
			},
			"expires": schema.StringAttribute{
				MarkdownDescription: "The expiration date of the credential",
				Computed:            true,
			},

			"tags": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "The name of the tag",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *CredentialDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *CredentialDataSource) fetchTags(res []PPSClient.V6TagResult) []models.Tag {
	var tags = []models.Tag{}
	for _, v := range res {
		tag := models.Tag{}
		tag.Name = types.StringValue(v.GetName())
		tags = append(tags, tag)
	}

	return tags

}

func (d *CredentialDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CredentialDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	credential_id := data.CredentialID.ValueString()

	client := d.client

	res, httpres, err := client.DefaultAPI.GetV6CredentialsByID(*d.ctx, credential_id).Execute()
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

	data.Id = types.StringValue(res.GetId())
	data.Name = types.StringValue(res.GetName())
	data.Username = types.StringValue(res.GetUsername())
	data.Url = types.StringValue(res.GetUrl())
	data.Notes = types.StringValue(res.GetNotes())
	data.FolderId = types.StringValue(res.GetGroupId())
	data.Created = types.StringValue("Not implemented")
	data.Modified = types.StringValue("Not implemented")
	data.Expires = types.StringValue(res.GetExpires())
	data.Tags = d.fetchTags(res.Tags)

	pwdres, httpres, err := client.DefaultAPI.GetV6CredentialPasswordByID(*d.ctx, credential_id).Execute()
	if err != nil {
		resp.Diagnostics.AddError("failure to invoke API: ", err.Error())
		return
	}
	if httpres.StatusCode != 200 {
		resp.Diagnostics.AddError("Got an unexpected response code", fmt.Sprintf("Got an unexpected response code %v", httpres.StatusCode))
		return
	}

	sanitypassword, err := strconv.Unquote(pwdres) // used to remove the quotes and escape characters from the password
	if err != nil {
		sanitypassword = pwdres
	}
	data.Password = types.StringValue(sanitypassword)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
