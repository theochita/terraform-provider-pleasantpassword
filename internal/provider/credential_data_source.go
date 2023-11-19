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
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &CredentialDataSource{}

func NewCredentialDataSource() datasource.DataSource {
	return &CredentialDataSource{}
}

// ExampleDataSource defines the data source implementation.
type CredentialDataSource struct {
	client *PPSClient.APIClient
	ctx    *context.Context
}

// ExampleDataSourceModel describes the data source data model.
type CredentialDataSourceModel struct {
	Id           types.String `tfsdk:"id"`
	CredentialID types.String `tfsdk:"credential_id"`
	Tags         []Tag        `tfsdk:"tags"`
	Name         types.String `tfsdk:"name"`
	Username     types.String `tfsdk:"username"`
	Password     types.String `tfsdk:"password"`
	Url          types.String `tfsdk:"url"`
	Notes        types.String `tfsdk:"notes"`
	GroupId      types.String `tfsdk:"groupid"`
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
		MarkdownDescription: "Example data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Example identifier",
				Computed:            true,
			},
			"credential_id": schema.StringAttribute{
				MarkdownDescription: "Example identifier",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Example identifier",
				Computed:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Example identifier",
				Computed:            true,
			},
			"password": schema.StringAttribute{
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
			"groupid": schema.StringAttribute{
				MarkdownDescription: "Example identifier",
				Computed:            true,
			},
			"created": schema.StringAttribute{
				MarkdownDescription: "Example identifier",
				Computed:            true,
			},
			"modified": schema.StringAttribute{
				MarkdownDescription: "Example identifier",
				Computed:            true,
			},
			"expires": schema.StringAttribute{
				MarkdownDescription: "Example identifier",
				Computed:            true,
			},

			"tags": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "Example identifier",
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

func (d *CredentialDataSource) fetchTags(res []PPSClient.V6TagResult) []Tag {
	var tags = []Tag{}
	for _, v := range res {
		tag := Tag{}
		tag.Name = types.StringValue(v.GetName())
		tags = append(tags, tag)
	}

	return tags

}

func (d *CredentialDataSource) fetchCredentials(res []PPSClient.V6CredentialResult) []Credential {
	var creds = []Credential{}
	for _, v := range res {
		cred := Credential{}
		cred.Id = types.StringValue(v.GetId())
		cred.Name = types.StringValue(v.GetName())
		cred.Username = types.StringValue(v.GetUsername())
		cred.Url = types.StringValue(v.GetUrl())
		cred.Notes = types.StringValue(v.GetNotes())
		cred.GroupId = types.StringValue(v.GetGroupId())
		cred.Created = types.StringValue("Not implemented")
		cred.Modified = types.StringValue("Not implemented")
		cred.Expires = types.StringValue(v.GetExpires())

		cred.Tags = d.fetchTags(v.Tags)

		creds = append(creds, cred)

	}

	return creds

}

func (d *CredentialDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CredentialDataSourceModel

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
	data.GroupId = types.StringValue(res.GetGroupId())
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
