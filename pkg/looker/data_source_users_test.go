package looker

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceUsers(t *testing.T) {
	dataSourceName := "data.looker_users.test"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceUsersConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "users.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "users.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "users.0.email"),
					resource.TestCheckResourceAttrSet(dataSourceName, "users.0.first_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "users.0.last_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "users.0.is_disabled"),
				),
			},
		},
	})
}

func testAccDataSourceUsersConfig() string {
	return `
data "looker_users" "test" {
}
`
}
