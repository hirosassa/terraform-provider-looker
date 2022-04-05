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

func TestAcc_Group(t *testing.T) {
	name1 := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	name2 := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: groupConfig(name1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("looker_group.test", "name", name1),
				),
			},
			{
				Config: groupConfig(name2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("looker_group.test", "name", name2),
				),
			},
			{
				ResourceName:      "looker_group.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
		CheckDestroy: testAccCheckGroupDestroy,
	})
}

func testAccCheckGroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*apiclient.LookerSDK)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "looker_group" {
			continue
		}

		groupID := rs.Primary.ID

		group, err := client.Group(groupID, "", nil)
		if err != nil {
			if strings.Contains(err.Error(), "404") {
				return nil // successfully destroyed
			}
			return err
		}

		if *group.Name == rs.Primary.Attributes["name"] {
			return fmt.Errorf("group still exists: %s", rs.Primary.ID)
		}
	}

	return nil

}

func groupConfig(name string) string {
	return fmt.Sprintf(`
	resource "looker_group" "test" {
		name = "%s"
	}
	`, name)
}
