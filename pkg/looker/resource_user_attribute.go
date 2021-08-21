package looker

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apiclient "github.com/looker-open-source/sdk-codegen/go/sdk/v4"
)

func resourceUserAttribute() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserAttributeCreate,
		ReadContext:   resourceUserAttributeRead,
		UpdateContext: resourceUserAttributeUpdate,
		DeleteContext: resourceUserAttributeDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"label": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceUserAttributeCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)
	userAttributeName := d.Get("name").(string)
	userAttributeLabel := d.Get("label").(string)
	userAttributeType := d.Get("type").(string)

	writeUserAttribute := apiclient.WriteUserAttribute{
		Name:  &userAttributeName,
		Label: &userAttributeLabel,
		Type:  &userAttributeType,
	}

	userAttribute, err := client.CreateUserAttribute(writeUserAttribute, "", nil)
	if err != nil {
		return diag.FromErr(err)
	}

	userAttributeID := *userAttribute.Id
	d.SetId(strconv.Itoa(int(userAttributeID)))

	return resourceUserAttributeRead(ctx, d, m)
}

func resourceUserAttributeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	userAttributeID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	userAttribute, err := client.UserAttribute(userAttributeID, "", nil)
	if err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("name", userAttribute.Name); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("type", userAttribute.Type); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("label", userAttribute.Label); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceUserAttributeUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	userAttributeID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	userAttributeName := d.Get("name").(string)
	userAttributeType := d.Get("type").(string)
	userAttributeLabel := d.Get("type").(string)

	writeUserAttribute := apiclient.WriteUserAttribute{
		Name:  &userAttributeName,
		Label: &userAttributeLabel,
		Type:  &userAttributeType,
	}

	_, err = client.UpdateUserAttribute(userAttributeID, writeUserAttribute, "", nil)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceUserAttributeRead(ctx, d, m)
}

func resourceUserAttributeDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	userAttributeID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.DeleteUserAttribute(userAttributeID, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
