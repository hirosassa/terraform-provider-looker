package looker

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apiclient "github.com/looker-open-source/sdk-codegen/go/sdk/v4"
)

func resourceRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRoleCreate,
		ReadContext:   resourceRoleRead,
		UpdateContext: resourceRoleUpdate,
		DeleteContext: resourceRoleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"permission_set_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"model_set_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceRoleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	roleName := d.Get("name").(string)
	permissionSetID := d.Get("permission_set_id").(string)
	modelSetID := d.Get("model_set_id").(string)

	writeRole := apiclient.WriteRole{
		Name:            &roleName,
		PermissionSetId: &permissionSetID,
		ModelSetId:      &modelSetID,
	}

	role, err := client.CreateRole(writeRole, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	roleID := *role.Id
	d.SetId(roleID)

	return resourceRoleRead(ctx, d, m)
}

func resourceRoleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	roleID := d.Id()

	role, err := client.Role(roleID, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("name", role.Name); err != nil {
		return diag.FromErr(err)
	}
	pSetID := *role.PermissionSet.Id
	if err = d.Set("permission_set_id", pSetID); err != nil {
		return diag.FromErr(err)
	}
	mSetID := *role.ModelSet.Id
	if err = d.Set("model_set_id", mSetID); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceRoleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	roleID := d.Id()

	roleName := d.Get("name").(string)
	permissionSetID := d.Get("permission_set_id").(string)
	modelSetID := d.Get("model_set_id").(string)
	writeRole := apiclient.WriteRole{
		Name:            &roleName,
		PermissionSetId: &permissionSetID,
		ModelSetId:      &modelSetID,
	}
	_, err := client.UpdateRole(roleID, writeRole, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceRoleRead(ctx, d, m)
}

func resourceRoleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	roleID := d.Id()

	_, err := client.DeleteRole(roleID, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
