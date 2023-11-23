// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	PPSClient "github.com/theochita/go-pleasant-password"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &CredentialResource{}
var _ resource.ResourceWithImportState = &CredentialResource{}

func NewCredentialResource() resource.Resource {
	return &CredentialResource{}
}

type CredentialResource struct {
	client *PPSClient.APIClient
	ctx    *context.Context
}

type CredentialResourceModel struct {
	Id types.String `tfsdk:"id"`
	//Tags     []Tag        `tfsdk:"tags"`
	Name     types.String `tfsdk:"name"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
	Url      types.String `tfsdk:"url"`
	Notes    types.String `tfsdk:"notes"`
	FolderId types.String `tfsdk:"folderid"`
	Created  types.String `tfsdk:"created"`
	Modified types.String `tfsdk:"modified"`
	Expires  types.String `tfsdk:"expires"`
}

func (r *CredentialResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_credential"
}

func (r *CredentialResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Example resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Example identifier",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Example identifier",
				Required:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Example identifier",
				Computed:            true,
				Optional:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "Example identifier",
				Computed:            true,
				Optional:            true,
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "Example identifier",
				Computed:            true,
				Optional:            true,
			},
			"notes": schema.StringAttribute{
				MarkdownDescription: "Example identifier",
				Computed:            true,
				Optional:            true,
			},
			"folderid": schema.StringAttribute{
				MarkdownDescription: "Example identifier",
				Required:            true,
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
				Optional:            true,
			},
		},
	}
}

func (r *CredentialResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = &providerclient.Client
	r.ctx = &providerclient.Ctx

}

func (r *CredentialResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CredentialResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	param := PPSClient.NewV6CredentialInputWithDefaults()
	param.Name = data.Name.ValueStringPointer()
	param.Notes = data.Notes.ValueStringPointer()
	param.GroupId = data.FolderId.ValueStringPointer()
	param.Username = data.Username.ValueStringPointer()
	param.Password = data.Password.ValueStringPointer()
	param.Url = data.Url.ValueStringPointer()

	// expire and tags not implemented

	res, httpres, err := r.client.DefaultAPI.PostV6Credentials(*r.ctx).V6CredentialInput(*param).Execute()
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

	sanityresult, err := strconv.Unquote(res)
	if err != nil {
		sanityresult = res
	}

	data.Id = types.StringValue(sanityresult)
	data.Name = types.StringValue(param.GetName())
	data.Notes = types.StringValue(param.GetNotes())
	data.FolderId = types.StringValue(param.GetGroupId())
	data.Username = types.StringValue(param.GetUsername())
	data.Password = types.StringValue(param.GetPassword())
	data.Url = types.StringValue(param.GetUrl())
	data.Created = types.StringValue("Not implemented")
	data.Modified = types.StringValue("Not implemented")
	data.Expires = types.StringValue("Not implemented")
	//data.Tags = types.ListValue(param.GetTags())

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CredentialResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CredentialResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	res, httpres, err := r.client.DefaultAPI.GetV6CredentialsByID(*r.ctx, data.Id.ValueString()).Execute()

	if err != nil {
		resp.State.RemoveResource(ctx)
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
	data.Expires = types.StringValue("Not implemented")
	//data.Tags = r.fetchTags(res.Tags)

	pwdres, httpres, err := r.client.DefaultAPI.GetV6CredentialPasswordByID(*r.ctx, data.Id.ValueString()).Execute()
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

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CredentialResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data CredentialResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	param := PPSClient.NewV6CredentialInputWithDefaults()
	param.Name = data.Name.ValueStringPointer()
	param.Notes = data.Notes.ValueStringPointer()
	param.GroupId = data.FolderId.ValueStringPointer()
	param.Username = data.Username.ValueStringPointer()
	param.Password = data.Password.ValueStringPointer()
	param.Url = data.Url.ValueStringPointer()

	httpres, err := r.client.DefaultAPI.PatchV6CredentialsByID(*r.ctx, data.Id.ValueString()).V6CredentialInput(*param).Execute()

	if err != nil {
		resp.Diagnostics.AddError("failure to invoke API: ", err.Error())
		return
	}
	if httpres.StatusCode != 204 {
		resp.Diagnostics.AddError("Got an unexpected response code", fmt.Sprintf("Got an unexpected response code %v", httpres.StatusCode))
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	data.Name = types.StringValue(param.GetName())
	data.Notes = types.StringValue(param.GetNotes())
	data.FolderId = types.StringValue(param.GetGroupId())
	data.Username = types.StringValue(param.GetUsername())
	data.Password = types.StringValue(param.GetPassword())
	data.Url = types.StringValue(param.GetUrl())
	data.Created = types.StringValue("Not implemented")
	data.Modified = types.StringValue("Not implemented")
	data.Expires = types.StringValue("Not implemented")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CredentialResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CredentialResourceModel

	//Delete removes current folder and all subfolders and credentials

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	httpres, err := r.client.DefaultAPI.DeleteV6CredentialsByID(*r.ctx, data.Id.ValueString()).Execute()

	if err != nil {
		resp.Diagnostics.AddError("failure to invoke API: ", err.Error())
		return
	}
	if httpres.StatusCode != 204 {
		resp.Diagnostics.AddError("Got an unexpected response code", fmt.Sprintf("Got an unexpected response code %v", httpres.StatusCode))
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *CredentialResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// not implemented
}
