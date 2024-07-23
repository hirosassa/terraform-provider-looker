package looker

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	apiclient "github.com/looker-open-source/sdk-codegen/go/sdk/v4"
)

func resourceProjectGitRepository() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProjectGitRepositoryCreate,
		ReadContext:   resourceProjectGitRepositoryRead,
		UpdateContext: resourceProjectGitRepositoryUpdate,
		DeleteContext: resourceProjectGitRepositoryDelete,
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
			"git_service_name": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "bare",
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{"bare", "github"}, false),
			},
			"git_remote_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"git_deploy_key": {
				Type:      schema.TypeString,
				Computed:  true,
			},
		},
	}
}

func resourceProjectGitRepositoryCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)
	// Set client to dev workspace
	if err := updateApiSessionWorkspaceId(client, "dev"); err != nil {
		return diag.FromErr(err)
	}

	// Update the project git service name
	projectName := d.Get("project").(string)
	d.SetId(projectName)
	gitServiceName := d.Get("git_service_name").(string)
	var gitRemoteUrl string
	if gitServiceName != "bare" {
		gitRemoteUrl = d.Get("git_remote_url").(string)
	}
	updateProject := apiclient.WriteProject{
		GitServiceName: &gitServiceName,
		GitRemoteUrl: &gitRemoteUrl,
	}
	_, err := client.UpdateProject(projectName, updateProject, "name,git_service_name,git_remote_url", nil)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceProjectGitRepositoryRead(ctx, d, m)
}

func resourceProjectGitRepositoryRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)
	// Set client to dev workspace
	if err := updateApiSessionWorkspaceId(client, "dev"); err != nil {
		return diag.FromErr(err)
	}

	project, err := client.Project(d.Id(), "name,git_service_name,git_remote_url", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("git_service_name", project.GitServiceName)
	d.Set("git_remote_url", project.GitRemoteUrl)

	return nil
}

func resourceProjectGitRepositoryUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)
	// Set client to dev workspace
	if err := updateApiSessionWorkspaceId(client, "dev"); err != nil {
		return diag.FromErr(err)
	}

	gitServiceName := d.Get("git_service_name").(string)
	var gitRemoteUrl string
	if gitServiceName != "bare" {
		gitRemoteUrl = d.Get("git_remote_url").(string)
	}
	writeProject := apiclient.WriteProject{
		GitServiceName:  &gitServiceName,
		GitRemoteUrl:    &gitRemoteUrl,
	}
	_, err := client.UpdateProject(d.Id(), writeProject, "name,git_service_name,git_remote_url", nil)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceProjectGitRepositoryDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}

