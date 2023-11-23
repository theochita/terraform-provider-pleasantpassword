data "pleasantpassword_folder_root" "get_root_folder" {
}

resource "pleasantpassword_folder" "create_folder" {
  name     = "example_folder"
  parentid = data.pleasantpassword_folder_root.get_root_folder.id
  notes    = " example notes"
}