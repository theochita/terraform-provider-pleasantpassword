// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFolderDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccFolderDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.pleasantpassword_folder.fetch_root_folder_test", "parent_id", "00000000-0000-0000-0000-000000000000"),
					resource.TestCheckResourceAttrPair("data.pleasantpassword_folder.fetch_root_folder_test", "folder_id", "data.pleasantpassword_folder_root.root_folder_id_test", "id"),
				),
			},
		},
	})
}

const testAccFolderDataSourceConfig = `



data "pleasantpassword_folder_root" "root_folder_id_test" {
}


data "pleasantpassword_folder" "fetch_root_folder_test" {
   folder_id = data.pleasantpassword_folder_root.root_folder_id_test.id
  
}
`
