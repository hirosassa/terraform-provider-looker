package looker

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
		return "", "", fmt.Errorf("Unexpected ID format (%q). Expected project:key", id)
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

func expandInt64ListFromSet(set interface{}) []int64 {
	var ints []int64
	for _, v := range set.(*schema.Set).List() {
		ints = append(ints, int64(v.(int)))
	}
	return ints
}
