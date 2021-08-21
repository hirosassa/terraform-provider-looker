package looker

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apiclient "github.com/looker-open-source/sdk-codegen/go/sdk/v4"
)

func resourceRoleGroups() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRoleGroupsCreate,
		ReadContext:   resourceRoleGroupsRead,
		UpdateContext: resourceRoleGroupsUpdate,
		DeleteContext: resourceRoleGroupsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"role_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"group_ids": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceRoleGroupsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	roleIDString := d.Get("role_id").(string)

	roleID, err := strconv.ParseInt(roleIDString, 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	var groupIDs []int64
	for _, groupID := range d.Get("group_ids").(*schema.Set).List() {
		gID, err := strconv.ParseInt(groupID.(string), 10, 64)
		if err != nil {
			return diag.FromErr(err)
		}
		groupIDs = append(groupIDs, gID)
	}

	_, err = client.SetRoleGroups(roleID, groupIDs, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(roleIDString)

	return resourceRoleGroupsRead(ctx, d, m)
}

func resourceRoleGroupsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	roleID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	groups, err := client.RoleGroups(roleID, "", nil)
	if err != nil {
		return diag.FromErr(err)
	}

	var groupIDs []string
	for _, group := range groups {
		gID := strconv.Itoa(int(*group.Id))
		groupIDs = append(groupIDs, gID)
	}

	if err = d.Set("role_id", strconv.Itoa(int(roleID))); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("group_ids", groupIDs); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceRoleGroupsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	roleID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	var groupIDs []int64
	for _, groupID := range d.Get("group_ids").(*schema.Set).List() {
		gID, err := strconv.ParseInt(groupID.(string), 10, 64)
		if err != nil {
			return diag.FromErr(err)
		}
		groupIDs = append(groupIDs, gID)
	}

	_, err = client.SetRoleGroups(roleID, groupIDs, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceRoleGroupsRead(ctx, d, m)
}

func resourceRoleGroupsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	roleID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	groupIDs := []int64{}
	_, err = client.SetRoleGroups(roleID, groupIDs, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
