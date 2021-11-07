package looker

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apiclient "github.com/looker-open-source/sdk-codegen/go/sdk/v4"
)

func resourceGroupMembership() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGroupMembershipCreate,
		ReadContext:   resourceGroupMembershipRead,
		UpdateContext: resourceGroupMembershipUpdate,
		DeleteContext: resourceGroupMembershipDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"target_group_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"user_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
				Set:      schema.HashInt,
			},
			"group_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
				Set:      schema.HashInt,
			},
		},
	}
}

func resourceGroupMembershipCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	targetGroupID := int64(d.Get("target_group_id").(int))

	// add users
	userIDs := expandInt64ListFromSet(d.Get("user_ids"))
	err := addGroupUsers(m, targetGroupID, userIDs)
	if err != nil {
		return diag.FromErr(err)
	}

	// add groups
	groupIDs := expandInt64ListFromSet(d.Get("group_ids"))
	err = addGroupGroups(m, targetGroupID, groupIDs)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(int(targetGroupID)))

	return resourceGroupMembershipRead(ctx, d, m)
}

func resourceGroupMembershipRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	targetGroupID := int64(d.Get("target_group_id").(int))

	req := apiclient.RequestAllGroupUsers{
		GroupId: targetGroupID,
	}

	users, err := client.AllGroupUsers(req, nil) // todo: imeplement paging
	if err != nil {
		return diag.FromErr(err)
	}

	groups, err := client.AllGroupGroups(targetGroupID, "", nil) // todo: imeplement paging
	if err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("target_group_id", int(targetGroupID)); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("user_ids", flattenUserIDs(users)); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("group_ids", flattenGroupIDs(groups)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceGroupMembershipUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	targetGroupID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	err = removeAllUsersFromGroup(m, targetGroupID)
	if err != nil {
		return diag.FromErr(err)
	}

	err = removeAllGroupsFromGroup(m, targetGroupID)
	if err != nil {
		return diag.FromErr(err)
	}

	userIDs := expandInt64ListFromSet(d.Get("user_ids"))
	err = addGroupUsers(m, targetGroupID, userIDs)
	if err != nil {
		return diag.FromErr(err)
	}

	groupIDs := expandInt64ListFromSet(d.Get("group_ids"))
	err = addGroupGroups(m, targetGroupID, groupIDs)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceGroupMembershipRead(ctx, d, m)
}

func resourceGroupMembershipDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	targetGroupID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	err = removeAllUsersFromGroup(m, targetGroupID)
	if err != nil {
		return diag.FromErr(err)
	}

	err = removeAllGroupsFromGroup(m, targetGroupID)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceGroupMembershipRead(ctx, d, m)
}

func addGroupUsers(m interface{}, targetGroupID int64, userIDs []int64) error {
	client := m.(*apiclient.LookerSDK)

	for _, userID := range userIDs {
		body := apiclient.GroupIdForGroupUserInclusion{
			UserId: &userID,
		}

		_, err := client.AddGroupUser(targetGroupID, body, nil)
		if err != nil {
			return err
		}
	}

	return nil
}

func addGroupGroups(m interface{}, targetGroupID int64, groupIDs []int64) error {
	client := m.(*apiclient.LookerSDK)

	for _, groupID := range groupIDs {
		body := apiclient.GroupIdForGroupInclusion{
			GroupId: &groupID,
		}

		_, err := client.AddGroupGroup(targetGroupID, body, nil)
		if err != nil {
			return err
		}
	}

	return nil
}

func removeAllUsersFromGroup(m interface{}, groupID int64) error {
	client := m.(*apiclient.LookerSDK)
	req := apiclient.RequestAllGroupUsers{
		GroupId: groupID,
	}

	users, err := client.AllGroupUsers(req, nil) // todo: imeplement paging
	if err != nil {
		return err
	}

	for _, user := range users {
		err = client.DeleteGroupUser(groupID, *user.Id, nil)
		if err != nil {
			return err
		}
	}

	return nil
}

func removeAllGroupsFromGroup(m interface{}, groupID int64) error {
	client := m.(*apiclient.LookerSDK)
	groups, err := client.AllGroupGroups(groupID, "", nil) // todo: imeplement paging
	if err != nil {
		return err
	}

	for _, group := range groups {
		err = client.DeleteGroupFromGroup(groupID, *group.Id, nil)
		if err != nil {
			return err
		}
	}

	return nil
}

func flattenUserIDs(users []apiclient.User) []int {
	userIDs := make([]int, 0, len(users))
	for _, user := range users {
		userIDs = append(userIDs, int(*user.Id))
	}
	return userIDs
}

func flattenGroupIDs(groups []apiclient.Group) []int {
	groupIDs := make([]int, 0, len(groups))
	for _, group := range groups {
		groupIDs = append(groupIDs, int(*group.Id))
	}
	return groupIDs
}
