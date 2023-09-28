package looker

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	apiclient "github.com/looker-open-source/sdk-codegen/go/sdk/v4"
)

func resourceConnection() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceConnectionCreate,
		ReadContext:   resourceConnectionRead,
		UpdateContext: resourceConnectionUpdate,
		DeleteContext: resourceConnectionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceConnectionImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.EqualFold(old, new) // case-insensive comparing
				},
				ValidateFunc: validation.StringDoesNotContainAny(" "),
			},
			"host": {
				Type:     schema.TypeString,
				Required: true,
			},
			"port": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				ForceNew:  true,
				Sensitive: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Id() != ""
				},
			},
			"certificate": {
				Type:        schema.TypeString,
				Description: "Base64 encoded certificate body for server authentication (when " +
										 "appropriate for the dialect). Due to limitations in the Looker " +
										 "API, changes made outside of Terraform cannot be detected.",
				Optional:    true,
				Sensitive:   true,
				StateFunc:   hash,
			},
			"file_type": {
				Type:         schema.TypeString,
				Description:  "Certificate key file type (.json or .p12).",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{".json", ".p12"}, false),
			},
			"database": {
				Type:     schema.TypeString,
				Required: true,
			},
			"db_timezone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"query_timezone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"schema": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"max_connections": {
				Type:     schema.TypeInt,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Id() != ""
				},
			},
			"max_billing_gigabytes": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ssl": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"verify_ssl": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"tmp_db_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"jdbc_additional_params": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"pool_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"dialect_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"user_db_credentials": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"user_attribute_fields": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"maintenance_cron": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sql_runner_precache_tables": {
				Type:     schema.TypeBool,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Id() != ""
				},
			},
			"sql_writing_with_info_schema": {
				Type:     schema.TypeBool,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Id() != ""
				},
			},
			"after_connect_statements": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"pdt_context_override": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"context": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"pdt"}, false),
						},
						"host": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"port": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"username": {
							Type:     schema.TypeString,
							ForceNew: true,
							Optional: true,
						},
						"passowrd": {
							Type:      schema.TypeString,
							ForceNew:  true,
							Optional:  true,
							Sensitive: true,
						},
						"certificate": {
							Type:      schema.TypeString,
							ForceNew:  true,
							Optional:  true,
							Sensitive: true,
						},
						"file_type": {
							Type:     schema.TypeString,
							ForceNew: true,
							Optional: true,
						},
						"database": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"schema": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"jdbc_additional_params": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"after_connect_statements": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"tunnel_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"pdt_concurrency": {
				Type:     schema.TypeInt,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Id() != ""
				},
			},
			"disable_context_comment": {
				Type:     schema.TypeBool,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Id() != ""
				},
			},
			"oauth_application_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceConnectionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	body, err := expandWriteDBConnection(d)
	if err != nil {
		return diag.FromErr(err)
	}

	result, err := client.CreateConnection(*body, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*result.Name)

	return resourceConnectionRead(ctx, d, m)
}

func resourceConnectionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)
	connectionName := d.Id()

	connection, err := client.Connection(connectionName, "", nil)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	return diag.FromErr(flattenConnection(connection, d))
}

func resourceConnectionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	name := d.Id()
	body, err := expandWriteDBConnection(d)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.UpdateConnection(name, *body, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceConnectionRead(ctx, d, m)
}

func resourceConnectionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	connectionName := d.Id()

	_, err := client.DeleteConnection(connectionName, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceConnectionImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := resourceConnectionRead(ctx, d, m); err != nil {
		return nil, fmt.Errorf("failed to read connection: %v", err)
	}
	return []*schema.ResourceData{d}, nil
}

func expandWriteDBConnection(d *schema.ResourceData) (*apiclient.WriteDBConnection, error) {
	// required values
	name := d.Get("name").(string)
	host := d.Get("host").(string)
	username := d.Get("username").(string)
	database := d.Get("database").(string)
	dialectName := d.Get("dialect_name").(string)
	writeDBConnection := &apiclient.WriteDBConnection{
		Name:        &name,
		Host:        &host,
		Username:    &username,
		Database:    &database,
		DialectName: &dialectName,
	}

	// optional values
	if v, ok := d.GetOk("port"); ok {
		port := v.(string) // for api breaking change
		writeDBConnection.Port = &port
	}
	if v, ok := d.GetOk("password"); ok {
		password := v.(string)
		writeDBConnection.Password = &password
	}
	if v, ok := d.GetOk("certificate"); ok {
		certificate := v.(string)
		writeDBConnection.Certificate = &certificate
	}
	if v, ok := d.GetOk("file_type"); ok {
		fileType := v.(string)
		writeDBConnection.FileType = &fileType
	}
	if v, ok := d.GetOk("db_timezone"); ok {
		dbTimezone := v.(string)
		writeDBConnection.DbTimezone = &dbTimezone
	}
	if v, ok := d.GetOk("query_timezone"); ok {
		queryTimezone := v.(string)
		writeDBConnection.QueryTimezone = &queryTimezone
	}
	if v, ok := d.GetOk("schema"); ok {
		schema := v.(string)
		writeDBConnection.Schema = &schema
	}
	if v, ok := d.GetOk("max_connections"); ok {
		maxConnections := int64(v.(int))
		writeDBConnection.MaxConnections = &maxConnections
	}
	if v, ok := d.GetOk("max_billing_gigabytes"); ok {
		maxBillingGigabytes := v.(string)
		writeDBConnection.MaxBillingGigabytes = &maxBillingGigabytes
	}
	if v, ok := d.GetOk("ssl"); ok {
		ssl := v.(bool)
		writeDBConnection.Ssl = &ssl
	}
	if v, ok := d.GetOk("verify_ssl"); ok {
		verifySsl := v.(bool)
		writeDBConnection.VerifySsl = &verifySsl
	}
	if v, ok := d.GetOk("tmp_db_name"); ok {
		tmpDbName := v.(string)
		writeDBConnection.TmpDbName = &tmpDbName
	}
	if v, ok := d.GetOk("jdbc_additional_params"); ok {
		jdbcAdditionalParams := v.(string)
		writeDBConnection.JdbcAdditionalParams = &jdbcAdditionalParams
	}

	if v, ok := d.GetOk("pool_timeout"); ok {
		poolTimeout := int64(v.(int))
		writeDBConnection.PoolTimeout = &poolTimeout
	}
	if v, ok := d.GetOk("user_db_credentials"); ok {
		userDbCredentials := v.(bool)
		writeDBConnection.UserDbCredentials = &userDbCredentials
	}
	if v, ok := d.GetOk("maintenance_cron"); ok {
		maintenanceCron := v.(string)
		writeDBConnection.MaintenanceCron = &maintenanceCron
	}
	if v, ok := d.GetOk("sql_runner_precache_tables"); ok {
		sqlRunnerPrecacheTables := v.(bool)
		writeDBConnection.SqlRunnerPrecacheTables = &sqlRunnerPrecacheTables
	}
	if v, ok := d.GetOk("sql_writing_with_info_schema"); ok {
		sqlWritingWithInfoSchema := v.(bool)
		writeDBConnection.SqlWritingWithInfoSchema = &sqlWritingWithInfoSchema
	}
	if v, ok := d.GetOk("after_connect_statements"); ok {
		afterConnectStatements := v.(string)
		writeDBConnection.AfterConnectStatements = &afterConnectStatements
	}
	if v, ok := d.GetOk("tunnel_id"); ok {
		tunnelId := v.(string)
		writeDBConnection.TunnelId = &tunnelId
	}
	if v, ok := d.GetOk("pdt_concurrency"); ok {
		pdtConcurrency := int64(v.(int))
		writeDBConnection.PdtConcurrency = &pdtConcurrency
	}
	if v, ok := d.GetOk("disable_context_comment"); ok {
		disable_context_comment := v.(bool)
		writeDBConnection.DisableContextComment = &disable_context_comment
	}
	if v, ok := d.GetOk("oauth_application_id"); ok {
		oauthApplicationId := v.(string) // for api breaking change
		writeDBConnection.OauthApplicationId = &oauthApplicationId
	}

	userAttributeFields := expandStringListFromSet(d.Get("user_attribute_fields").(*schema.Set))
	writeDBConnection.UserAttributeFields = &userAttributeFields

	if _, ok := d.GetOk("pdt_context_override"); ok {
		var pdtContextOverride apiclient.WriteDBConnectionOverride
		if v, ok := d.GetOk("pdt_context_override.0.context"); ok {
			pdtContextOverride.Context = v.(*string)
		}
		if v, ok := d.GetOk("pdt_context_override.0.host"); ok {
			pdtContextOverride.Host = v.(*string)
		}
		if v, ok := d.GetOk("pdt_context_override.0.port"); ok {
			pdtContextOverride.Port = v.(*string)
		}
		if v, ok := d.GetOk("pdt_context_override.0.username"); ok {
			pdtContextOverride.Username = v.(*string)
		}
		if v, ok := d.GetOk("pdt_context_override.0.password"); ok {
			pdtContextOverride.Password = v.(*string)
		}
		if v, ok := d.GetOk("pdt_context_override.0.certificate"); ok {
			pdtContextOverride.Certificate = v.(*string)
		}
		if v, ok := d.GetOk("pdt_context_override.0.file_type"); ok {
			pdtContextOverride.FileType = v.(*string)
		}
		if v, ok := d.GetOk("pdt_context_override.0.database"); ok {
			pdtContextOverride.Database = v.(*string)
		}
		if v, ok := d.GetOk("pdt_context_override.0.schema"); ok {
			pdtContextOverride.Schema = v.(*string)
		}
		if v, ok := d.GetOk("pdt_context_override.0.jdbc_additional_params"); ok {
			pdtContextOverride.JdbcAdditionalParams = v.(*string)
		}
		if v, ok := d.GetOk("pdt_context_override.0.after_connect_statements"); ok {
			pdtContextOverride.AfterConnectStatements = v.(*string)
		}

		writeDBConnection.PdtContextOverride = &pdtContextOverride
	}

	return writeDBConnection, nil
}

func flattenConnection(connection apiclient.DBConnection, d *schema.ResourceData) error {
	if err := d.Set("name", connection.Name); err != nil {
		return err
	}
	if err := d.Set("host", connection.Host); err != nil {
		return err
	}
	if err := d.Set("port", connection.Port); err != nil {
		return err
	}
	if err := d.Set("username", connection.Username); err != nil {
		return err
	}
	if err := d.Set("password", connection.Password); err != nil {
		return err
	}
	if err := d.Set("database", connection.Database); err != nil {
		return err
	}
	if err := d.Set("db_timezone", connection.DbTimezone); err != nil {
		return err
	}
	if err := d.Set("query_timezone", connection.QueryTimezone); err != nil {
		return err
	}
	if err := d.Set("schema", connection.Schema); err != nil {
		return err
	}
	if err := d.Set("max_connections", connection.MaxConnections); err != nil {
		return err
	}
	if err := d.Set("max_billing_gigabytes", connection.MaxBillingGigabytes); err != nil {
		return err
	}
	if err := d.Set("ssl", connection.Ssl); err != nil {
		return err
	}
	if err := d.Set("verify_ssl", connection.VerifySsl); err != nil {
		return err
	}
	if err := d.Set("tmp_db_name", connection.TmpDbName); err != nil {
		return err
	}
	if err := d.Set("jdbc_additional_params", connection.JdbcAdditionalParams); err != nil {
		return err
	}
	if err := d.Set("pool_timeout", connection.PoolTimeout); err != nil {
		return err
	}
	if err := d.Set("dialect_name", connection.DialectName); err != nil {
		return err
	}
	if connection.UserAttributeFields != nil {
		if err := d.Set("user_attribute_fields", flattenStringListToSet(*connection.UserAttributeFields)); err != nil {
			return err
		}
	}
	if err := d.Set("maintenance_cron", connection.MaintenanceCron); err != nil {
		return err
	}
	if err := d.Set("sql_runner_precache_tables", connection.SqlRunnerPrecacheTables); err != nil {
		return err
	}
	if err := d.Set("sql_writing_with_info_schema", connection.SqlWritingWithInfoSchema); err != nil {
		return err
	}
	if err := d.Set("after_connect_statements", connection.AfterConnectStatements); err != nil {
		return err
	}

	if connection.PdtContextOverride != nil {
		pdtContextOverride := make(map[string]interface{})

		if connection.PdtContextOverride.Context != nil {
			pdtContextOverride["context"] = *connection.PdtContextOverride.Context
		}
		if connection.PdtContextOverride.Host != nil {
			pdtContextOverride["host"] = *connection.PdtContextOverride.Host
		}
		if connection.PdtContextOverride.Port != nil {
			pdtContextOverride["port"] = *connection.PdtContextOverride.Port
		}
		if connection.PdtContextOverride.Username != nil {
			pdtContextOverride["username"] = *connection.PdtContextOverride.Username
		}
		if connection.PdtContextOverride.Password != nil {
			pdtContextOverride["password"] = *connection.PdtContextOverride.Password
		}
		if connection.PdtContextOverride.Certificate != nil {
			pdtContextOverride["certificate"] = *connection.PdtContextOverride.Certificate
		}
		if connection.PdtContextOverride.FileType != nil {
			pdtContextOverride["file_type"] = *connection.PdtContextOverride.FileType
		}
		if connection.PdtContextOverride.Database != nil {
			pdtContextOverride["database"] = *connection.PdtContextOverride.Database
		}
		if connection.PdtContextOverride.Schema != nil {
			pdtContextOverride["schema"] = *connection.PdtContextOverride.Schema
		}

		if err := d.Set("pdt_context_override", []map[string]interface{}{pdtContextOverride}); err != nil {
			return err
		}
	}
	if err := d.Set("tunnel_id", connection.TunnelId); err != nil {
		return err
	}
	if err := d.Set("pdt_concurrency", connection.PdtConcurrency); err != nil {
		return err
	}
	if err := d.Set("disable_context_comment", connection.DisableContextComment); err != nil {
		return err
	}
	if err := d.Set("oauth_application_id", connection.OauthApplicationId); err != nil {
		return err
	}
	return nil
}
