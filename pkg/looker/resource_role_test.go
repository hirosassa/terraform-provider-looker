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

func TestAcc_Role(t *testing.T) {
	name1 := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: roleConfig(name1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("looker_role.role_test", "name", name1),
				),
			},
			{
				ResourceName:      "looker_role.role_test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
		CheckDestroy: testAccCheckRoleDestroy,
	})
}

func testAccCheckRoleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*apiclient.LookerSDK)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "looker_role" {
			continue
		}

		roleID := rs.Primary.ID

		role, err := client.Role(roleID, nil)
		if err != nil {
			if strings.Contains(err.Error(), "404") {
				return nil // successfully destroyed
			}
			return err
		}

		if *role.Name == rs.Primary.Attributes["name"] {
			return fmt.Errorf("role '%s' still exists", rs.Primary.ID)
		}

	}

	return nil
}

func roleConfig(name string) string {
	return fmt.Sprintf(`
	resource "looker_model_set" "role_test" {
		name = "%s"
		models = ["test"]
	}
	resource "looker_permission_set" "role_test" {
		name = "%s"
		permissions = ["test"]
	}
	resource "looker_role" "role_test" {
		name = "%s"
		permission_set_id = looker_permission_set.role_test.id
		model_set_id = looker_model_set.role_test.id
	}
	`, name, name, name)
}
