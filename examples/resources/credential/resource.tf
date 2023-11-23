data "pleasantpassword_folder_root" "get_root_folder" {
}

resource "pleasantpassword_folder" "create_folder" {
  name     = "example_folder"
  parentid = data.pleasantpassword_folder_root.get_root_folder.id
  notes    = "testnotes"
}


resource "pleasantpassword_credential" "cred1" {
  name     = "example_credential"
  folderid = pleasantpassword_folder.create_folder.id
  password = "example_password"
  notes    = "example notes"
  username = "example_username1"


}