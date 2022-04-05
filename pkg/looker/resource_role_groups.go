package looker

import (
	"context"

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

	roleID := d.Get("role_id").(string)

	var groupIDs []string
	for _, groupID := range d.Get("group_ids").(*schema.Set).List() {
		gID := groupID.(string)
		groupIDs = append(groupIDs, gID)
	}

	_, err := client.SetRoleGroups(roleID, groupIDs, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(roleID)

	return resourceRoleGroupsRead(ctx, d, m)
}

func resourceRoleGroupsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	roleID := d.Id()

	groups, err := client.RoleGroups(roleID, "", nil)
	if err != nil {
		return diag.FromErr(err)
	}

	var groupIDs []string
	for _, group := range groups {
		gID := *group.Id
		groupIDs = append(groupIDs, gID)
	}

	if err = d.Set("role_id", roleID); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("group_ids", groupIDs); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceRoleGroupsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	roleID := d.Id()

	var groupIDs []string
	for _, groupID := range d.Get("group_ids").(*schema.Set).List() {
		gID := groupID.(string)
		groupIDs = append(groupIDs, gID)
	}

	_, err := client.SetRoleGroups(roleID, groupIDs, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceRoleGroupsRead(ctx, d, m)
}

func resourceRoleGroupsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	roleID := d.Id()

	groupIDs := []string{}
	_, err := client.SetRoleGroups(roleID, groupIDs, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
