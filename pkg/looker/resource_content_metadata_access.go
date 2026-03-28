package looker

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	apiclient "github.com/looker-open-source/sdk-codegen/go/sdk/v4"
)

func resourceContentMetadataAccess() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceContentMetadataAccessCreate,
		ReadContext:   resourceContentMetadataAccessRead,
		UpdateContext: resourceContentMetadataAccessUpdate,
		DeleteContext: resourceContentMetadataAccessDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				idParts := strings.SplitN(d.Id(), "/", 2)
				if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
					return nil, fmt.Errorf("unexpected format of ID (%q), expected content_metadata_id/access_id", d.Id())
				}

				if err := d.Set("content_metadata_id", idParts[0]); err != nil {
					return nil, err
				}

				d.SetId(idParts[1])

				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"content_metadata_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"permission_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"view", "edit"}, false),
			},
			"user_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"group_id"},
				AtLeastOneOf:  []string{"user_id", "group_id"},
			},
			"group_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"user_id"},
				AtLeastOneOf:  []string{"user_id", "group_id"},
			},
		},
	}
}

func resourceContentMetadataAccessCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	contentMetadataID := d.Get("content_metadata_id").(string)
	permissionType := apiclient.PermissionType(d.Get("permission_type").(string))

	access := apiclient.ContentMetaGroupUser{
		ContentMetadataId: &contentMetadataID,
		PermissionType:    &permissionType,
	}

	if v, ok := d.GetOk("user_id"); ok {
		userID := v.(string)
		access.UserId = &userID
	}
	if v, ok := d.GetOk("group_id"); ok {
		groupID := v.(string)
		access.GroupId = &groupID
	}

	result, err := client.CreateContentMetadataAccess(access, false, nil)
	if err != nil {
		return diag.FromErr(wrapSDKError(err, "CreateContentMetadataAccess", "content_metadata_access", "metadata_id=%s", contentMetadataID))
	}

	if result.Id == nil {
		return diag.Errorf("Content Metadata Access ID not returned from API")
	}

	d.SetId(*result.Id)

	return resourceContentMetadataAccessRead(ctx, d, m)
}

func resourceContentMetadataAccessRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	accessID := d.Id()
	contentMetadataID := d.Get("content_metadata_id").(string)

	accesses, err := client.AllContentMetadataAccesses(contentMetadataID, "", nil)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			d.SetId("")
			return nil
		}
		return diag.FromErr(wrapSDKError(err, "AllContentMetadataAccesses", "content_metadata_access", "metadata_id=%s", contentMetadataID))
	}

	var found *apiclient.ContentMetaGroupUser
	for _, a := range accesses {
		if a.Id != nil && *a.Id == accessID {
			found = &a
			break
		}
	}

	if found == nil {
		d.SetId("")
		return nil
	}

	if found.ContentMetadataId == nil {
		d.SetId("")
		return nil
	}
	if err = d.Set("content_metadata_id", *found.ContentMetadataId); err != nil {
		return diag.FromErr(err)
	}

	if found.PermissionType == nil {
		d.SetId("")
		return nil
	}
	if err = d.Set("permission_type", string(*found.PermissionType)); err != nil {
		return diag.FromErr(err)
	}
	if found.UserId != nil {
		if err = d.Set("user_id", *found.UserId); err != nil {
			return diag.FromErr(err)
		}
	} else {
		if err = d.Set("user_id", ""); err != nil {
			return diag.FromErr(err)
		}
	}
	if found.GroupId != nil {
		if err = d.Set("group_id", *found.GroupId); err != nil {
			return diag.FromErr(err)
		}
	} else {
		if err = d.Set("group_id", ""); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func resourceContentMetadataAccessUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	accessID := d.Id()

	if d.HasChange("permission_type") {
		permissionType := apiclient.PermissionType(d.Get("permission_type").(string))
		access := apiclient.ContentMetaGroupUser{
			PermissionType: &permissionType,
		}

		_, err := client.UpdateContentMetadataAccess(accessID, access, nil)
		if err != nil {
			return diag.FromErr(wrapSDKError(err, "UpdateContentMetadataAccess", "content_metadata_access", "id=%s", accessID))
		}
	}

	return resourceContentMetadataAccessRead(ctx, d, m)
}

func resourceContentMetadataAccessDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	accessID := d.Id()

	_, err := client.DeleteContentMetadataAccess(accessID, nil)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return nil
		}
		return diag.FromErr(wrapSDKError(err, "DeleteContentMetadataAccess", "content_metadata_access", "id=%s", accessID))
	}

	return nil
}
