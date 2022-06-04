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

func TestAcc_UserAttribute(t *testing.T) {
	name1 := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	name2 := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: userAttributeConfig(name1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserAttributeExists("looker_user_attribute.test"),
					resource.TestCheckResourceAttr("looker_user_attribute.test", "name", name1),
				),
			},
			// Test: Update
			{
				Config: userAttributeGroupValueConfig(name2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserAttributeGroupValueExists("looker_user_attribute.test"),
					resource.TestCheckResourceAttr("looker_user_attribute.test", "name", name2),
				),
			},
			// Test: Import
			{
				ResourceName:      "looker_user_attribute.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
		CheckDestroy: testAccCheckUserAttributeDestroy,
	})
}

func TestAcc_UserAttributeWithDefaultValue(t *testing.T) {
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	type1 := "advanced_filter_string"
	type2 := "advanced_filter_number"
	defaultValue1 := "%, NULL"
	defaultValue2 := "<0, >=0, NULL"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: userAttributeConfigWithDefaultValue(name, type1, defaultValue1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserAttributeExists("looker_user_attribute.test_with_default"),
					resource.TestCheckResourceAttr("looker_user_attribute.test_with_default", "name", name),
					resource.TestCheckResourceAttr("looker_user_attribute.test_with_default", "type", type1),
					resource.TestCheckResourceAttr("looker_user_attribute.test_with_default", "default_value", defaultValue1),
				),
			},
			// Test: Update
			{
				Config: userAttributeConfigWithDefaultValue(name, type2, defaultValue2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserAttributeExists("looker_user_attribute.test_with_default"),
					resource.TestCheckResourceAttr("looker_user_attribute.test_with_default", "name", name),
					resource.TestCheckResourceAttr("looker_user_attribute.test_with_default", "type", type2),
					resource.TestCheckResourceAttr("looker_user_attribute.test_with_default", "default_value", defaultValue2),
				),
			},
			// Test: Import
			{
				ResourceName:      "looker_user_attribute.test_with_default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
		CheckDestroy: testAccCheckUserAttributeDestroy,
	})
}

func testAccCheckUserAttributeExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("user attribute setting not found from resources: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no user attribute setting ID is set")
		}

		client := testAccProvider.Meta().(*apiclient.LookerSDK)
		userAttribute, err := client.UserAttribute(rs.Primary.ID, "", nil)
		if err != nil {
			return err
		}

		if userAttribute.Name != rs.Primary.Attributes["name"] {
			return fmt.Errorf("looker_user_attribute '%s' does not exist", rs.Primary.ID)
		}

		return nil
	}
}

func testAccCheckUserAttributeDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*apiclient.LookerSDK)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "looker_user_attribute" {
			continue
		}

		userAttribute, err := client.UserAttribute(rs.Primary.ID, "", nil)
		if err != nil {
			if strings.Contains(err.Error(), "404") {
				return nil // successfully destroyed
			}
			return err
		}

		if userAttribute.Name == rs.Primary.Attributes["name"] {
			return fmt.Errorf("looker_user_attribute '%s' still exists", rs.Primary.ID)
		}
	}

	return nil
}

func userAttributeConfig(name string) string {
	return fmt.Sprintf(`
	resource "looker_user_attribute" "test" {
        name  = "%s"
        type  = "string"
        label = "testing"
	}
	`, name)
}

func userAttributeConfigWithDefaultValue(name, dataType, defaultValue string) string {
	return fmt.Sprintf(`
	resource "looker_user_attribute" "test_with_default" {
        name  = "%s"
        type  = "%s"
        label = "testing"
        default_value = "%s"
	}
	`, name, dataType, defaultValue)
}
