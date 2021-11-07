package looker

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	apiclient "github.com/looker-open-source/sdk-codegen/go/sdk/v4"
)

func TestAcc_GroupMembership(t *testing.T) {
	target := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	group1 := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	group2 := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	user1 := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	user2 := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	user3 := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: groupMembershipConfig(target, group1, group2, user1, user2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGroupMembershipExists("looker_group_membership.test"),
					resource.TestCheckResourceAttr("looker_group_membership.test", "group_ids.#", "2"),
					resource.TestCheckResourceAttr("looker_group_membership.test", "user_ids.#", "2"),
				),
			},
			// Test: Create
			{
				Config: groupMembershipConfigUpdate(user3),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGroupMembershipExists("looker_group_membership.test"),
					resource.TestCheckResourceAttr("looker_group_membership.test", "group_ids.#", "1"),
					resource.TestCheckResourceAttr("looker_group_membership.test", "user_ids.#", "3"),
				),
			},
			// Test: Import
			{
				ResourceName:      "looker_group_membership.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
		CheckDestroy: testAccCheckGroupMembershipDestroy,
	})
}

func testAccCheckGroupMembershipExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("group membership setting not found from resources: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no group membership setting ID is set")
		}

		client := testAccProvider.Meta().(*apiclient.LookerSDK)
		targetGroupID, err := strconv.ParseInt(rs.Primary.ID, 10, 64)
		if err != nil {
			return err
		}

		users, err := client.AllGroupUsers(apiclient.RequestAllGroupUsers{GroupId: targetGroupID}, nil)
		if err != nil {
			return err
		}

		groups, err := client.AllGroupGroups(targetGroupID, "", nil)
		if err != nil {
			return err
		}

		if len(users) == 0 && len(groups) == 0 {
			return fmt.Errorf("no group members are set: %s", n)
		}

		return nil
	}
}

func testAccCheckGroupMembershipDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*apiclient.LookerSDK)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "looker_membership" {
			continue
		}

		targetGroupID, err := strconv.ParseInt(rs.Primary.ID, 10, 64)
		if err != nil {
			return err
		}

		users, err := client.AllGroupUsers(apiclient.RequestAllGroupUsers{GroupId: targetGroupID}, nil)
		if err != nil {
			return err
		}

		if len(users) != 0 {
			return fmt.Errorf("group_membership '%s' still exists", rs.Primary.ID)
		}

		groups, err := client.AllGroupGroups(targetGroupID, "", nil)
		if err != nil {
			return err
		}

		if len(groups) != 0 {
			return fmt.Errorf("group_membership '%s' still exists", rs.Primary.ID)
		}

	}

	return nil
}

func groupMembershipConfig(target, group1, group2, user1, user2 string) string {
	return fmt.Sprintf(`
	resource "looker_group" "target_group" {
		name = "%s"
	}
	resource "looker_user" "user1" {
        email = "%s@example.com"
	}
	resource "looker_user" "user2" {
        email = "%s@example.com"
	}
	resource "looker_group" "group1" {
		name = "%s"
	}
	resource "looker_group" "group2" {
		name = "%s"
	}
	resource "looker_group_membership" "test" {
		target_group_id = looker_group.target_group.id
		user_ids        = [looker_user.user1.id, looker_user.user2.id]
		group_ids       = [looker_group.group1.id, looker_group.group2.id]
	}
	`, target, user1, user2, group1, group2)
}

func groupMembershipConfigUpdate(user string) string {
	return fmt.Sprintf(`
	resource "looker_user" "user3" {
        email = "%s@example.com"
	}
	resource "looker_group_membership" "test" {
		target_group_id = looker_group.target_group.id
		user_ids        = [looker_user.user1.id, looker_user.user2.id, looker_user.user3.id]
		group_ids       = [looker_group.group1.id]
	}
	`, user)
}
