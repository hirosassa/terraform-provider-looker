package looker

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceGroup(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	dataSourceName := "data.looker_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// Look up with a lower-cased name to verify case-insensitive matching.
				Config: groupConfig(name) + `
data "looker_group" "test" {
  name = lower(looker_group.test.name)
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "id", "looker_group.test", "id"),
					resource.TestCheckResourceAttr(dataSourceName, "name", name),
				),
			},
		},
		CheckDestroy: testAccCheckGroupDestroy,
	})
}

func TestAccDataSourceGroup_notFound(t *testing.T) {
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
data "looker_group" "test" {
  name = "%s"
}
`, name),
				ExpectError: regexp.MustCompile("no group found with name"),
			},
		},
	})
}
