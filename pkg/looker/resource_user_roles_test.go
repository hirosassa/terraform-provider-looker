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

func TestAcc_UserRoles(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
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
		CheckDestroy: testAccCheckUserRoleDestroy,
	})
}

func testAccCheckUserRoleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*apiclient.LookerSDK)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "looker_user_role" {
			continue
		}

		userID := rs.Primary.ID

		request := apiclient.RequestUserRoles{UserId: userID}

		userRoles, err := client.UserRoles(request, nil)
		if err != nil {
			if strings.Contains(err.Error(), "404") {
				return nil // successfully destroyed
			}
			return err
		}

		if len(userRoles) != 0 {
			return fmt.Errorf("user_role '%s' still exists", rs.Primary.ID)
		}

	}

	return nil
}

func userRolesConfig(name string) string {
	return fmt.Sprintf(`
	resource "looker_user" "user_role_test" {
        first_name = "%s"
        last_name  = "%s"
		email      = "%s@example.com"
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
	`, name, name, name, name, name, name)
}
