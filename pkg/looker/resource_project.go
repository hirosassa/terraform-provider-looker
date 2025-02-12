package looker

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apiclient "github.com/looker-open-source/sdk-codegen/go/sdk/v4"
)

func resourceProject() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProjectCreate,
		ReadContext:   resourceProjectRead,
		UpdateContext: resourceProjectUpdate,
		DeleteContext: resourceProjectDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceProjectCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)
	workspace := "dev"

	writeApiSession := apiclient.WriteApiSession{
		WorkspaceId: &workspace,
	}

	_, errSession := client.UpdateSession(writeApiSession, nil)
	if errSession != nil {
		return diag.FromErr(errSession)
	}

	name := d.Get("name").(string)
	writeProject := apiclient.WriteProject{
		Name: &name,
	}

	_, err := client.CreateProject(writeProject, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(name)

	return resourceProjectRead(ctx, d, m)

}

func resourceProjectRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)
	workspace := "dev"

	writeApiSession := apiclient.WriteApiSession{
		WorkspaceId: &workspace,
	}
	_, errSession := client.UpdateSession(writeApiSession, nil)
	if errSession != nil {
		return diag.FromErr(errSession)
	}

	projectId := d.Id()

	project, err := client.Project(projectId, "", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("name", project.Name); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceProjectUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)
	workspace := "dev"

	writeApiSession := apiclient.WriteApiSession{
		WorkspaceId: &workspace,
	}

	_, errSession := client.UpdateSession(writeApiSession, nil)
	if errSession != nil {
		return diag.FromErr(errSession)
	}

	projectId := d.Id()
	name := d.Get("name").(string)
	writeProject := apiclient.WriteProject{
		Name: &name,
	}

	_, err := client.UpdateProject(projectId, writeProject, "", nil)
	if err != nil {
	}

	d.SetId(name)

	return resourceProjectRead(ctx, d, m)
}

func resourceProjectDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// TODO: Looker doesn't appear to support deleting projects from the API
	return nil
}
