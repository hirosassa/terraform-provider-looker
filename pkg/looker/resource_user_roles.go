package looker

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apiclient "github.com/looker-open-source/sdk-codegen/go/sdk/v4"
)

func resourceUserRoles() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserRolesCreate,
		ReadContext:   resourceUserRolesRead,
		UpdateContext: resourceUserRolesUpdate,
		DeleteContext: resourceUserRolesDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"role_ids": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceUserRolesCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	userIDString := d.Get("user_id").(string)

	userID, err := strconv.ParseInt(userIDString, 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	var roleIDs []int64
	for _, roleID := range d.Get("role_ids").(*schema.Set).List() {
		rID, err := strconv.ParseInt(roleID.(string), 10, 64)
		if err != nil {
			return diag.FromErr(err)
		}
		roleIDs = append(roleIDs, rID)
	}

	_, err = client.SetUserRoles(userID, roleIDs, "", nil)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(userIDString)

	return resourceUserRolesRead(ctx, d, m)
}

func resourceUserRolesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	userID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	request := apiclient.RequestUserRoles{UserId: userID}

	userRoles, err := client.UserRoles(request, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	var roleIDs []string
	for _, role := range userRoles {
		rID := strconv.Itoa(int(*role.Id))
		roleIDs = append(roleIDs, rID)
	}

	if err = d.Set("user_id", d.Id()); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("role_ids", roleIDs); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceUserRolesUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	userID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	var roleIDs []int64
	for _, roleID := range d.Get("role_ids").(*schema.Set).List() {
		rID, err := strconv.ParseInt(roleID.(string), 10, 64)
		if err != nil {
			return diag.FromErr(err)
		}
		roleIDs = append(roleIDs, rID)
	}

	_, err = client.SetUserRoles(userID, roleIDs, "", nil)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceUserRolesRead(ctx, d, m)
}

func resourceUserRolesDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	userID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	roleIDs := []int64{}
	_, err = client.SetUserRoles(userID, roleIDs, "", nil)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
