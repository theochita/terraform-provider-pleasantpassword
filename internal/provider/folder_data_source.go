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
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &FolderDataSource{}

func NewFolderDataSource() datasource.DataSource {
	return &FolderDataSource{}
}

// ExampleDataSource defines the data source implementation.
type FolderDataSource struct {
	client *PPSClient.APIClient
	ctx    *context.Context
}

// ExampleDataSourceModel describes the data source data model.
type FolderDataSourceModel struct {
	Id          types.String      `tfsdk:"id"`
	FolderID    types.String      `tfsdk:"folderid"`
	Name        types.String      `tfsdk:"name"`
	ParentID    types.String      `tfsdk:"parentid"`
	Credentials []Credential      `tfsdk:"credentials"`
	Children    []CredentialGroup `tfsdk:"children"`
	Tags        []Tag             `tfsdk:"tags"`
	Notes       types.String      `tfsdk:"notes"`
	Created     types.String      `tfsdk:"created"`
	Modified    types.String      `tfsdk:"modified"`
	Expires     types.String      `tfsdk:"expires"`
}

type CredentialGroup struct {
	//
	//CustomUserFields map[string]interface{} `json:"CustomUserFields,omitempty"`
	//
	//CustomApplicationFields    map[string]interface{}    `json:"CustomApplicationFields,omitempty"`
	Children    []CredentialGroup `tfsdk:"children"`
	Credentials []Credential      `tfsdk:"credentials"`
	Tags        []Tag             `tfsdk:"tags"`
	//HasModifyEntriesAccess     *bool                     `json:"HasModifyEntriesAccess,omitempty"`
	//HasViewEntryContentsAccess *bool                     `json:"HasViewEntryContentsAccess,omitempty"`
	//CommentPrompts             *V6CommentPromptResult    `json:"CommentPrompts,omitempty"`
	Id       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	ParentId types.String `tfsdk:"parentid"`
	Notes    types.String `tfsdk:"notes"`
	Created  types.String `tfsdk:"created"`
	Modified types.String `tfsdk:"modified"`
	Expires  types.String `tfsdk:"expires"`
}

type Credential struct {
	//
	//	CustomUserFields map[string]interface{} `json:"CustomUserFields,omitempty"`
	//
	//	CustomApplicationFields    map[string]interface{} `json:"CustomApplicationFields,omitempty"`
	Tags []Tag        `tfsdk:"tags"`
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`

	Username types.String `tfsdk:"username"`
	Url      types.String `tfsdk:"url"`
	Notes    types.String `tfsdk:"notes"`
	GroupId  types.String `tfsdk:"groupid"`
	Created  types.String `tfsdk:"created"`
	Modified types.String `tfsdk:"modified"`
	Expires  types.String `tfsdk:"expires"`
}

type Tag struct {
	Name types.String `tfsdk:"name"`
}

func (d FolderDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_folder"
}

func (d *FolderDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Example data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Example configurable attribute",
				Optional:            true,
				Computed:            true,
			},
			"folderid": schema.StringAttribute{
				MarkdownDescription: "Example identifier",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Example identifier",
				Computed:            true,
			},
			"parentid": schema.StringAttribute{
				MarkdownDescription: "Example identifier",
				Computed:            true,
			},
			"notes": schema.StringAttribute{
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
				},
			},
			"children": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "Example configurable attribute",
							Optional:            true,
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Example identifier",
							Computed:            true,
						},
						"parentid": schema.StringAttribute{
							MarkdownDescription: "Example identifier",
							Computed:            true,
						},
						"notes": schema.StringAttribute{
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
						"children": schema.ListNestedAttribute{
							MarkdownDescription: "Example configurable attribute",
							Optional:            true,
							Computed:            true,
							NestedObject:        schema.NestedAttributeObject{},
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
							},
						}},
				},
			},
		},
	}
}

func (d *FolderDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *FolderDataSource) fetchTags(res []PPSClient.V6TagResult) []Tag {
	var tags = []Tag{}
	for _, v := range res {
		tag := Tag{}
		tag.Name = types.StringValue(*v.Name)
		tags = append(tags, tag)
	}

	return tags

}

func (d *FolderDataSource) fetchCredentials(res []PPSClient.V6CredentialResult) []Credential {
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

func (d *FolderDataSource) fetchChildren(res []PPSClient.V6CredentialGroupOutput) []CredentialGroup {

	var children = []CredentialGroup{}
	for _, v := range res {
		child := CredentialGroup{}
		child.Id = types.StringValue(v.GetId())
		child.Name = types.StringValue(v.GetName())
		child.ParentId = types.StringValue(v.GetParentId())
		child.Notes = types.StringValue(v.GetNotes())
		child.Created = types.StringValue("Not implemented")
		child.Modified = types.StringValue("Not implemented")
		child.Expires = types.StringValue(v.GetExpires())

		child.Tags = d.fetchTags(v.Tags)

		child.Credentials = d.fetchCredentials(v.Credentials)

		children = append(children, child)

	}

	return children

}

func (d *FolderDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data FolderDataSourceModel

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

	folderid := data.FolderID.ValueString()

	client := d.client

	res, httpres, err := client.DefaultAPI.GetV6FoldersByID(*d.ctx, folderid).Execute()
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

	data.Id = types.StringValue(*res.Id)
	data.Name = types.StringValue(*res.Name)
	data.ParentID = types.StringValue(*res.ParentId)
	data.Notes = types.StringValue(*res.Notes)
	data.Created = types.StringValue("Not implemented")
	data.Modified = types.StringValue("Not implemented")
	data.Expires = types.StringValue(res.GetExpires())

	data.Tags = d.fetchTags(res.Tags)

	data.Credentials = d.fetchCredentials(res.Credentials)
	data.Children = d.fetchChildren(res.Children)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
