package looker

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apiclient "github.com/looker-open-source/sdk-codegen/go/sdk/v4"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		UpdateContext: resourceUserUpdate,
		DeleteContext: resourceUserDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"email": {
				Type:     schema.TypeString,
				Required: true,
			},
			"first_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"last_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_disabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"send_setup_link_on_create": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"can_manage_api3_creds": {
				Type:     schema.TypeBool,
				Optional: true,
				// Computed so that users whose flag was turned on outside Terraform don't get a diff back to false when this is omitted.
				Computed: true,
				Description: "Whether the user can create, view, and delete API keys for their own account. " +
					"When omitted, the current value on Looker is left unchanged. " +
					"This field is marked as experimental in the Looker API and may not be available on your instance.",
			},
		},
	}
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)
	firstName := d.Get("first_name").(string)
	lastName := d.Get("last_name").(string)
	email := d.Get("email").(string)
	isDisabled := d.Get("is_disabled").(bool)

	writeUser := apiclient.WriteUser{
		FirstName:  &firstName,
		LastName:   &lastName,
		IsDisabled: &isDisabled,
	}
	// Why not d.Get: it can't distinguish "unset" from false, and we avoid sending this experimental field unless it is explicitly configured.
	if raw := d.GetRawConfig().GetAttr("can_manage_api3_creds"); !raw.IsNull() {
		canManageAPI3Creds := raw.True()
		writeUser.CanManageApi3Creds = &canManageAPI3Creds
	}

	// CreateUser sometimes returns 500 error
	var user apiclient.User
	err := resource.RetryContext(ctx, 1*time.Minute, func() *resource.RetryError {
		var err error

		user, err = client.CreateUser(writeUser, "", nil)
		if err != nil {
			if d.IsNewResource() && strings.Contains(err.Error(), "500") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return diag.FromErr(wrapSDKError(err, "CreateUser", "user", "%s", email))
	}

	userID := *user.Id
	d.SetId(userID)

	writeCredentialsEmail := apiclient.WriteCredentialsEmail{
		Email: &email,
	}

	_, err = client.CreateUserCredentialsEmail(userID, writeCredentialsEmail, "", nil)
	if err != nil {
		if _, err = client.DeleteUser(userID, nil); err != nil {
			return diag.FromErr(wrapSDKError(err, "DeleteUser", "user", "email=%s, id=%s", email, userID))
		}
		return diag.FromErr(wrapSDKError(err, "CreateUserCredentialsEmail", "user", "%s", email))
	}

	// Send setup mail if requested
	sendSetupMail := d.Get("send_setup_link_on_create").(bool)
	if sendSetupMail {
		_, err = client.SendUserCredentialsEmailPasswordReset(userID, "", nil)
		if err != nil {
			// Log the error but don't fail the resource creation
			// since the user was successfully created
			return diag.Errorf("User created successfully but failed to send setup email: %v",
				wrapSDKError(err, "SendUserCredentialsEmailPasswordReset", "user", "%s", email))
		}
	}

	return resourceUserRead(ctx, d, m)
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	userID := d.Id()

	user, err := client.User(userID, "", nil)
	if err != nil {
		return diag.FromErr(wrapSDKError(err, "User", "user", "%s", userID))
	}

	if err = d.Set("email", user.Email); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("first_name", user.FirstName); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("last_name", user.LastName); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("is_disabled", user.IsDisabled); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("can_manage_api3_creds", user.CanManageApi3Creds); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	userID := d.Id()

	if d.HasChanges("first_name", "last_name", "is_disabled", "can_manage_api3_creds") {
		firstName := d.Get("first_name").(string)
		lastName := d.Get("last_name").(string)
		isDisabled := d.Get("is_disabled").(bool)
		email := d.Get("email").(string)
		writeUser := apiclient.WriteUser{
			FirstName:  &firstName,
			LastName:   &lastName,
			IsDisabled: &isDisabled,
		}
		if d.HasChange("can_manage_api3_creds") {
			canManageAPI3Creds := d.Get("can_manage_api3_creds").(bool)
			writeUser.CanManageApi3Creds = &canManageAPI3Creds
		}
		_, err := client.UpdateUser(userID, writeUser, "", nil)
		if err != nil {
			return diag.FromErr(wrapSDKError(err, "UpdateUser", "user", "email=%s, id=%s", email, userID))
		}
	}

	if d.HasChange("email") {
		email := d.Get("email").(string)
		writeCredentialsEmail := apiclient.WriteCredentialsEmail{
			Email: &email,
		}
		_, err := client.UpdateUserCredentialsEmail(userID, writeCredentialsEmail, "", nil)
		if err != nil {
			return diag.FromErr(wrapSDKError(err, "UpdateUserCredentialsEmail", "user", "email=%s, id=%s", email, userID))
		}
	}

	return resourceUserRead(ctx, d, m)
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	userID := d.Id()
	email := d.Get("email").(string)

	_, err := client.DeleteUser(userID, nil)
	if err != nil {
		return diag.FromErr(wrapSDKError(err, "DeleteUser", "user", "email=%s, id=%s", email, userID))
	}

	return nil
}
