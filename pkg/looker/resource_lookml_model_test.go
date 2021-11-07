package looker

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_LookMLModel(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	connectionName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	projectName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: lookMLModelConfig(name, connectionName, projectName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("looker_lookml_model.test", "name", name),
					resource.TestCheckResourceAttr("looker_lookml_model.test", "project_name", projectName),
				),
			},
			{
				ResourceName:      "looker_lookml_model.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func lookMLModelConfig(name, connectionName, projectName string) string {
	return fmt.Sprintf(`
	resource "looker_connection" "test" {
		name         = "%s"
		host         = "testproject"
		username     = "test@testproject.iam.gserviceaccount.com"
		certificate  = filebase64("testdata/gcp-sa.json")
		file_type    = ".json"
		database     = "test_dataset"
		tmp_db_name  = "tmp_test_dataset"
		dialect_name = "bigquery_standard_sql"
	}
	resource "looker_lookml_model" "test" {
		name = "%s"
		allowed_db_connection_names = [looker_connection.test.name]
		project_name = "%s"
	}
	`, connectionName, name, projectName)
}
