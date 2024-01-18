---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "pleasantpassword_folder Data Source - terraform-provider-pleasantpassword"
subcategory: ""
description: |-
  The folder data source can be used to access information about a folder.
---

# pleasantpassword_folder (Data Source)

The `folder` data source can be used to access information about a folder.

## Example Usage

```terraform
data "pleasantpassword_folder_root" "root_folder_id" {
}


data "pleasantpassword_folder" "fetch_root_folder" {
  folder_id = data.pleasantpassword_folder_root.root_folder_id.id

}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `folder_id` (String) Required: Id of the folder

### Read-Only

- `children` (Attributes List) (see [below for nested schema](#nestedatt--children))
- `created` (String) Creation timestamp of the folder
- `credentials` (Attributes List) (see [below for nested schema](#nestedatt--credentials))
- `expires` (String) Expiration timestamp of the folder
- `id` (String) Identifier of the folder
- `modified` (String) Last modified timestamp of the folder
- `name` (String) Name of the folder
- `notes` (String) Notes for the folder
- `parent_id` (String) Identifier of the parent folder
- `tags` (Attributes List) (see [below for nested schema](#nestedatt--tags))

<a id="nestedatt--children"></a>
### Nested Schema for `children`

Read-Only:

- `children` (Attributes List) Empty list of child folders (see [below for nested schema](#nestedatt--children--children))
- `created` (String) Creation timestamp of the child folder
- `credentials` (Attributes List) (see [below for nested schema](#nestedatt--children--credentials))
- `expires` (String) Expiration timestamp of the child folder
- `id` (String) Identifier of the child folder
- `modified` (String) Last modified timestamp of the child folder
- `name` (String) Name of the child folder
- `notes` (String) Notes for the child folder
- `parent_id` (String) Identifier of the  parent folder
- `tags` (Attributes List) (see [below for nested schema](#nestedatt--children--tags))

<a id="nestedatt--children--children"></a>
### Nested Schema for `children.children`


<a id="nestedatt--children--credentials"></a>
### Nested Schema for `children.credentials`

Read-Only:

- `created` (String) The creation date of the credential
- `expires` (String) The expiration date of the credential
- `folder_id` (String) The folder ID of the credential
- `id` (String) The unique identifier of the credential
- `modified` (String) The modification date of the credential
- `name` (String) The name of the credential
- `notes` (String) The notes of the credential
- `tags` (Attributes List) (see [below for nested schema](#nestedatt--children--credentials--tags))
- `url` (String) The URL of the credential
- `username` (String) The username of the credential

<a id="nestedatt--children--credentials--tags"></a>
### Nested Schema for `children.credentials.tags`

Read-Only:

- `name` (String) The name of the tag



<a id="nestedatt--children--tags"></a>
### Nested Schema for `children.tags`

Read-Only:

- `name` (String) Name of the tag



<a id="nestedatt--credentials"></a>
### Nested Schema for `credentials`

Read-Only:

- `created` (String) The creation date of the credential
- `expires` (String) The expiration date of the credential
- `folder_id` (String) The folder ID of the credential
- `id` (String) The unique identifier of the credential
- `modified` (String) The modification date of the credential
- `name` (String) The name of the credential
- `notes` (String) The notes of the credential
- `tags` (Attributes List) (see [below for nested schema](#nestedatt--credentials--tags))
- `url` (String) The URL of the credential
- `username` (String) The username of the credential

<a id="nestedatt--credentials--tags"></a>
### Nested Schema for `credentials.tags`

Read-Only:

- `name` (String) The name of the tag



<a id="nestedatt--tags"></a>
### Nested Schema for `tags`

Read-Only:

- `name` (String) Name of the tag