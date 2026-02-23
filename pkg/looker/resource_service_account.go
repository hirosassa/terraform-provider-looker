package looker

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apiclient "github.com/looker-open-source/sdk-codegen/go/sdk/v4"
)

// resourceServiceAccount manages Looker API-only service accounts.
// NOTE: This uses an alpha API endpoint and may be subject to breaking changes.
func resourceServiceAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceServiceAccountCreate,
		ReadContext:   resourceServiceAccountRead,
		UpdateContext: resourceServiceAccountUpdate,
		DeleteContext: resourceServiceAccountDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"service_account_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"is_disabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceServiceAccountCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	name := d.Get("service_account_name").(string)
	isDisabled := d.Get("is_disabled").(bool)

	body := apiclient.WriteServiceAccount{
		ServiceAccountName: &name,
		IsDisabled:         &isDisabled,
	}

	sa, err := client.CreateServiceAccount(body, "", nil)
	if err != nil {
		return diag.FromErr(wrapSDKError(err, "CreateServiceAccount", "service_account", "%s", name))
	}

	d.SetId(*sa.Id)

	return resourceServiceAccountRead(ctx, d, m)
}

func resourceServiceAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	userID := d.Id()

	user, err := client.User(userID, "", nil)
	if err != nil {
		return diag.FromErr(wrapSDKError(err, "User", "service_account", "%s", userID))
	}

	if err = d.Set("service_account_name", user.ServiceAccountName); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("is_disabled", user.IsDisabled); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceServiceAccountUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	userID := d.Id()
	name := d.Get("service_account_name").(string)
	isDisabled := d.Get("is_disabled").(bool)

	body := apiclient.WriteServiceAccount{
		ServiceAccountName: &name,
		IsDisabled:         &isDisabled,
	}

	_, err := client.UpdateServiceAccount(userID, body, "", nil)
	if err != nil {
		return diag.FromErr(wrapSDKError(err, "UpdateServiceAccount", "service_account", "name=%s, id=%s", name, userID))
	}

	return resourceServiceAccountRead(ctx, d, m)
}

func resourceServiceAccountDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	userID := d.Id()
	name := d.Get("service_account_name").(string)

	_, err := client.DeleteServiceAccount(userID, nil)
	if err != nil {
		return diag.FromErr(wrapSDKError(err, "DeleteServiceAccount", "service_account", "name=%s, id=%s", name, userID))
	}

	return nil
}
