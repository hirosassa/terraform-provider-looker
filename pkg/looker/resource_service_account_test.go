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

func TestAcc_ServiceAccount(t *testing.T) {
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckServiceAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: serviceAccountConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("looker_service_account.test", "service_account_name", name),
					resource.TestCheckResourceAttr("looker_service_account.test", "is_disabled", "false"),
				),
			},
			{
				ResourceName:      "looker_service_account.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_ServiceAccountUpdate(t *testing.T) {
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	updatedName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckServiceAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: serviceAccountConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("looker_service_account.test", "service_account_name", name),
					resource.TestCheckResourceAttr("looker_service_account.test", "is_disabled", "false"),
				),
			},
			{
				Config: serviceAccountConfigDisabled(updatedName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("looker_service_account.test", "service_account_name", updatedName),
					resource.TestCheckResourceAttr("looker_service_account.test", "is_disabled", "true"),
				),
			},
		},
	})
}

func testAccCheckServiceAccountDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*apiclient.LookerSDK)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "looker_service_account" {
			continue
		}

		userID := rs.Primary.ID

		_, err := client.User(userID, "", nil)
		if err != nil {
			if strings.Contains(err.Error(), "404") {
				return nil // successfully destroyed
			}
			return err
		}

		return fmt.Errorf("service account '%s' still exists", userID)
	}
	return nil
}

func serviceAccountConfig(name string) string {
	return fmt.Sprintf(`
resource "looker_service_account" "test" {
  service_account_name = "%s"
}
`, name)
}

func serviceAccountConfigDisabled(name string) string {
	return fmt.Sprintf(`
resource "looker_service_account" "test" {
  service_account_name = "%s"
  is_disabled          = true
}
`, name)
}
