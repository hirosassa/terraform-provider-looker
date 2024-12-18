package looker

import (
	"context"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apiclient "github.com/looker-open-source/sdk-codegen/go/sdk/v4"
)

func dataSourceUsers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUsersRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier for the resource.",
			},
			"users": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"email": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"first_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_disabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceUsersRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)
	request := apiclient.RequestAllUsers{}
	users, err := client.AllUsers(request, nil)

	if err != nil {
		return diag.FromErr(err)
	}

	userList := make([]map[string]interface{}, len(users))
	userEmails := make([]string, len(users))
	for i, user := range users {
		userList[i] = map[string]interface{}{
			"id":          *user.Id,
			"email":       *user.Email,
			"first_name":  *user.FirstName,
			"last_name":   *user.LastName,
			"is_disabled": *user.IsDisabled,
		}
		userEmails[i] = *user.Email
	}

	if err := d.Set("users", userList); err != nil {
		return diag.FromErr(err)
	}

	// Generate a hash of the user emails to use as the resource ID
	sort.Strings(userEmails)
	d.SetId(hash(strings.Join(userEmails, ",")))

	return nil
}
