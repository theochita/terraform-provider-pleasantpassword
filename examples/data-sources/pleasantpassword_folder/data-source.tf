data "pleasantpassword_folder_root" "root_folder_id" {
}


data "pleasantpassword_folder" "fetch_root_folder" {
  folderid = data.pleasantpassword_folder_root.root_folder_id.id

}