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

func TestAcc_ModelSet(t *testing.T) {
	name1 := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	name2 := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: modelSetConfig(name1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("looker_model_set.test", "name", name1),
					resource.TestCheckResourceAttr("looker_model_set.test", "models.#", "1"),
				),
			},
			{
				Config: modelSetConfig(name2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("looker_model_set.test", "name", name2),
					resource.TestCheckResourceAttr("looker_model_set.test", "models.#", "1"),
				),
			},
			{
				ResourceName:      "looker_model_set.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
		CheckDestroy: testAccCheckModelSetDestroy,
	})
}

func testAccCheckModelSetDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*apiclient.LookerSDK)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "looker_model_set" {
			continue
		}

		modelSetID := rs.Primary.ID

		modelSet, err := client.ModelSet(modelSetID, "", nil)
		if err != nil {
			if strings.Contains(err.Error(), "404") {
				return nil // successfully destroyed
			}
			return err
		}

		if *modelSet.Name == rs.Primary.Attributes["name"] {
			return fmt.Errorf("model_set '%s' still exists", rs.Primary.ID)
		}

	}

	return nil
}

func modelSetConfig(name string) string {
	return fmt.Sprintf(`
	resource "looker_model_set" "test" {
		name = "%s"
		models = ["test"]
	}
	`, name)
}
