package looker

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

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

// func expandInt64ListFromSet(set interface{}) []int64 {
// 	var ints []int64
// 	for _, v := range set.(*schema.Set).List() {
// 		ints = append(ints, int64(v.(int)))
// 	}
// 	return ints
// }

func updateApiSessionWorkspaceId(client *apiclient.LookerSDK, workspaceId string) error {
	writeApiSession := apiclient.WriteApiSession{
		WorkspaceId: &workspaceId,
	}
	_, err := client.UpdateSession(writeApiSession, nil)
	return err
}

func hash(val interface{}) string {
	if val == nil || val.(string) == "" {
		return ""
	}
	sha := sha256.Sum256([]byte(val.(string)))
	return hex.EncodeToString(sha[:])
}

