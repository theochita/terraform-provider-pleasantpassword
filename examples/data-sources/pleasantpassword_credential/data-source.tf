data "pleasantpassword_folder_root" "root_folder_id" {
}

resource "pleasantpassword_folder" "create_folder" {
  name     = "example_folder"
  parentid = data.pleasantpassword_folder_root.root_folder_id.id
  notes    = "example notes"

}

resource "pleasantpassword_credential" "cred1" {
  name     = "example_credential1"
  groupid  = pleasantpassword_folder.create_folder.id
  password = "example_password1"
  notes    = "example notes"
  username = "example_username1"


}

data "pleasantpassword_credential" "get_credential" {
  credential_id = pleasantpassword_credential.cred1.id

}