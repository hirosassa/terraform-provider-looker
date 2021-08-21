package looker

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apiclient "github.com/looker-open-source/sdk-codegen/go/sdk/v4"
)

func resourceLookMLModel() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLookMLModelCreate,
		ReadContext:   resourceLookMLModelRead,
		UpdateContext: resourceLookMLModelUpdate,
		DeleteContext: resourceLookMLModelDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"allowed_db_connection_names": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"project_name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceLookMLModelCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	body, err := expandWriteLookmlModel(d)
	if err != nil {
		return diag.FromErr(err)
	}

	result, err := client.CreateLookmlModel(*body, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*result.Name)

	return resourceLookMLModelRead(ctx, d, m)
}

func resourceLookMLModelRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	model, err := client.LookmlModel(d.Id(), "", nil)
	if err != nil {
		if strings.Contains(err.Error(), "Not found") {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	return diag.FromErr(flattenLookMLModel(model, d))
}

func resourceLookMLModelUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	body, err := expandWriteLookmlModel(d)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.UpdateLookmlModel(d.Id(), *body, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceLookMLModelRead(ctx, d, m)
}

func resourceLookMLModelDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	_, err := client.DeleteLookmlModel(d.Id(), nil)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func expandWriteLookmlModel(d *schema.ResourceData) (*apiclient.WriteLookmlModel, error) {
	modelName := d.Get("name").(string)
	projectName := d.Get("project_name").(string)
	var connections []string
	for _, modelName := range d.Get("allowed_db_connection_names").(*schema.Set).List() {
		connections = append(connections, modelName.(string))
	}
	return &apiclient.WriteLookmlModel{
		Name:                     &modelName,
		ProjectName:              &projectName,
		AllowedDbConnectionNames: &connections,
	}, nil
}

func flattenLookMLModel(model apiclient.LookmlModel, d *schema.ResourceData) error {
	if err := d.Set("name", *model.Name); err != nil {
		return err
	}
	if err := d.Set("project_name", *model.ProjectName); err != nil {
		return err
	}
	if err := d.Set("allowed_db_connection_names", *model.AllowedDbConnectionNames); err != nil {
		return err
	}
	return nil
}
