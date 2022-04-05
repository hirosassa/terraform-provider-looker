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

func TestAcc_PermissionSet(t *testing.T) {
	name1 := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	name2 := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: permissionSetConfig(name1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("looker_permission_set.test", "name", name1),
					resource.TestCheckResourceAttr("looker_permission_set.test", "permissions.#", "1"),
				),
			},
			{
				Config: permissionSetConfig(name2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("looker_permission_set.test", "name", name2),
					resource.TestCheckResourceAttr("looker_permission_set.test", "permissions.#", "1"),
				),
			},
			{
				ResourceName:      "looker_permission_set.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
		CheckDestroy: testAccCheckPermissionSetDestroy,
	})
}

func testAccCheckPermissionSetDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*apiclient.LookerSDK)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "looker_permission_set" {
			continue
		}

		permissionSetID := rs.Primary.ID

		permissionSet, err := client.PermissionSet(permissionSetID, "", nil)
		if err != nil {
			if strings.Contains(err.Error(), "404") {
				return nil // successfully destroyed
			}
			return err
		}

		if *permissionSet.Name == rs.Primary.Attributes["name"] {
			return fmt.Errorf("permission_set '%s' still exists", rs.Primary.ID)
		}

	}

	return nil
}

func permissionSetConfig(name string) string {
	return fmt.Sprintf(`
	resource "looker_permission_set" "test" {
		name = "%s"
		permissions = ["test"]
	}
	`, name)
}
