package looker

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apiclient "github.com/looker-open-source/sdk-codegen/go/sdk/v4"
)

func resourceGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGroupCreate,
		ReadContext:   resourceGroupRead,
		UpdateContext: resourceGroupUpdate,
		DeleteContext: resourceGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)
	groupName := d.Get("name").(string)

	writeGroup := apiclient.WriteGroup{
		Name: &groupName,
	}

	group, err := client.CreateGroup(writeGroup, "", nil)
	if err != nil {
		return diag.FromErr(wrapSDKError(err, "CreateGroup", "group", "%s", groupName))
	}

	groupID := *group.Id
	d.SetId(groupID)

	return resourceGroupRead(ctx, d, m)
}

func resourceGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	groupID := d.Id()

	group, err := client.Group(groupID, "", nil)
	if err != nil {
		return diag.FromErr(wrapSDKError(err, "Group", "group", "%s", groupID))
	}

	if err = d.Set("name", group.Name); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	groupID := d.Id()

	groupName := d.Get("name").(string)
	writeGroup := apiclient.WriteGroup{
		Name: &groupName,
	}
	_, err := client.UpdateGroup(groupID, writeGroup, "", nil)
	if err != nil {
		return diag.FromErr(wrapSDKError(err, "UpdateGroup", "group", "name=%s, id=%s", groupName, groupID))
	}

	return resourceGroupRead(ctx, d, m)
}

func resourceGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	groupID := d.Id()
	groupName := d.Get("name").(string)

	_, err := client.DeleteGroup(groupID, nil)
	if err != nil {
		return diag.FromErr(wrapSDKError(err, "DeleteGroup", "group", "name=%s, id=%s", groupName, groupID))
	}

	return nil
}
