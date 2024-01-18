---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "pleasantpassword_search Data Source - terraform-provider-pleasantpassword"
subcategory: ""
description: |-
  The search data source can be used to search for credentials and folders.
---

# pleasantpassword_search (Data Source)

The `search` data source can be used to search for credentials and folders.

## Example Usage

```terraform
data "pleasantpassword_search" "search_cred1" {
  search = "example_cred1"

}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `search` (String) The search query for credentials and folders.

### Read-Only

- `credentials` (Attributes List) (see [below for nested schema](#nestedatt--credentials))
- `folders` (Attributes List) (see [below for nested schema](#nestedatt--folders))

<a id="nestedatt--credentials"></a>
### Nested Schema for `credentials`

Read-Only:

- `folder_id` (String) The identifier of the folder that the credential belongs to.
- `id` (String) The identifier of the credential.
- `name` (String) The name of the credential.
- `notes` (String) The notes of the credential.
- `path` (String) The path of the credential.
- `url` (String) The URL of the credential.
- `username` (String) The username of the credential.


<a id="nestedatt--folders"></a>
### Nested Schema for `folders`

Read-Only:

- `fullpath` (String) The full path of the folder.
- `id` (String) The identifier of the folder.
- `name` (String) The name of the folder.