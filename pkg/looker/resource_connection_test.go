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

func TestAcc_Connection(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: connectionConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConnectionExists("looker_connection.test"),
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
		CheckDestroy: testAccCheckConnectionDestroy,
	})
}

func testAccCheckConnectionExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("connection setting not found from resources: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no connection setting ID is set")
		}

		client := testAccProvider.Meta().(*apiclient.LookerSDK)
		connectionName := rs.Primary.ID

		_, err := client.Connection(connectionName, "", nil)
		if err !=nil {
			return err
		}

		return nil
	}
}
func testAccCheckConnectionDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*apiclient.LookerSDK)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "looker_connection" {
			continue
		}

		connectionName := rs.Primary.ID
		_, err := client.Connection(connectionName, "", nil)
		if err != nil {
			if strings.Contains(err.Error(), "404") {
				return nil // successfully destroyed
			}
			return err
		}
	}

	return nil
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
