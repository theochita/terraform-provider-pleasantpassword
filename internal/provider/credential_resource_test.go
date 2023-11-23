// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCredentialResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccCredentialResourceConfig("one"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pleasantpassword_credential.cred1_test", "name", "acctest_credentialone"),
					resource.TestCheckResourceAttr("pleasantpassword_credential.cred1_test", "password", "acctest_passwordone"),
				),
			},

			// Update and Read testing
			{
				Config: testAccCredentialResourceConfig("two"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pleasantpassword_credential.cred1_test", "name", "acctest_credentialtwo"),
					resource.TestCheckResourceAttr("pleasantpassword_credential.cred1_test", "password", "acctest_passwordtwo"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccCredentialResourceConfig(configurableAttribute string) string {
	return fmt.Sprintf(`

data "pleasantpassword_folder_root" "get_root_folder" {
}

resource "pleasantpassword_folder" "create_folder" {
	name = "acctest_folder"
	parent_id = data.pleasantpassword_folder_root.get_root_folder.id
	notes = "testnotes"
 }


 resource "pleasantpassword_credential" "cred1_test" {
	name = "acctest_credential%[1]s"
	folder_id =  pleasantpassword_folder.create_folder.id
	password = "acctest_password%[1]s"
	notes = "acctest notes"
	username = "acctest_username1"
	
   
 }


`, configurableAttribute)
}
