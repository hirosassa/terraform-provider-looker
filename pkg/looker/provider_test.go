package looker

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	testAccProviders map[string]*schema.Provider
	testAccProvider  *schema.Provider
)

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"looker": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	if err := os.Getenv("LOOKER_API_BASE_URL"); err == "" {
		t.Fatal("LOOKER_API_BASE_URL must be set for acceptance tests")
	}
	if err := os.Getenv("LOOKER_API_CLIENT_ID"); err == "" {
		t.Fatal("LOOKER_API_CLIENT_ID must be set for acceptance tests")
	}
	if err := os.Getenv("LOOKER_API_CLIENT_SECRET"); err == "" {
		t.Fatal("LOOKER_API_CLIENT_SECRET must be set for acceptance tests")
	}
}
