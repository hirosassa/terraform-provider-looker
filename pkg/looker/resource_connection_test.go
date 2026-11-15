package looker

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	apiclient "github.com/looker-open-source/sdk-codegen/go/sdk/v4"
	"github.com/stretchr/testify/assert"
)

func TestAcc_Connection(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: connectionConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConnectionExists("looker_connection.test"),
					resource.TestCheckResourceAttr("looker_connection.test", "name", strings.ToLower(name)),
					resource.TestCheckResourceAttr("looker_connection.test", "host", "test_project"),
				),
			},
			{
				ResourceName:      "looker_connection.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
		CheckDestroy: testAccCheckConnectionDestroy,
	})
}

func testAccCheckConnectionExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("connection setting not found from resources: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no connection setting ID is set")
		}

		client := testAccProvider.Meta().(*apiclient.LookerSDK)
		connectionName := rs.Primary.ID

		_, err := client.Connection(connectionName, "", nil)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckConnectionDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*apiclient.LookerSDK)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "looker_connection" {
			continue
		}

		connectionName := rs.Primary.ID
		_, err := client.Connection(connectionName, "", nil)
		if err != nil {
			if strings.Contains(err.Error(), "404") {
				return nil // successfully destroyed
			}
			return err
		}

		return fmt.Errorf("connection '%s' still exists", connectionName)
	}

	return nil
}

func connectionConfig(name string) string {
	return fmt.Sprintf(`
	resource "looker_connection" "test" {
		name         = "%s"
		host         = "test_project"
		username     = "test@testproject.iam.gserviceaccount.com"
		certificate  = filebase64("testdata/gcp-sa.json")
		file_type    = ".json"
		database     = "test_dataset"
		tmp_db_name  = "tmp_test_dataset"
		dialect_name = "bigquery_standard_sql"
	}
	`, name)
}

func TestAcc_ConnectionSnowflakeKeyPair(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: connectionConfigSnowflakeKeyPair(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConnectionExists("looker_connection.test"),
					resource.TestCheckResourceAttr("looker_connection.test", "uses_key_pair_auth", "true"),
				),
			},
			{
				ResourceName:            "looker_connection.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"certificate", "file_type", "password"},
			},
		},
		CheckDestroy: testAccCheckConnectionDestroy,
	})
}

func connectionConfigSnowflakeKeyPair(name string) string {
	return fmt.Sprintf(`
	resource "looker_connection" "test" {
		name               = "%s"
		host               = "example.snowflakecomputing.com"
		port               = "443"
		username           = "TEST_USER"
		certificate        = base64encode("dummy")
		file_type          = ".p8"
		uses_key_pair_auth = true
		database           = "TEST_DATABASE"
		dialect_name       = "snowflake"
	}
	`, name)
}

func TestExpandWriteDBConnection_KeyPairAuth(t *testing.T) {
	base := map[string]interface{}{
		"name":         "test",
		"host":         "example.snowflakecomputing.com",
		"username":     "TEST_USER",
		"database":     "TEST_DATABASE",
		"dialect_name": "snowflake",
	}
	usesKeyPairAuth := true
	fileType := ".p8"

	tests := map[string]struct {
		raw                 map[string]interface{}
		wantUsesKeyPairAuth *bool
		wantFileType        *string
	}{
		"key pair configured": {
			raw: map[string]interface{}{
				"certificate":        "ZHVtbXk=",
				"file_type":          ".p8",
				"uses_key_pair_auth": true,
			},
			wantUsesKeyPairAuth: &usesKeyPairAuth,
			wantFileType:        &fileType,
		},
		"key pair not configured": {
			raw:                 map[string]interface{}{},
			wantUsesKeyPairAuth: nil,
			wantFileType:        nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			raw := map[string]interface{}{}
			for k, v := range base {
				raw[k] = v
			}
			for k, v := range tt.raw {
				raw[k] = v
			}
			d := schema.TestResourceDataRaw(t, resourceConnection().Schema, raw)

			got := expandWriteDBConnection(d)

			assert.Equal(t, tt.wantUsesKeyPairAuth, got.UsesKeyPairAuth)
			assert.Equal(t, tt.wantFileType, got.FileType)
		})
	}
}
