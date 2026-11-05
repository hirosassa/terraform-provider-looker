package looker

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceUser(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	dataSourceName := "data.looker_user.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceUserConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "id", "looker_user.test", "id"),
					resource.TestCheckResourceAttr(dataSourceName, "first_name", name),
					resource.TestCheckResourceAttr(dataSourceName, "last_name", name),
					resource.TestCheckResourceAttr(dataSourceName, "is_disabled", "false"),
				),
			},
		},
		CheckDestroy: testAccCheckUserDestroy,
	})
}

func TestAccDataSourceUser_notFound(t *testing.T) {
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
data "looker_user" "test" {
  email = "%s@example.com"
}
`, name),
				ExpectError: regexp.MustCompile("no user found with email"),
			},
		},
	})
}

func testAccDataSourceUserConfig(name string) string {
	// Look up with a lower-cased email to verify case-insensitive matching.
	return fmt.Sprintf(`
resource "looker_user" "test" {
  first_name = "%[1]s"
  last_name  = "%[1]s"
  email      = "%[1]s@example.com"
}

data "looker_user" "test" {
  email = lower(looker_user.test.email)
}
`, name)
}
