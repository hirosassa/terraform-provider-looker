package looker

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apiclient "github.com/looker-open-source/sdk-codegen/go/sdk/v4"
)

func dataSourceUser() *schema.Resource {
	return &schema.Resource{
		Description: "Looks up a single Looker user by email address. " +
			"Use `looker_users` instead when you need to list all users.",
		ReadContext: dataSourceLookupRead("email", "user",
			func(client *apiclient.LookerSDK, email string) ([]apiclient.User, error) {
				return client.SearchUsers(apiclient.RequestSearchUsers{Email: &email}, nil)
			},
			func(x apiclient.User) *string { return x.Email },
			func(x apiclient.User) *string { return x.Id },
			flattenUser,
		),
		Schema: dataSourceUserSchema(),
	}
}

func dataSourceUserSchema() map[string]*schema.Schema {
	s := dataSourceSchemaFromResourceSchema(resourceUser().Schema, "email", "The email address of the user.")
	// Only meaningful when creating a user via the resource.
	delete(s, "send_setup_link_on_create")
	// The resource's description is about omitting the argument, which makes no sense for a data source.
	s["can_manage_api3_creds"].Description = "Whether the user can create, view, and delete API keys for their own account. " +
		"This field is marked as experimental in the Looker API and may not be available on your instance."
	return s
}
