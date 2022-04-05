package looker

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	apiclient "github.com/looker-open-source/sdk-codegen/go/sdk/v4"
)

func TestAcc_RoleGroups(t *testing.T) {
	name1 := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: roleGroupsConfig(name1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("looker_role_groups.test", "group_ids.#", "1"),
				),
			},
			{
				ResourceName:      "looker_role_groups.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
		CheckDestroy: testAccCheckRoleGroupsDestroy,
	})
}

func testAccCheckRoleGroupsDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*apiclient.LookerSDK)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "looker_role_groups" {
			continue
		}

		roleGroupsID := rs.Primary.ID

		groups, err := client.RoleGroups(roleGroupsID, "", nil)
		if err != nil {
			if strings.Contains(err.Error(), "404") {
				return nil // successfully destroyed
			}
			return err
		}

		if len(groups) != 0 {
			return fmt.Errorf("role_groups '%s' still exists", rs.Primary.ID)
		}

	}

	return nil
}

func roleGroupsConfig(name string) string {
	return fmt.Sprintf(`
	resource "looker_group" "test" {
		name = "%s"
	}
	resource "looker_model_set" "test" {
		name = "%s"
		models = ["test"]
	}
	resource "looker_permission_set" "test" {
		name = "%s"
		permissions = ["test"]
	}
	resource "looker_role" "test" {
		name = "%s"
		permission_set_id = looker_permission_set.test.id
		model_set_id = looker_model_set.test.id
	}
	resource "looker_role_groups" "test" {
		role_id   = looker_role.test.id
		group_ids = [looker_group.test.id]
	}
	`, name, name, name, name)
}
