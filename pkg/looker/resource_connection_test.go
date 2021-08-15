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
	locals {
		gcp_service_account_email = "test@testproject.iam.gserviceaccount.com"
		gcp_service_account_json = <<EOT{
  "type": "service_account",
  "project_id": "testproject",
  "private_key_id": "dummydummydummydummydummydummydummy",
  "private_key": "dummydummydummydummydummydummydummydummydummy",
  "client_email": "test@testproject.iam.gserviceaccount.com",
  "client_id": "1234567890123456789012345678901234567890",
  "auth_uri": "https://accounts.google.com/o/oauth2/auth",
  "token_uri": "https://oauth2.googleapis.com/robot/v1/metadata/x509/test@testproject.iam.gserviceaccount.com"
}
EOT
	}
	resource "looker_connection" "test" {
		name = %s
		host = "testproject"
		user = locals.gcp_service_account_email
		certificate = locals.gcp_service_account_json
		file_type = ".json"
		database = "test_dataset"
		tmp_db_name = "tmp_test_dataset"
		dialetct_name = "bigquery_standard_sql"
	}
	`, name)
}
