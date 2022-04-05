package looker

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apiclient "github.com/looker-open-source/sdk-codegen/go/sdk/v4"
)

func resourcePermissionSet() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePermissionSetCreate,
		ReadContext:   resourcePermissionSetRead,
		UpdateContext: resourcePermissionSetUpdate,
		DeleteContext: resourcePermissionSetDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"permissions": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourcePermissionSetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	permissionSetName := d.Get("name").(string)

	var permissions []string
	for _, permission := range d.Get("permissions").(*schema.Set).List() {
		permissions = append(permissions, permission.(string))
	}

	writePermissionSet := apiclient.WritePermissionSet{
		Name:        &permissionSetName,
		Permissions: &permissions,
	}

	permissionSet, err := client.CreatePermissionSet(writePermissionSet, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	permissionSetID := *permissionSet.Id
	d.SetId(permissionSetID)

	return resourcePermissionSetRead(ctx, d, m)
}

func resourcePermissionSetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	permissionSetID := d.Id()

	permissionSet, err := client.PermissionSet(permissionSetID, "", nil)
	if err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("name", permissionSet.Name); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("permissions", permissionSet.Permissions); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourcePermissionSetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	permissionSetID := d.Id()

	permissionSetName := d.Get("name").(string)
	var permissions []string
	for _, permission := range d.Get("permissions").(*schema.Set).List() {
		permissions = append(permissions, permission.(string))
	}
	writePermissionSet := apiclient.WritePermissionSet{
		Name:        &permissionSetName,
		Permissions: &permissions,
	}
	_, err := client.UpdatePermissionSet(permissionSetID, writePermissionSet, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourcePermissionSetRead(ctx, d, m)
}

func resourcePermissionSetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	permissionSetID := d.Id()

	_, err := client.DeletePermissionSet(permissionSetID, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
