package looker

import (
	"context"
	"strconv"

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
				Type:     schema.TypeInt,
				Required: true,
			},
			"group_id": {
				Type:     schema.TypeInt,
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

	groupID := int64(d.Get("group_id").(int))
	userAttributeID := int64(d.Get("user_attribute_id").(int))
	value := d.Get("value").(string)

	body := apiclient.UserAttributeGroupValue{
		GroupId:         &groupID,
		UserAttributeId: &userAttributeID,
		Value:           &value,
	}
	userAttributeGroupValue, err := client.UpdateUserAttributeGroupValue(groupID, userAttributeID, body, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	groupIDString := strconv.Itoa(int(*userAttributeGroupValue.GroupId))
	userAttributeIDString := strconv.Itoa(int(*userAttributeGroupValue.UserAttributeId))
	id := buildTwoPartID(&groupIDString, &userAttributeIDString)

	d.SetId(id)

	return resourceUserAttributeGroupValueRead(ctx, d, m)
}

func resourceUserAttributeGroupValueRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	groupIDString, userAttributeIDString, err := parseTwoPartID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	groupID, err := strconv.ParseInt(groupIDString, 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}
	userAttributeID, err := strconv.ParseInt(userAttributeIDString, 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	userAttributeGroupValues, err := client.AllUserAttributeGroupValues(userAttributeID, "", nil)
	if err != nil {
		return diag.FromErr(err)
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

	groupIDString, userAttributeIDString, err := parseTwoPartID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	groupID, err := strconv.ParseInt(groupIDString, 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}
	userAttributeID, err := strconv.ParseInt(userAttributeIDString, 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	value := d.Get("value").(string)

	body := apiclient.UserAttributeGroupValue{
		GroupId:         &groupID,
		UserAttributeId: &userAttributeID,
		Value:           &value,
	}
	_, err = client.UpdateUserAttributeGroupValue(groupID, userAttributeID, body, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceUserAttributeGroupValueRead(ctx, d, m)
}

func resourceUserAttributeGroupValueDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	groupIDString, userAttributeIDString, err := parseTwoPartID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	groupID, err := strconv.ParseInt(groupIDString, 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}
	userAttributeID, err := strconv.ParseInt(userAttributeIDString, 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.DeleteUserAttributeGroupValue(groupID, userAttributeID, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
