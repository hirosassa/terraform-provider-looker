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
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceGroupMembershipCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	groupIDString := d.Get("group_id").(string)

	groupID, err := strconv.ParseInt(groupIDString, 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	userIDString := d.Get("user_id").(string)

	userID, err := strconv.ParseInt(userIDString, 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	body := apiclient.GroupIdForGroupUserInclusion{
		UserId: &userID,
	}

	_, err = client.AddGroupUser(groupID, body, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(buildTwoPartID(&groupIDString, &userIDString))

	return resourceGroupMembershipRead(ctx, d, m)
}

func resourceGroupMembershipRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	id := d.Id()
	groupID, userID, err := groupIDAndUserIDFromID(id)
	if err != nil {
		return diag.FromErr(err)
	}

	req := apiclient.RequestAllGroupUsers{
		GroupId: groupID,
	}

	users, err := client.AllGroupUsers(req, nil) // todo: imeplement paging
	if err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("group_id", strconv.Itoa(int(groupID))); err != nil {
		return diag.FromErr(err)
	}

	if isContained(users, userID) {
		if err = d.Set("user_id", strconv.Itoa(int(userID))); err != nil {
			return diag.FromErr(err)
		}
	} else {
		if err = d.Set("user_id", ""); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func resourceGroupMembershipUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// there's no update in this resource
	return resourceGroupMembershipCreate(ctx, d, m)
}

func resourceGroupMembershipDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	id := d.Id()
	groupID, userID, err := groupIDAndUserIDFromID(id)
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.DeleteGroupUser(groupID, userID, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceGroupMembershipRead(ctx, d, m)
}

func isContained(users []apiclient.User, userID int64) bool {
	for _, user := range users {
		if user.Id == &userID {
			return true
		}
	}
	return false
}

func groupIDAndUserIDFromID(id string) (int64, int64, error) {
	groupIDString, userIDString, err := parseTwoPartID(id)
	if err != nil {
		return 0, 0, err
	}

	groupID, err := strconv.ParseInt(groupIDString, 10, 64)
	if err != nil {
		return 0, 0, err
	}

	userID, err := strconv.ParseInt(userIDString, 10, 64)
	if err != nil {
		return 0, 0, err
	}

	return groupID, userID, err
}
