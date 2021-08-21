package looker

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/looker-open-source/sdk-codegen/go/rtl"
	apiclient "github.com/looker-open-source/sdk-codegen/go/sdk/v4"
)

const (
	defaultAPIVersion = "4.0"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"client_id": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("LOOKER_API_CLIENT_ID", nil),
				Description: "Client ID to authenticate with Looker",
			},
			"client_secret": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("LOOKER_API_CLIENT_SECRET", nil),
				Description: "Client Secret to authenticate with Looker",
			},
			"base_url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("LOOKER_API_BASE_URL", nil),
				Description: "Looker API Base URL",
			},
			"api_version": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("LOOKER_API_VERSION", defaultAPIVersion),
			},
			"verify_ssl": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("LOOKER_VERIFY_SSL", true),
			},
			"timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("LOOKER_TIMEOUT", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"looker_user":             resourceUser(),
			"looker_user_roles":       resourceUserRoles(),
			"looker_permission_set":   resourcePermissionSet(),
			"looker_model_set":        resourceModelSet(),
			"looker_group":            resourceGroup(),
			"looker_group_membership": resourceGroupMembership(),
			"looker_role":             resourceRole(),
			"looker_role_groups":      resourceRoleGroups(),
			"looker_user_attribute":   resourceUserAttribute(),
			"looker_connection":       resourceConnection(),
		},

		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	baseUrl := d.Get("base_url").(string)
	clientID := d.Get("client_id").(string)
	clientSecret := d.Get("client_secret").(string)
	apiVersion := d.Get("api_version").(string)
	timeout := d.Get("timeout").(int)

	apiSettings := rtl.ApiSettings{
		BaseUrl:      baseUrl,
		ClientId:     clientID,
		ClientSecret: clientSecret,
		ApiVersion:   apiVersion,
		VerifySsl:    d.Get("verify_ssl").(bool),
		Timeout:      int32(timeout),
	}
	authSession := rtl.NewAuthSession(apiSettings)
	client := apiclient.NewLookerSDK(authSession)

	return client, diag.Diagnostics{}
}
