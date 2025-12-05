package looker

import (
	"context"

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
				Optional: true,
				Default:  nil,
			},
		},
	}
}

func resourceFolderCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	folderName := d.Get("name").(string)
	parentId := d.Get("parent_id").(string)

	createFolder := apiclient.CreateFolder{
		Name:     folderName,
		ParentId: parentId,
	}

	folder, err := client.CreateFolder(createFolder, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	folderId := *folder.Id
	d.SetId(folderId)

	return resourceFolderRead(ctx, d, m)
}

func resourceFolderRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	folderId := d.Id()

	folder, err := client.Folder(folderId, "", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("name", folder.Name); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("parent_id", folder.ParentId); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceFolderUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	folderId := d.Id()
	folderName := d.Get("name").(string)
	parentId := d.Get("parent_id").(string)

	updateFolder := apiclient.UpdateFolder{
		Name:     &folderName,
		ParentId: &parentId,
	}

	_, err := client.UpdateFolder(folderId, updateFolder, nil)
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceFolderRead(ctx, d, m)
}

func resourceFolderDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)
	folderId := d.Id()

	_, err := client.DeleteFolder(folderId, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil

}
