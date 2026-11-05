package looker

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apiclient "github.com/looker-open-source/sdk-codegen/go/sdk/v4"
)

func dataSourcePermissionSet() *schema.Resource {
	return &schema.Resource{
		Description: "Looks up a Looker permission set by name.",
		ReadContext: dataSourceLookupRead("name", "permission set",
			func(client *apiclient.LookerSDK, name string) ([]apiclient.PermissionSet, error) {
				return client.SearchPermissionSets(apiclient.RequestSearchPermissionSets{Name: &name}, nil)
			},
			func(x apiclient.PermissionSet) *string { return x.Name },
			func(x apiclient.PermissionSet) *string { return x.Id },
			flattenPermissionSet,
		),
		Schema: dataSourceSchemaFromResourceSchema(resourcePermissionSet().Schema, "name", "The name of the permission set."),
	}
}
