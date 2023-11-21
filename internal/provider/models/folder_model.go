package models

import "github.com/hashicorp/terraform-plugin-framework/types"

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
