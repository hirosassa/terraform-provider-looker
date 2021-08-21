package looker

import (
	"context"
	"strconv"

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
		return diag.FromErr(err)
	}

	modelSetID := *modelSet.Id
	d.SetId(strconv.Itoa(int(modelSetID)))

	return resourceModelSetRead(ctx, d, m)
}

func resourceModelSetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	modelSetID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	modelSet, err := client.ModelSet(modelSetID, "", nil)
	if err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("name", modelSet.Name); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("models", modelSet.Models); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceModelSetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	modelSetID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}
	modelSetName := d.Get("name").(string)
	var modelNames []string
	for _, modelName := range d.Get("models").(*schema.Set).List() {
		modelNames = append(modelNames, modelName.(string))
	}
	writeModelSet := apiclient.WriteModelSet{
		Name:   &modelSetName,
		Models: &modelNames,
	}
	_, err = client.UpdateModelSet(modelSetID, writeModelSet, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceModelSetRead(ctx, d, m)
}

func resourceModelSetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	modelSetID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.DeleteModelSet(modelSetID, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
