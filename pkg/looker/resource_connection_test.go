package looker

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_Connection(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: connectionConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("looker_connection.test", "name", strings.ToLower(name)),
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
		name         = "%s"
		host         = "test_project"
		username     = "test@testproject.iam.gserviceaccount.com"
		certificate  = filebase64("testdata/gcp-sa.json")
		file_type    = ".json"
		database     = "test_dataset"
		tmp_db_name  = "tmp_test_dataset"
		dialect_name = "bigquery_standard_sql"
	}
	`, name)
}
