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
	"github.com/theochita/terraform-provider-pleasantpassword/internal/provider/models"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &FolderDataSource{}

func NewFolderDataSource() datasource.DataSource {
	return &FolderDataSource{}
}

type FolderDataSource struct {
	client *PPSClient.APIClient
	ctx    *context.Context
}

type FolderDataSourceModel struct {
	Id          types.String             `tfsdk:"id"`
	FolderID    types.String             `tfsdk:"folder_id"`
	Name        types.String             `tfsdk:"name"`
	ParentID    types.String             `tfsdk:"parent_id"`
	Credentials []models.Credential      `tfsdk:"credentials"`
	Children    []models.CredentialGroup `tfsdk:"children"`
	Tags        []models.Tag             `tfsdk:"tags"`
	Notes       types.String             `tfsdk:"notes"`
	Created     types.String             `tfsdk:"created"`
	Modified    types.String             `tfsdk:"modified"`
	Expires     types.String             `tfsdk:"expires"`
}

func (d FolderDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_folder"
}

func (d *FolderDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "The `folder` data source can be used to access information about a folder.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier of the folder",
				Computed:            true,
			},
			"folder_id": schema.StringAttribute{
				MarkdownDescription: "Required: Id of the folder",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the folder",
				Computed:            true,
			},
			"parent_id": schema.StringAttribute{
				MarkdownDescription: "Identifier of the parent folder",
				Computed:            true,
			},
			"notes": schema.StringAttribute{
				MarkdownDescription: "Notes for the folder",
				Computed:            true,
			},
			"created": schema.StringAttribute{
				MarkdownDescription: "Creation timestamp of the folder",
				Computed:            true,
			},
			"modified": schema.StringAttribute{
				MarkdownDescription: "Last modified timestamp of the folder",
				Computed:            true,
			},
			"expires": schema.StringAttribute{
				MarkdownDescription: "Expiration timestamp of the folder",
				Computed:            true,
			},

			"tags": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "Name of the tag",
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
							MarkdownDescription: "The unique identifier of the credential",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "The name of the credential",
							Computed:            true,
						},
						"username": schema.StringAttribute{
							MarkdownDescription: "The username of the credential",
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
						"folder_id": schema.StringAttribute{
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
				},
			},
			"children": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "Identifier of the child folder",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Name of the child folder",
							Computed:            true,
						},
						"parent_id": schema.StringAttribute{
							MarkdownDescription: "Identifier of the  parent folder",
							Computed:            true,
						},
						"notes": schema.StringAttribute{
							MarkdownDescription: "Notes for the child folder",
							Computed:            true,
						},
						"created": schema.StringAttribute{
							MarkdownDescription: "Creation timestamp of the child folder",
							Computed:            true,
						},
						"modified": schema.StringAttribute{
							MarkdownDescription: "Last modified timestamp of the child folder",
							Computed:            true,
						},
						"expires": schema.StringAttribute{
							MarkdownDescription: "Expiration timestamp of the child folder",
							Computed:            true,
						},
						"children": schema.ListNestedAttribute{
							MarkdownDescription: "Empty list of child folders",
							Computed:            true,
							NestedObject:        schema.NestedAttributeObject{},
						},
						"tags": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										MarkdownDescription: "Name of the tag",
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
										MarkdownDescription: "The unique identifier of the credential",
										Computed:            true,
									},
									"name": schema.StringAttribute{
										MarkdownDescription: "The name of the credential",
										Computed:            true,
									},
									"username": schema.StringAttribute{
										MarkdownDescription: "The username of the credential",
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
									"folder_id": schema.StringAttribute{
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

func (d *FolderDataSource) fetchTags(res []PPSClient.V6TagResult) []models.Tag {
	var tags = []models.Tag{}
	for _, v := range res {
		tag := models.Tag{}
		tag.Name = types.StringValue(v.GetName())
		tags = append(tags, tag)
	}

	return tags

}

func (d *FolderDataSource) fetchCredentials(res []PPSClient.V6CredentialResult) []models.Credential {
	var creds = []models.Credential{}
	for _, v := range res {
		cred := models.Credential{}
		cred.Id = types.StringValue(v.GetId())
		cred.Name = types.StringValue(v.GetName())
		cred.Username = types.StringValue(v.GetUsername())
		cred.Url = types.StringValue(v.GetUrl())
		cred.Notes = types.StringValue(v.GetNotes())
		cred.Folderid = types.StringValue(v.GetGroupId())
		cred.Created = types.StringValue("Not implemented")
		cred.Modified = types.StringValue("Not implemented")
		cred.Expires = types.StringValue(v.GetExpires())

		cred.Tags = d.fetchTags(v.Tags)

		creds = append(creds, cred)

	}

	return creds

}

func (d *FolderDataSource) fetchChildren(res []PPSClient.V6CredentialGroupOutput) []models.CredentialGroup {

	var children = []models.CredentialGroup{}
	for _, v := range res {
		child := models.CredentialGroup{}
		child.Id = types.StringValue(v.GetId())
		child.Name = types.StringValue(v.GetName())
		child.ParentId = types.StringValue(v.GetParentId())
		child.Notes = types.StringValue(v.GetNotes())
		child.Created = types.StringValue("Not implemented")
		child.Modified = types.StringValue("Not implemented")
		child.Expires = types.StringValue(v.GetExpires())

		child.Tags = d.fetchTags(v.GetTags())

		child.Credentials = d.fetchCredentials(v.GetCredentials())

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

	data.Id = types.StringValue(res.GetId())
	data.Name = types.StringValue(res.GetName())
	data.ParentID = types.StringValue(res.GetParentId())
	data.Notes = types.StringValue(res.GetNotes())
	data.Created = types.StringValue("Not implemented")
	data.Modified = types.StringValue("Not implemented")
	data.Expires = types.StringValue(res.GetExpires())

	data.Tags = d.fetchTags(res.GetTags())

	data.Credentials = d.fetchCredentials(res.GetCredentials())
	data.Children = d.fetchChildren(res.GetChildren())

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
