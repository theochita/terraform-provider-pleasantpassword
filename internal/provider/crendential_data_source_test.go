// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCredentialDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccCredentialDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.pleasantpassword_credential.get_credential_test", "name", "acctest_credential1"),
					resource.TestCheckResourceAttr("data.pleasantpassword_credential.get_credential_test", "password", "acctest_password1"),
					//resource.TestCheckResourceAttrPair("data.pleasantpassword_folder.fetch_root_folder_test", "parentid", "data.pleasantpassword_folder_root.root_folder_id_test", "id"),
				),
			},
		},
	})
}

const testAccCredentialDataSourceConfig = `



data "pleasantpassword_folder_root" "root_folder_id_test" {
}

resource "pleasantpassword_folder" "create_folder" {
	name = "acctest_folder"
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

 data "pleasantpassword_credential" "get_credential_test" {
	credential_id = pleasantpassword_credential.cred1_test.id
   
 }




`
