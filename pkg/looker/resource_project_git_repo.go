package looker

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apiclient "github.com/looker-open-source/sdk-codegen/go/sdk/v4"
)

func resourceProjectGitRepo() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProjectGitRepoCreate,
		ReadContext:   resourceProjectGitRepoRead,
		UpdateContext: resourceProjectGitRepoUpdate,
		DeleteContext: resourceProjectGitRepoDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"git_remote_url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"git_username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"git_password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"git_service_name": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "github",
			},
			"pull_request_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "off",
			},
		},
	}
}

func setProjectGitRepo(d *schema.ResourceData, m interface{}, create bool) error {
	client := m.(*apiclient.LookerSDK)
	workspace := "dev"

	writeApiSession := apiclient.WriteApiSession{
		WorkspaceId: &workspace,
	}

	_, errSession := client.UpdateSession(writeApiSession, nil)
	if errSession != nil {
		return errSession
	}

	projectId := d.Get("project_id").(string)
	gitRemoteURL := d.Get("git_remote_url").(string)
	gitUsername := d.Get("git_username").(string)
	gitPassword := d.Get("git_password").(string)
	gitServiceName := d.Get("git_service_name").(string)
	pullRequestMode := apiclient.PullRequestMode(d.Get("pull_request_mode").(string))

	writeProject := apiclient.WriteProject{
		GitRemoteUrl:    &gitRemoteURL,
		GitUsername:     &gitUsername,
		GitPassword:     &gitPassword,
		GitServiceName:  &gitServiceName,
		PullRequestMode: &pullRequestMode,
	}

	_, err := client.UpdateProject(projectId, writeProject, "", nil)
	if err != nil {
		return err
	}
	return nil

}

func resourceProjectGitRepoCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	error := setProjectGitRepo(d, m, true)
	if error != nil {
		return diag.FromErr(error)
	}
	d.SetId(d.Get("project_id").(string))
	return resourceProjectGitRepoRead(ctx, d, m)
}

func resourceProjectGitRepoRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)
	workspace := "dev"

	writeApiSession := apiclient.WriteApiSession{
		WorkspaceId: &workspace,
	}

	_, errSession := client.UpdateSession(writeApiSession, nil)
	if errSession != nil {
		return diag.FromErr(errSession)
	}

	projectId := d.Get("project_id").(string)

	project, err := client.Project(projectId, "", nil)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("project_id", *project.Id)
	d.Set("git_remote_url", *project.GitRemoteUrl)
	d.Set("git_username", *project.GitUsername)
	d.Set("git_service_name", *project.GitServiceName)
	d.Set("pull_request_mode", string(*project.PullRequestMode))

	return nil
}

func resourceProjectGitRepoUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	error := setProjectGitRepo(d, m, true)
	if error != nil {
		return diag.FromErr(error)
	}
	return resourceProjectGitRepoRead(ctx, d, m)
}

func resourceProjectGitRepoDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// TODO: Looker doesn't appear to support deleting projects from the API
	return nil
}
