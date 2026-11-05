package looker

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceRole(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	dataSourceName := "data.looker_role.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: roleConfig(name) + `
data "looker_role" "test" {
  name = lower(looker_role.role_test.name)
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "id", "looker_role.role_test", "id"),
					resource.TestCheckResourceAttr(dataSourceName, "name", name),
					resource.TestCheckResourceAttrPair(dataSourceName, "permission_set_id", "looker_permission_set.role_test", "id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "model_set_id", "looker_model_set.role_test", "id"),
				),
			},
		},
		CheckDestroy: testAccCheckRoleDestroy,
	})
}
