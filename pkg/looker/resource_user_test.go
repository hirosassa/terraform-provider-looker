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

func TestAcc_User(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: userConfig(name, name, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("looker_user.test", "first_name", name),
					resource.TestCheckResourceAttr("looker_user.test", "last_name", name),
					resource.TestCheckResourceAttr("looker_user.test", "email", fmt.Sprintf("%s@example.com", name)),
				),
			},
			{
				ResourceName:      "looker_user.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
		CheckDestroy: testAccCheckUserDestroy,
	})
}

func testAccCheckUserDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*apiclient.LookerSDK)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "looker_user" {
			continue
		}

		userID := rs.Primary.ID

		user, err := client.User(userID, "", nil)
		if err != nil {
			if strings.Contains(err.Error(), "404") {
				return nil // successfully destroyed
			}
			return err
		}

		if *user.Email == rs.Primary.Attributes["email"] {
			return fmt.Errorf("user '%s' still exists", rs.Primary.ID)
		}

	}

	return nil
}

func userConfig(firstName, lastName, email string) string {
	return fmt.Sprintf(`
	resource "looker_user" "test" {
		first_name = "%s"
		last_name = "%s"
		email = "%s@example.com"
	}
	`, firstName, lastName, email)
}
