package looker

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_UserRoles(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: userRolesConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("looker_user_roles.user_role_test", "role_ids.#", "1"),
				),
			},
			{
				ResourceName:      "looker_user_roles.user_role_test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func userRolesConfig(name string) string {
	return fmt.Sprintf(`
	resource "looker_user" "user_role_test" {
		email = "%s@example.com"
	}
	resource "looker_model_set" "user_role_test" {
		name = "%s"
		models = ["test"]
	}
	resource "looker_permission_set" "user_role_test" {
		name = "%s"
		permissions = ["access_data"]
	}
	resource "looker_role" "user_role_test" {
		name = "%s"
		permission_set_id = looker_permission_set.user_role_test.id
		model_set_id = looker_model_set.user_role_test.id
	}
	resource "looker_user_roles" "user_role_test" {
		user_id  = looker_user.user_role_test.id
		role_ids = [looker_role.user_role_test.id]
	}
	`, name, name, name, name)
}
