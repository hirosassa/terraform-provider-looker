package looker

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apiclient "github.com/looker-open-source/sdk-codegen/go/sdk/v4"
)

// format the strings into an id `a:b`
// this function borrowed from https://github.com/gitlabhq/terraform-provider-gitlab
func buildTwoPartID(a, b *string) string {
	return fmt.Sprintf("%s:%s", *a, *b)
}

// return the pieces of id `a:b` as a, b
// this function borrowed from https://github.com/gitlabhq/terraform-provider-gitlab
func parseTwoPartID(id string) (string, string, error) {
	parts := strings.SplitN(id, ":", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("unexpected ID format (%q). Expected project:key", id)
	}

	return parts[0], parts[1], nil
}

func expandStringListFromSet(set interface{}) []string {
	var strings []string
	for _, v := range set.(*schema.Set).List() {
		strings = append(strings, v.(string))
	}
	return strings
}

func flattenStringList(strings []string) []interface{} {
	vs := make([]interface{}, 0, len(strings))
	for _, v := range strings {
		vs = append(vs, v)
	}
	return vs
}

func flattenStringListToSet(strings []string) *schema.Set {
	return schema.NewSet(schema.HashString, flattenStringList(strings))
}

// func expandInt64ListFromSet(set interface{}) []int64 {
// 	var ints []int64
// 	for _, v := range set.(*schema.Set).List() {
// 		ints = append(ints, int64(v.(int)))
// 	}
// 	return ints
// }

// Looker's search APIs treat '_' and '%' as SQL LIKE wildcards ("a_b" also matches "axb"),
// so the results have to be narrowed down here.
// Case-insensitive because Looker rejects names/emails that differ only in case,
// which also guarantees at most one match.
func findExactMatch[T any](items []T, value string, getKey func(T) *string) (T, bool) {
	for _, item := range items {
		if k := getKey(item); k != nil && strings.EqualFold(*k, value) {
			return item, true
		}
	}
	var zero T
	return zero, false
}

func dataSourceLookupRead[T any](
	key, kind string,
	search func(client *apiclient.LookerSDK, value string) ([]T, error),
	getKey, getID func(T) *string,
	flatten func(T, *schema.ResourceData) error,
) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		value := d.Get(key).(string)

		items, err := search(m.(*apiclient.LookerSDK), value)
		if err != nil {
			return diag.FromErr(wrapSDKError(err, "Search", kind, "%s=%s", key, value))
		}

		item, ok := findExactMatch(items, value, getKey)
		if !ok {
			return diag.Errorf("no %s found with %s %q", kind, key, value)
		}

		d.SetId(*getID(item))
		return diag.FromErr(flatten(item, d))
	}
}

// Derived rather than hand-written so that attributes added to the resource later are exposed by the data source as well.
func dataSourceSchemaFromResourceSchema(rs map[string]*schema.Schema, key, keyDescription string) map[string]*schema.Schema {
	ds := make(map[string]*schema.Schema, len(rs))
	for k, v := range rs {
		ds[k] = dataSourceSchemaFromResourceSchemaElem(v)
	}
	ds[key] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: keyDescription + " Matched case-insensitively.",
	}
	return ds
}

func dataSourceSchemaFromResourceSchemaElem(v *schema.Schema) *schema.Schema {
	return &schema.Schema{
		Type:        v.Type,
		Description: v.Description,
		Sensitive:   v.Sensitive,
		Computed:    true,
		Set:         v.Set,
		// Nested blocks are copied as is (not made computed-only) and InternalValidate doesn't catch it.
		// Convert them here once a resource with nested blocks needs a data source.
		Elem: v.Elem,
	}
}

func hash(val interface{}) string {
	if val == nil || val.(string) == "" {
		return ""
	}
	sha := sha256.Sum256([]byte(val.(string)))
	return hex.EncodeToString(sha[:])
}
