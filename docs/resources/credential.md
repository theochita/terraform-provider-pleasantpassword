---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "pleasantpassword_credential Resource - terraform-provider-pleasantpassword"
subcategory: ""
description: |-
  The credential resource allows you to create and manage credentials in Pleasant Password Server.
---

# pleasantpassword_credential (Resource)

The `credential` resource allows you to create and manage credentials in Pleasant Password Server.

## Example Usage

```terraform
data "pleasantpassword_folder_root" "get_root_folder" {
}

resource "pleasantpassword_folder" "create_folder" {
  name      = "example_folder"
  parent_id = data.pleasantpassword_folder_root.get_root_folder.id
  notes     = "testnotes"
}


resource "pleasantpassword_credential" "cred1" {
  name      = "example_credential"
  folder_id = pleasantpassword_folder.create_folder.id
  password  = "example_password"
  notes     = "example notes"
  username  = "example_username1"


}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `folder_id` (String) The folder ID where the credential is stored.
- `name` (String) The name of the credential.

### Optional

- `expires` (String) The expiration date of the credential.
- `notes` (String) Additional notes for the credential.
- `password` (String) The password associated with the credential.
- `url` (String) The URL associated with the credential.
- `username` (String) The username associated with the credential.

### Read-Only

- `created` (String) The creation timestamp of the credential.
- `id` (String) The unique identifier of the credential.
- `modified` (String) The last modification timestamp of the credential.
