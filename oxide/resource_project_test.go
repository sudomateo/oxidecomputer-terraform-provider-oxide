// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package oxide

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	oxideSDK "github.com/oxidecomputer/oxide.go/oxide"
)

func TestAccResourceProject_full(t *testing.T) {
	resourceName := "oxide_project.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		CheckDestroy:             testAccProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testResourceProjectConfig,
				Check:  checkResourceProject(resourceName),
			},
			{
				Config: testResourceProjectUpdateConfig,
				Check:  checkResourceProjectUpdate(resourceName),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var testResourceProjectConfig = `
resource "oxide_project" "test" {
	description       = "a test project"
	name              = "terraform-acc-myproject"
	timeouts = {
		read   = "1m"
		create = "3m"
		delete = "2m"
		update = "4m"
	}
  }
`

func checkResourceProject(resourceName string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc([]resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(resourceName, "id"),
		resource.TestCheckResourceAttr(resourceName, "description", "a test project"),
		resource.TestCheckResourceAttr(resourceName, "name", "terraform-acc-myproject"),
		resource.TestCheckResourceAttrSet(resourceName, "time_created"),
		resource.TestCheckResourceAttrSet(resourceName, "time_modified"),
		resource.TestCheckResourceAttr(resourceName, "timeouts.read", "1m"),
		resource.TestCheckResourceAttr(resourceName, "timeouts.delete", "2m"),
		resource.TestCheckResourceAttr(resourceName, "timeouts.create", "3m"),
		resource.TestCheckResourceAttr(resourceName, "timeouts.update", "4m"),
	}...)
}

var testResourceProjectUpdateConfig = `
resource "oxide_project" "test" {
	description       = "a new description for project"
	name              = "terraform-acc-myproject2"
  }
`

func checkResourceProjectUpdate(resourceName string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc([]resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(resourceName, "id"),
		resource.TestCheckResourceAttr(resourceName, "description", "a new description for project"),
		resource.TestCheckResourceAttr(resourceName, "name", "terraform-acc-myproject2"),
		resource.TestCheckResourceAttrSet(resourceName, "time_created"),
		resource.TestCheckResourceAttrSet(resourceName, "time_modified"),
	}...)
}

func testAccProjectDestroy(s *terraform.State) error {
	client, err := newTestClient()
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "oxide_project" {
			continue
		}

		res, err := client.ProjectView(oxideSDK.ProjectViewParams{Project: "terraform-acc-myproject2"})
		if err != nil && is404(err) {
			continue
		}
		return fmt.Errorf("project (%v) still exists", &res.Name)
	}

	return nil
}
