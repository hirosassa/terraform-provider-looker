package looker

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/looker-open-source/sdk-codegen/go/rtl"
	apiclient "github.com/looker-open-source/sdk-codegen/go/sdk/v4"
)

func resourceUserAttributeUserValue() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserAttributeUserValueCreate,
		ReadContext:   resourceUserAttributeUserValueRead,
		UpdateContext: resourceUserAttributeUserValueUpdate,
		DeleteContext: resourceUserAttributeUserValueDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"user_attribute_id": {
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

func resourceUserAttributeUserValueCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)
	userID := int64(d.Get("user_id").(int))
	userAttributeID := int64(d.Get("user_attribute_id").(int))
	userAttributeValue := d.Get("value").(string)

	body := apiclient.WriteUserAttributeWithValue{
		Value: &userAttributeValue,
	}

	userAttributeWithValue, err := client.SetUserAttributeUserValue(userID, userAttributeID, body, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	userIDString := strconv.Itoa(int(*userAttributeWithValue.UserId))
	userAttributeIDString := strconv.Itoa(int(*userAttributeWithValue.UserAttributeId))
	id := buildTwoPartID(&userIDString, &userAttributeIDString)

	d.SetId(id)

	return resourceUserAttributeUserValueRead(ctx, d, m)
}

func resourceUserAttributeUserValueRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	userIDString, userAttributeIDString, err := parseTwoPartID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	userID, err := strconv.ParseInt(userIDString, 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}
	userAttributeID, err := strconv.ParseInt(userAttributeIDString, 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	userAttributeIDs := rtl.DelimInt64{userAttributeID}
	request := apiclient.RequestUserAttributeUserValues{
		UserId:           userID,
		UserAttributeIds: &userAttributeIDs,
	}

	userAttributeUserValues, err := client.UserAttributeUserValues(request, nil)
	if err != nil {
		return diag.FromErr(err)
	}
	if len(userAttributeUserValues) != 1 { // the number of the result should be one
		return diag.FromErr(err)
	}

	if err = d.Set("user_id", userAttributeUserValues[0].UserId); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("user_attribute_id", userAttributeUserValues[0].UserAttributeId); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("value", userAttributeUserValues[0].Value); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceUserAttributeUserValueUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	userIDString, userAttributeIDString, err := parseTwoPartID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	userID, err := strconv.ParseInt(userIDString, 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}
	userAttributeID, err := strconv.ParseInt(userAttributeIDString, 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	userAttributeValue := d.Get("value").(string)
	body := apiclient.WriteUserAttributeWithValue{
		Value: &userAttributeValue,
	}

	_, err = client.SetUserAttributeUserValue(userID, userAttributeID, body, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceUserAttributeUserValueRead(ctx, d, m)
}

func resourceUserAttributeUserValueDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	userIDString, userAttributeIDString, err := parseTwoPartID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	userID, err := strconv.ParseInt(userIDString, 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}
	userAttributeID, err := strconv.ParseInt(userAttributeIDString, 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	body := apiclient.WriteUserAttributeWithValue{
		Value: nil,
	}

	// fill the target value null
	_, err = client.SetUserAttributeUserValue(userID, userAttributeID, body, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
