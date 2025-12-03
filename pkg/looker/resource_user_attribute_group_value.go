package looker

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apiclient "github.com/looker-open-source/sdk-codegen/go/sdk/v4"
)

func resourceUserAttributeGroupValue() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserAttributeGroupValueCreate,
		ReadContext:   resourceUserAttributeGroupValueRead,
		UpdateContext: resourceUserAttributeGroupValueUpdate,
		DeleteContext: resourceUserAttributeGroupValueDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"user_attribute_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"value": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceUserAttributeGroupValueCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	groupID := d.Get("group_id").(string)
	userAttributeID := d.Get("user_attribute_id").(string)
	value := d.Get("value").(string)

	log.Printf("[DEBUG] Create user attribute group value %s for %s by %s", userAttributeID, groupID, value)

	body := apiclient.UserAttributeGroupValue{
		GroupId:         &groupID,
		UserAttributeId: &userAttributeID,
		Value:           &value,
	}
	userAttributeGroupValue, err := client.UpdateUserAttributeGroupValue(groupID, userAttributeID, body, nil)
	if err != nil {
		return diag.FromErr(wrapSDKError(err, "UpdateUserAttributeGroupValue", "user_attribute_group_value", "%s:%s", groupID, userAttributeID))
	}

	groupIDString := *userAttributeGroupValue.GroupId
	userAttributeIDString := *userAttributeGroupValue.UserAttributeId
	id := buildTwoPartID(&groupIDString, &userAttributeIDString)

	d.SetId(id)

	return resourceUserAttributeGroupValueRead(ctx, d, m)
}

func resourceUserAttributeGroupValueRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	groupID, userAttributeID, err := parseTwoPartID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Read user attribute group value %s for %s", userAttributeID, groupID)

	userAttributeGroupValues, err := client.AllUserAttributeGroupValues(userAttributeID, "", nil)
	if err != nil {
		return diag.FromErr(wrapSDKError(err, "AllUserAttributeGroupValues", "user_attribute_group_value", "%s", userAttributeID))
	}

	var userAttributeGroupValue apiclient.UserAttributeGroupValue
	for _, groupValue := range userAttributeGroupValues {
		if *groupValue.GroupId == groupID {
			userAttributeGroupValue = groupValue
			break
		}
	}

	if err = d.Set("group_id", userAttributeGroupValue.GroupId); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("user_attribute_id", userAttributeGroupValue.UserAttributeId); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("value", userAttributeGroupValue.Value); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceUserAttributeGroupValueUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	groupID, userAttributeID, err := parseTwoPartID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	value := d.Get("value").(string)

	log.Printf("[DEBUG] Update user attribute group value %s for %s by %s", userAttributeID, groupID, value)

	body := apiclient.UserAttributeGroupValue{
		GroupId:         &groupID,
		UserAttributeId: &userAttributeID,
		Value:           &value,
	}
	_, err = client.UpdateUserAttributeGroupValue(groupID, userAttributeID, body, nil)
	if err != nil {
		return diag.FromErr(wrapSDKError(err, "UpdateUserAttributeGroupValue", "user_attribute_group_value", "%s:%s", groupID, userAttributeID))
	}

	return resourceUserAttributeGroupValueRead(ctx, d, m)
}

func resourceUserAttributeGroupValueDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	groupID, userAttributeID, err := parseTwoPartID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Delete user attribute group value %s for %s", userAttributeID, groupID)

	err = client.DeleteUserAttributeGroupValue(groupID, userAttributeID, nil)
	if err != nil {
		log.Printf("[DEBUG] %+v", err)
		return diag.FromErr(wrapSDKError(err, "DeleteUserAttributeGroupValue", "user_attribute_group_value", "%s:%s", groupID, userAttributeID))
	}

	return nil
}
