package looker

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	apiclient "github.com/looker-open-source/sdk-codegen/go/sdk/v4"
)

func resourceProjectGitDeployKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProjectGitDeployKeyCreate,
		ReadContext:   resourceProjectGitDeployKeyRead,
		DeleteContext: resourceProjectGitDeployKeyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringDoesNotContainAny(" "),
			},
			"git_deploy_key": {
				Type:      schema.TypeString,
				Computed:  true,
			},
		},
	}
}

func resourceProjectGitDeployKeyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)
	// Set client to dev workspace
	if err := updateApiSessionWorkspaceId(client, "dev"); err != nil {
		return diag.FromErr(err)
	}

	projectName := d.Get("project").(string)
	d.SetId(projectName)
	gitDeployKey, err := client.CreateGitDeployKey(projectName, nil)
	if err != nil {
		// Graceful fallback if there is already a git deploy key setup
		if strings.Contains(err.Error(), "409") {
			return resourceProjectGitDeployKeyRead(ctx, d, m)
		}
		return diag.FromErr(err)
	}
	d.Set("git_deploy_key", gitDeployKey)

	return nil
}

func resourceProjectGitDeployKeyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)
	// Set client to dev workspace
	if err := updateApiSessionWorkspaceId(client, "dev"); err != nil {
		return diag.FromErr(err)
	}

	gitDeployKey, err := client.GitDeployKey(d.Id(), nil)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("git_deploy_key", gitDeployKey)

	return nil
}

func resourceProjectGitDeployKeyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}

