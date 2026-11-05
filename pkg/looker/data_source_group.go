package looker

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apiclient "github.com/looker-open-source/sdk-codegen/go/sdk/v4"
)

func dataSourceGroup() *schema.Resource {
	return &schema.Resource{
		Description: "Looks up a Looker group by name.",
		ReadContext: dataSourceLookupRead("name", "group",
			func(client *apiclient.LookerSDK, name string) ([]apiclient.Group, error) {
				return client.SearchGroups(apiclient.RequestSearchGroups{Name: &name}, nil)
			},
			func(x apiclient.Group) *string { return x.Name },
			func(x apiclient.Group) *string { return x.Id },
			flattenGroup,
		),
		Schema: dataSourceSchemaFromResourceSchema(resourceGroup().Schema, "name", "The name of the group."),
	}
}
