package looker

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	apiclient "github.com/looker-open-source/sdk-codegen/go/sdk/v4"
)

func TestAcc_ContentMetadataAccess_Group(t *testing.T) {
	folderName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	groupName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: contentMetadataAccessConfigGroup(folderName, groupName, "view"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("looker_content_metadata_access.test", "permission_type", "view"),
					resource.TestCheckResourceAttrSet("looker_content_metadata_access.test", "group_id"),
				),
			},
			{
				Config: contentMetadataAccessConfigGroup(folderName, groupName, "edit"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("looker_content_metadata_access.test", "permission_type", "edit"),
				),
			},
			// Test: Import
			{
				ResourceName:      "looker_content_metadata_access.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: contentMetadataAccessImportStateIdFunc("looker_content_metadata_access.test"),
			},
		},
		CheckDestroy: testAccCheckContentMetadataAccessDestroy,
	})
}

func TestAcc_ContentMetadataAccess_User(t *testing.T) {
	folderName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	userName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: contentMetadataAccessConfigUser(folderName, userName, "view"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("looker_content_metadata_access.test", "permission_type", "view"),
					resource.TestCheckResourceAttrSet("looker_content_metadata_access.test", "user_id"),
				),
			},
			// Test: Import
			{
				ResourceName:      "looker_content_metadata_access.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: contentMetadataAccessImportStateIdFunc("looker_content_metadata_access.test"),
			},
		},
		CheckDestroy: testAccCheckContentMetadataAccessDestroy,
	})
}

func TestAcc_ContentMetadataAccess_InvalidID(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `
resource "looker_content_metadata_access" "test" {
  content_metadata_id = "999999"
  group_id            = "1"
  permission_type     = "view"
}
`,
				ExpectError: regexp.MustCompile("404"),
			},
		},
	})
}

func contentMetadataAccessImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}
		return rs.Primary.Attributes["content_metadata_id"] + "/" + rs.Primary.ID, nil
	}
}

func testAccCheckContentMetadataAccessDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*apiclient.LookerSDK)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "looker_content_metadata_access" {
			continue
		}

		accessID := rs.Primary.ID
		contentMetadataID := rs.Primary.Attributes["content_metadata_id"]

		accesses, err := client.AllContentMetadataAccesses(contentMetadataID, "", nil)
		if err != nil {
			if strings.Contains(err.Error(), "404") {
				continue
			}
			return err
		}

		for _, a := range accesses {
			if a.Id != nil && *a.Id == accessID {
				return fmt.Errorf("content metadata access still exists: %s", accessID)
			}
		}
	}

	return nil
}

func contentMetadataAccessConfigGroup(folderName, groupName, permission string) string {
	return fmt.Sprintf(`
	resource "looker_folder" "test" {
		name      = "%s"
		parent_id = "1"
		inherits  = false
	}

	resource "looker_group" "test" {
		name = "%s"
	}

	resource "looker_content_metadata_access" "test" {
		content_metadata_id = looker_folder.test.content_metadata_id
		group_id            = looker_group.test.id
		permission_type     = "%s"
	}
	`, folderName, groupName, permission)
}

func contentMetadataAccessConfigUser(folderName, userName, permission string) string {
	return fmt.Sprintf(`
	resource "looker_folder" "test" {
		name      = "%s"
		parent_id = "1"
		inherits  = false
	}

	resource "looker_user" "test" {
		first_name = "%s"
		last_name  = "Test"
		email      = "%s@example.com"
	}

	resource "looker_content_metadata_access" "test" {
		content_metadata_id = looker_folder.test.content_metadata_id
		user_id             = looker_user.test.id
		permission_type     = "%s"
	}
	`, folderName, userName, userName, permission)
}
