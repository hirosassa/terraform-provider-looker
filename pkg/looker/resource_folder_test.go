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

func TestAcc_Folder(t *testing.T) {
	name1 := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	name2 := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	parentID := "1" // Default Shared folder ID

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: folderConfig(name1, parentID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("looker_folder.test", "name", name1),
					resource.TestCheckResourceAttr("looker_folder.test", "parent_id", parentID),
				),
			},
			{
				Config: folderConfig(name2, parentID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("looker_folder.test", "name", name2),
					resource.TestCheckResourceAttr("looker_folder.test", "parent_id", parentID),
				),
			},
			{
				ResourceName:      "looker_folder.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
		CheckDestroy: testAccCheckFolderDestroy,
	})
}

func testAccCheckFolderDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*apiclient.LookerSDK)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "looker_folder" {
			continue
		}

		folderID := rs.Primary.ID

		_, err := client.Folder(folderID, "", nil)
		if err != nil {
			if strings.Contains(err.Error(), "404") {
				continue // successfully destroyed
			}
			return err
		}

		return fmt.Errorf("folder still exists: %s", rs.Primary.ID)
	}

	return nil
}

func folderConfig(name, parentID string) string {
	return fmt.Sprintf(`
	resource "looker_folder" "test" {
		name      = "%s"
		parent_id = "%s"
	}
	`, name, parentID)
}

func folderConfigWithoutParentID(name string) string {
	return fmt.Sprintf(`
	resource "looker_folder" "test_no_parent" {
		name = "%s"
	}
	`, name)
}

func TestAcc_FolderWithoutParentID(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: folderConfigWithoutParentID(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("looker_folder.test_no_parent", "name", name),
				),
			},
		},
		CheckDestroy: testAccCheckFolderDestroy,
	})
}
