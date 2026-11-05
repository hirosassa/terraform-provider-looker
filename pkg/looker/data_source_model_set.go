package looker

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apiclient "github.com/looker-open-source/sdk-codegen/go/sdk/v4"
)

func dataSourceModelSet() *schema.Resource {
	return &schema.Resource{
		Description: "Looks up a Looker model set by name.",
		ReadContext: dataSourceLookupRead("name", "model set",
			func(client *apiclient.LookerSDK, name string) ([]apiclient.ModelSet, error) {
				return client.SearchModelSets(apiclient.RequestSearchModelSets{Name: &name}, nil)
			},
			func(x apiclient.ModelSet) *string { return x.Name },
			func(x apiclient.ModelSet) *string { return x.Id },
			flattenModelSet,
		),
		Schema: dataSourceSchemaFromResourceSchema(resourceModelSet().Schema, "name", "The name of the model set."),
	}
}
