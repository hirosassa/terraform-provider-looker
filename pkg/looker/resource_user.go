package looker

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
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
		},
	}
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)
	firstName := d.Get("first_name").(string)
	lastName := d.Get("last_name").(string)
	email := d.Get("email").(string)

	writeUser := apiclient.WriteUser{
		FirstName: &firstName,
		LastName:  &lastName,
	}

	user, err := client.CreateUser(writeUser, "", nil)
	if err != nil {
		return diag.FromErr(err)
	}

	userID := *user.Id
	d.SetId(strconv.Itoa(int(userID)))

	writeCredentialsEmail := apiclient.WriteCredentialsEmail{
		Email: &email,
	}
	_, err = client.CreateUserCredentialsEmail(userID, writeCredentialsEmail, "", nil)
	if err != nil {
		if _, err = client.DeleteUser(userID, nil); err != nil {
			return diag.FromErr(err)
		}
		return diag.FromErr(err)
	}

	return resourceUserRead(ctx, d, m)
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	userID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	user, err := client.User(userID, "", nil)
	if err != nil {
		return diag.FromErr(err)
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

	return nil
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	userID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChanges("first_name", "last_name") {
		firstName := d.Get("first_name").(string)
		lastName := d.Get("last_name").(string)
		writeUser := apiclient.WriteUser{
			FirstName: &firstName,
			LastName:  &lastName,
		}
		_, err = client.UpdateUser(userID, writeUser, "", nil)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("email") {
		email := d.Get("email").(string)
		writeCredentialsEmail := apiclient.WriteCredentialsEmail{
			Email: &email,
		}
		_, err = client.UpdateUserCredentialsEmail(userID, writeCredentialsEmail, "", nil)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceUserRead(ctx, d, m)
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	userID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.DeleteUser(userID, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
