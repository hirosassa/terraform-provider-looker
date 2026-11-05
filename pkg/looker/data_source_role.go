package looker

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apiclient "github.com/looker-open-source/sdk-codegen/go/sdk/v4"
)

func dataSourceRole() *schema.Resource {
	return &schema.Resource{
		Description: "Looks up a Looker role by name.",
		ReadContext: dataSourceLookupRead("name", "role",
			func(client *apiclient.LookerSDK, name string) ([]apiclient.Role, error) {
				return client.SearchRoles(apiclient.RequestSearchRoles{Name: &name}, nil)
			},
			func(x apiclient.Role) *string { return x.Name },
			func(x apiclient.Role) *string { return x.Id },
			flattenRole,
		),
		Schema: dataSourceSchemaFromResourceSchema(resourceRole().Schema, "name", "The name of the role."),
	}
}
