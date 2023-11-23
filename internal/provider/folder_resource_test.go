// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFolderResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccFolderResourceConfig("one"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pleasantpassword_folder.create_folder", "name", "acctest_folderone"),
				),
			},

			// Update and Read testing
			{
				Config: testAccFolderResourceConfig("two"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pleasantpassword_folder.create_folder", "name", "acctest_foldertwo"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccFolderResourceConfig(configurableAttribute string) string {
	return fmt.Sprintf(`

	



data "pleasantpassword_folder_root" "get_root_folder" {
}

resource "pleasantpassword_folder" "create_folder" {
	name = "acctest_folder%s"
	parent_id = data.pleasantpassword_folder_root.get_root_folder.id
	notes = "testnotes"
   
 }
`, configurableAttribute)
}
