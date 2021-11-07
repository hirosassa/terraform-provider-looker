package looker

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_User(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: providers(),
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
	})
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
