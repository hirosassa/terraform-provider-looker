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

func TestAcc_GroupMembership(t *testing.T) {
	target1 := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	target2 := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	group1 := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	group2 := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	user1 := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	user2 := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	user3 := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	user4 := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	user5 := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: groupMembershipConfig(target1, user1, user2, group1, group2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGroupMembershipExists("looker_group_membership.test"),
					resource.TestCheckResourceAttr("looker_group_membership.test", "group_ids.#", "2"),
					resource.TestCheckResourceAttr("looker_group_membership.test", "user_ids.#", "2"),
				),
			},
			{
				Config: groupMembershipConfigNoGroup(target2, user3, user4),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGroupMembershipExists("looker_group_membership.test_no_group"),
					resource.TestCheckResourceAttr("looker_group_membership.test_no_group", "user_ids.#", "2"),
				),
			},
			// Test: Update
			{
				Config: groupMembershipConfigUpdate(target1, user1, user2, user5, group1),
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
		targetGroupID := rs.Primary.ID

		users, _ := client.AllGroupUsers(apiclient.RequestAllGroupUsers{GroupId: targetGroupID}, nil)

		groups, _ := client.AllGroupGroups(targetGroupID, "", nil)

		if users != nil && len(users) == 0 && groups != nil && len(groups) == 0 {
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

		targetGroupID := rs.Primary.ID

		users, err := client.AllGroupUsers(apiclient.RequestAllGroupUsers{GroupId: targetGroupID}, nil)
		if err != nil {
			if strings.Contains(err.Error(), "404") {
				return nil // successfully destroyed
			}
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

func groupMembershipConfig(target, user1, user2, group1, group2 string) string {
	return fmt.Sprintf(`
	resource "looker_group" "target_group" {
		name = "%s"
	}
	resource "looker_user" "membership_user1" {
        first_name = "%s"
        last_name  = "%s"
        email      = "%s@example.com"
	}
	resource "looker_user" "membership_user2" {
        first_name = "%s"
        last_name  = "%s"
        email      = "%s@example.com"
	}
	resource "looker_group" "membership_group1" {
		name = "%s"
	}
	resource "looker_group" "membership_group2" {
		name = "%s"
	}
	resource "looker_group_membership" "test" {
		target_group_id = looker_group.target_group.id
		user_ids        = [looker_user.membership_user1.id, looker_user.membership_user2.id]
		group_ids       = [looker_group.membership_group1.id, looker_group.membership_group2.id]
	}
	`, target, user1, user1, user1, user2, user2, user2, group1, group2)
}

func groupMembershipConfigNoGroup(target, user1, user2 string) string {
	return fmt.Sprintf(`
	resource "looker_group" "target_no_group" {
		name = "%s"
	}
	resource "looker_user" "membership_user3" {
        first_name = "%s"
        last_name  = "%s"
        email      = "%s@example.com"
	}
	resource "looker_user" "membership_user4" {
        first_name = "%s"
        last_name  = "%s"
        email      = "%s@example.com"
	}
	resource "looker_group_membership" "test_no_group" {
		target_group_id = looker_group.target_no_group.id
		user_ids        = [looker_user.membership_user3.id, looker_user.membership_user4.id]
	}
	`, target, user1, user1, user1, user2, user2, user2)
}

func groupMembershipConfigUpdate(target, user1, user2, user5, group1 string) string {
	return fmt.Sprintf(`
	resource "looker_group" "target_group" {
		name = "%s"
	}
	resource "looker_user" "membership_user1" {
        first_name = "%s"
        last_name  = "%s"
        email      = "%s@example.com"
	}
	resource "looker_user" "membership_user2" {
        first_name = "%s"
        last_name  = "%s"
        email      = "%s@example.com"
	}
	resource "looker_user" "membership_user5" {
        first_name = "%s"
        last_name  = "%s"
        email      = "%s@example.com"
	}
	resource "looker_group" "membership_group1" {
		name = "%s"
	}
	resource "looker_group_membership" "test" {
		target_group_id = looker_group.target_group.id
		user_ids        = [looker_user.membership_user1.id, looker_user.membership_user2.id, looker_user.membership_user5.id]
		group_ids       = [looker_group.membership_group1.id]
	}
	`, target, user1, user1, user1, user2, user2, user2, user5, user5, user5, group1)
}
