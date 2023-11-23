package models

import "github.com/hashicorp/terraform-plugin-framework/types"

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
	Folderid types.String `tfsdk:"folderid"`
	Created  types.String `tfsdk:"created"`
	Modified types.String `tfsdk:"modified"`
	Expires  types.String `tfsdk:"expires"`
}

type Tag struct {
	Name types.String `tfsdk:"name"`
}
