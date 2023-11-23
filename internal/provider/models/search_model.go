package models

import "github.com/hashicorp/terraform-plugin-framework/types"

type V6CredentialGroupSearchResult struct {
	Id       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	FullPath types.String `tfsdk:"fullpath"`
}

type V6CredentialSearchResult struct {
	Id       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Username types.String `tfsdk:"username"`
	Url      types.String `tfsdk:"url"`
	Notes    types.String `tfsdk:"notes"`
	FolderId types.String `tfsdk:"folder_id"`
	Path     types.String `tfsdk:"path"`
}
