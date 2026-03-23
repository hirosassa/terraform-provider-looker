package looker

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
				Type:     schema.TypeString,
				Required: true,
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

	if folder.ParentId != nil {
		if err = d.Set("parent_id", *folder.ParentId); err != nil {
			return diag.FromErr(err)
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

	return resourceFolderRead(ctx, d, m)
}

func resourceFolderDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	folderID := d.Id()
	folderName := d.Get("name").(string)

	_, err := client.DeleteFolder(folderID, nil)
	if err != nil {
		return diag.FromErr(wrapSDKError(err, "DeleteFolder", "folder", "name=%s, id=%s", folderName, folderID))
	}

	return nil
}
