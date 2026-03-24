package looker

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	apiclient "github.com/looker-open-source/sdk-codegen/go/sdk/v4"
)

func resourceFolder() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFolderCreate,
		ReadContext:   resourceFolderRead,
		UpdateContext: resourceFolderUpdate,
		DeleteContext: resourceFolderDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"parent_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"content_metadata_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"inherits": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether content inherits its access levels from parent. Set to false to manage access with looker_content_metadata_access.",
			},
		},
	}
}

func resourceFolderCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)
	folderName := d.Get("name").(string)

	writeFolder := apiclient.CreateFolder{
		Name:     folderName,
		ParentId: d.Get("parent_id").(string),
	}

	folder, err := client.CreateFolder(writeFolder, nil)
	if err != nil {
		return diag.FromErr(wrapSDKError(err, "CreateFolder", "folder", "%s", folderName))
	}

	if folder.Id == nil {
		return diag.Errorf("Folder ID not returned from API")
	}

	d.SetId(*folder.Id)
	if folder.ContentMetadataId != nil {
		if err = d.Set("content_metadata_id", *folder.ContentMetadataId); err != nil {
			return diag.FromErr(err)
		}

		inherits := d.Get("inherits").(bool)
		_, err = client.UpdateContentMetadata(*folder.ContentMetadataId, apiclient.WriteContentMeta{
			Inherits: &inherits,
		}, nil)
		if err != nil {
			return diag.FromErr(wrapSDKError(err, "UpdateContentMetadata", "folder", "content_metadata_id=%s", *folder.ContentMetadataId))
		}
	}

	return resourceFolderRead(ctx, d, m)
}

func resourceFolderRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	folderID := d.Id()

	folder, err := client.Folder(folderID, "", nil)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			d.SetId("")
			return nil
		}
		return diag.FromErr(wrapSDKError(err, "Folder", "folder", "%s", folderID))
	}

	if err = d.Set("name", folder.Name); err != nil {
		return diag.FromErr(err)
	}

	if folder.ParentId == nil {
		return diag.Errorf("folder %s has no parent_id; root-level folders are not supported", folderID)
	}
	if err = d.Set("parent_id", *folder.ParentId); err != nil {
		return diag.FromErr(err)
	}

	if folder.ContentMetadataId != nil {
		if err = d.Set("content_metadata_id", *folder.ContentMetadataId); err != nil {
			return diag.FromErr(err)
		}

		contentMeta, err := client.ContentMetadata(*folder.ContentMetadataId, "", nil)
		if err != nil {
			return diag.FromErr(wrapSDKError(err, "ContentMetadata", "folder", "content_metadata_id=%s", *folder.ContentMetadataId))
		}
		if contentMeta.Inherits != nil {
			if err = d.Set("inherits", *contentMeta.Inherits); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return nil
}

func resourceFolderUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	folderID := d.Id()

	updateFolder := apiclient.UpdateFolder{}
	hasChanges := false

	if d.HasChange("name") {
		folderName := d.Get("name").(string)
		updateFolder.Name = &folderName
		hasChanges = true
	}

	if d.HasChange("parent_id") {
		parentID := d.Get("parent_id").(string)
		updateFolder.ParentId = &parentID
		hasChanges = true
	}

	if hasChanges {
		_, err := client.UpdateFolder(folderID, updateFolder, nil)
		if err != nil {
			return diag.FromErr(wrapSDKError(err, "UpdateFolder", "folder", "id=%s", folderID))
		}
	}

	if d.HasChange("inherits") {
		contentMetadataID := d.Get("content_metadata_id").(string)
		if contentMetadataID != "" {
			inherits := d.Get("inherits").(bool)
			_, err := client.UpdateContentMetadata(contentMetadataID, apiclient.WriteContentMeta{
				Inherits: &inherits,
			}, nil)
			if err != nil {
				return diag.FromErr(wrapSDKError(err, "UpdateContentMetadata", "folder", "content_metadata_id=%s", contentMetadataID))
			}
		}
	}

	return resourceFolderRead(ctx, d, m)
}

func resourceFolderDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	folderID := d.Id()
	folderName := d.Get("name").(string)

	_, err := client.DeleteFolder(folderID, nil)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return nil
		}
		return diag.FromErr(wrapSDKError(err, "DeleteFolder", "folder", "name=%s, id=%s", folderName, folderID))
	}

	return nil
}
