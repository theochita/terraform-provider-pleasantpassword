// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSearchDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccSearchDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.pleasantpassword_search.search_test", "credentials.#", "2"),
					resource.TestCheckResourceAttr("data.pleasantpassword_search.search_test", "folders.#", "2"),
				),
			},
		},
	})
}

const testAccSearchDataSourceConfig = `



data "pleasantpassword_folder_root" "root_folder_id_test" {
}

resource "pleasantpassword_folder" "create_folder" {
	name = "acctest_folder"
	parentid = data.pleasantpassword_folder_root.root_folder_id_test.id
	notes = "acctest notes"
   
 }

 resource "pleasantpassword_folder" "create_folder1" {
	name = "acctest_folder1"
	parentid = data.pleasantpassword_folder_root.root_folder_id_test.id
	notes = "acctest notes"
   
 }
resource "pleasantpassword_credential" "cred1_test" {
	name = "acctest_credential1"
	folderid =  pleasantpassword_folder.create_folder.id
	password = "acctest_password1"
	notes = "acctest notes"
	username = "acctest_username1"
	
   
 }

 resource "pleasantpassword_credential" "cred2_test" {
	name = "acctest_credential2"
	folderid =  pleasantpassword_folder.create_folder.id
	password = "acctest_password2"
	notes = "acctest notes"
	username = "acctest_username2"
	
   
 }

 data "pleasantpassword_search" "search_test" {
	search = "acctest"

	depends_on = [
		pleasantpassword_credential.cred1_test,
		pleasantpassword_credential.cred2_test
	  ]
   
 }




`
