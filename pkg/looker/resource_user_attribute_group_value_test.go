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

func TestAcc_UserAttributeGroupValue(t *testing.T) {
	groupValue1 := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	groupValue2 := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: userAttributeGroupValueConfig(groupValue1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserAttributeGroupValueExists("looker_user_attribute_group_value.test"),
					resource.TestCheckResourceAttr("looker_user_attribute_group_value.test", "value", groupValue1),
				),
			},
			// Test: Update
			{
				Config: userAttributeGroupValueConfig(groupValue2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserAttributeGroupValueExists("looker_user_attribute_group_value.test"),
					resource.TestCheckResourceAttr("looker_user_attribute_group_value.test", "value", groupValue2),
				),
			},
			// Test: Import
			{
				ResourceName:      "looker_user_attribute_group_value.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
		CheckDestroy: testAccCheckUserAttributeGroupValueDestroy,
	})
}

func testAccCheckUserAttributeGroupValueExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("user attribute group value setting not found from resources: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no user attribute group value setting ID is set")
		}

		_, userAttributeIDString, err := parseTwoPartID(rs.Primary.ID)
		if err != nil {
			return err
		}

		userAttributeID := userAttributeIDString

		client := testAccProvider.Meta().(*apiclient.LookerSDK)
		userAttributeGroupValues, err := client.AllUserAttributeGroupValues(userAttributeID, "", nil)
		if err != nil {
			return err
		}

		if len(userAttributeGroupValues) != 1 {
			return fmt.Errorf("looker_user_attribute_group_value '%s' doesn't exist", rs.Primary.ID)
		}

		return nil
	}
}

func testAccCheckUserAttributeGroupValueDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*apiclient.LookerSDK)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "looker_user_attribute_group_value" {
			continue
		}

		_, userAttributeID, err := parseTwoPartID(rs.Primary.ID)
		if err != nil {
			return err
		}

		userAttributeGroupValues, err := client.AllUserAttributeGroupValues(userAttributeID, "", nil)
		if err != nil {
			if strings.Contains(err.Error(), "404") {
				return nil // successfully destroyed
			}
			return err
		}

		if len(userAttributeGroupValues) != 0 {
			return fmt.Errorf("looker_user_attribute_group_value '%s' still exists", rs.Primary.ID)
		}
	}

	return nil
}

func userAttributeGroupValueConfig(groupValue string) string {
	return fmt.Sprintf(`
	resource "looker_group" "test" {
        name = "testing"
	}
	resource "looker_user_attribute" "test" {
        name  = "testing"
        type  = "string"
        label = "testing"
	}
	resource "looker_user_attribute_group_value" "test" {
		group_id          = looker_group.test.id
		user_attribute_id = looker_user_attribute.test.id
        value             = "%s"
	}
	`, groupValue)
}
