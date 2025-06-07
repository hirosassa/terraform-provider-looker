package looker

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	apiclient "github.com/looker-open-source/sdk-codegen/go/sdk/v4"
)

func resourceProject() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProjectCreate,
		ReadContext:   resourceProjectRead,
		DeleteContext: resourceProjectDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringDoesNotContainAny(" "),
			},
		},
	}
}

func resourceProjectCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)
	// Set client to dev workspace
	if err := updateApiSessionWorkspaceId(client, "dev"); err != nil {
		return diag.FromErr(err)
	}

	// Create the project
	projectName := d.Get("name").(string)
	writeProject := apiclient.WriteProject{
		Name:            &projectName,
	}
	project, err := client.CreateProject(writeProject, nil)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(*project.Id)

	return nil
}

func resourceProjectRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)
	// Set client to dev workspace
	if err := updateApiSessionWorkspaceId(client, "dev"); err != nil {
		return diag.FromErr(err)
	}

	project, err := client.Project(d.Id(),  "name", nil)
	if err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("name", project.Name); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceProjectDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	diags = append(diags, diag.Diagnostic{
		Severity: diag.Warning,
		Summary: "Deleting a Looker project is currently impossible via the API.",
		Detail: "In order to delete a Looker project, and to avoid naming conflicts the usage of random_id is advised.\nSee also: https://docs.looker.com/data-modeling/getting-started/manage-projects#deleting_a_project",
	})
	d.SetId("")
	return diags
}

