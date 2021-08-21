package looker

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_LookMLModel(t *testing.T) {
	name1 := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	name2 := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: lookMLModelConfig(name1, name2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("looker_lookml_model.test", "name", name1),
					resource.TestCheckResourceAttr("looker_lookml_model.test", "project_name", "name2"),
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

func lookMLModelConfig(name1, name2 string) string {
	return fmt.Sprintf(`
	resource "looker_lookml_model" "test" {
		name = %s
		allowed_db_connection_names = ["bigquery-connection"]
		project_name = %s
	}
	`, name1, name2)
}
