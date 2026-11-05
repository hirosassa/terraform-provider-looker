package looker

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apiclient "github.com/looker-open-source/sdk-codegen/go/sdk/v4"
)

func resourceModelSet() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceModelSetCreate,
		ReadContext:   resourceModelSetRead,
		UpdateContext: resourceModelSetUpdate,
		DeleteContext: resourceModelSetDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"models": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceModelSetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	modelSetName := d.Get("name").(string)

	var modelNames []string
	for _, modelName := range d.Get("models").(*schema.Set).List() {
		modelNames = append(modelNames, modelName.(string))
	}

	writeModelSet := apiclient.WriteModelSet{
		Name:   &modelSetName,
		Models: &modelNames,
	}

	modelSet, err := client.CreateModelSet(writeModelSet, nil)
	if err != nil {
		return diag.FromErr(wrapSDKError(err, "CreateModelSet", "model_set", "%s", modelSetName))
	}

	modelSetID := *modelSet.Id
	d.SetId(modelSetID)

	return resourceModelSetRead(ctx, d, m)
}

func resourceModelSetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	modelSetID := d.Id()

	modelSet, err := client.ModelSet(modelSetID, "", nil)
	if err != nil {
		return diag.FromErr(wrapSDKError(err, "ModelSet", "model_set", "%s", modelSetID))
	}

	return diag.FromErr(flattenModelSet(modelSet, d))
}

func flattenModelSet(modelSet apiclient.ModelSet, d *schema.ResourceData) error {
	if err := d.Set("name", modelSet.Name); err != nil {
		return err
	}
	return d.Set("models", modelSet.Models)
}

func resourceModelSetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	modelSetID := d.Id()
	modelSetName := d.Get("name").(string)
	var modelNames []string
	for _, modelName := range d.Get("models").(*schema.Set).List() {
		modelNames = append(modelNames, modelName.(string))
	}
	writeModelSet := apiclient.WriteModelSet{
		Name:   &modelSetName,
		Models: &modelNames,
	}
	_, err := client.UpdateModelSet(modelSetID, writeModelSet, nil)
	if err != nil {
		return diag.FromErr(wrapSDKError(err, "UpdateModelSet", "model_set", "name=%s, id=%s", modelSetName, modelSetID))
	}

	return resourceModelSetRead(ctx, d, m)
}

func resourceModelSetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	modelSetID := d.Id()
	modelSetName := d.Get("name").(string)

	_, err := client.DeleteModelSet(modelSetID, nil)
	if err != nil {
		return diag.FromErr(wrapSDKError(err, "DeleteModelSet", "model_set", "name=%s, id=%s", modelSetName, modelSetID))
	}

	return nil
}
