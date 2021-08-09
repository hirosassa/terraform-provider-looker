package looker

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_Connection(t *testing.T) {
	name1 := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: connectionConfig(name1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("looker_connection.test", "name", name1),
					resource.TestCheckResourceAttr("looker_connection.test", "host", "test_project"),
				),
			},
			{
				ResourceName:      "looker_connection.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func connectionConfig(name string) string {
	return fmt.Sprintf(`
	resource "looker_connection" "test" {
		name = %s
		host = "test_project"
		user = var.gcp_service_account_email
		certificate = var.gcp_service_account_json
		file_type = ".json"
		database = "test_dataset"
		tmp_db_name = "tmp_test_dataset"
		dialetct_name = "bigquery_standard_sql"
	}
	`, name)
}
