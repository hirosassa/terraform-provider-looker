package looker

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	apiclient "github.com/looker-open-source/sdk-codegen/go/sdk/v4"
)

func resourceConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceConnectionCreate,
		Read:   resourceConnectionRead,
		Update: resourceConnectionUpdate,
		Delete: resourceConnectionDelete,
		Importer: &schema.ResourceImporter{
			State: resourceConnectionImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.EqualFold(old, new) // case-insensive comparing
				},
			},
			"host": {
				Type:     schema.TypeString,
				Required: true,
			},
			"port": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "443",
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
			},
			"certificate": {
				Type:      schema.TypeString,
				Optional:  true,
				ForceNew:  true,
				Sensitive: true,
			},
			"file_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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
			},
			"max_billing_gigabyte": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ssl": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"verify_ssl": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
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
			},
			"sql_writing_with_info_schema": {
				Type:     schema.TypeBool,
				Optional: true,
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
			},
			"disable_context_comment": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"oauth_application_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceConnectionCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.LookerSDK)

	body := expandWriteDBConnection(d)

	result, err := client.CreateConnection(*body, nil)
	if err != nil {
		return err
	}

	d.SetId(*result.Name)

	return resourceConnectionRead(d, m)
}

func resourceConnectionRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.LookerSDK)
	connectionName := d.Get("name").(string)

	connection, err := client.Connection(connectionName, "", nil)
	if err != nil {
		if strings.Contains(err.Error(), "Not found") {
			d.SetId("")
			return nil
		}
		return err
	}

	return flattenConnection(connection, d)
}

func resourceConnectionUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.LookerSDK)

	name := d.Get("name").(string)
	body := expandWriteDBConnection(d)

	_, err := client.UpdateConnection(name, *body, nil)
	if err != nil {
		return err
	}

	return resourceConnectionRead(d, m)
}

func resourceConnectionDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.LookerSDK)

	connectionName := d.Id()

	_, err := client.DeleteConnection(connectionName, nil)
	if err != nil {
		return err
	}

	return nil
}

func resourceConnectionImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := resourceConnectionRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

func expandWriteDBConnection(d *schema.ResourceData) *apiclient.WriteDBConnection {
	name := d.Get("name").(string)
	host := d.Get("host").(string)
	port := d.Get("port").(int64)
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	certificate := d.Get("certificate").(string)
	fileType := d.Get("file_type").(string)
	database := d.Get("database").(string)
	dbTimezone := d.Get("db_timezone").(string)
	queryTimezone := d.Get("query_timezone").(string)
	schema := d.Get("schema").(string)
	maxConnections := d.Get("max_connctions").(int64)
	maxBillingGigabytes := d.Get("max_billing_gigabytes").(string)
	ssl := d.Get("ssl").(bool)
	verifySsl := d.Get("verify_ssl").(bool)
	tmpDbName := d.Get("tmp_db_name").(string)
	jdbcAdditionalParams := d.Get("jdbc_additional_params").(string)
	poolTimeout := d.Get("pool_timeout").(int64)
	dialectName := d.Get("dialect_name").(string)
	userDbCredentials := d.Get("user_db_credentials").(bool)
	maintenanceCron := d.Get("maintenance_cron").(string)
	sqlRunnerPrecacheTables := d.Get("sql_runner_precache_tables").(bool)
	sqlWritingWithInfoSchema := d.Get("sql_writing_with_info_schema").(bool)
	afterConnectStatements := d.Get("after_connect_statements").(string)
	tunnelId := d.Get("tunnel_id").(string)
	pdtConcurrency := d.Get("pdt_concurrency").(int64)
	disable_context_comment := d.Get("disable_context_comment").(bool)
	oauthApplicationId := d.Get("oauth_application_id").(int64)

	writeDBConnection := &apiclient.WriteDBConnection{
		Name:                     &name,
		Host:                     &host,
		Port:                     &port,
		Username:                 &username,
		Password:                 &password,
		Certificate:              &certificate,
		FileType:                 &fileType,
		Database:                 &database,
		DbTimezone:               &dbTimezone,
		QueryTimezone:            &queryTimezone,
		Schema:                   &schema,
		MaxConnections:           &maxConnections,
		MaxBillingGigabytes:      &maxBillingGigabytes,
		Ssl:                      &ssl,
		VerifySsl:                &verifySsl,
		TmpDbName:                &tmpDbName,
		JdbcAdditionalParams:     &jdbcAdditionalParams,
		PoolTimeout:              &poolTimeout,
		DialectName:              &dialectName,
		UserDbCredentials:        &userDbCredentials,
		UserAttributeFields:      nil,
		MaintenanceCron:          &maintenanceCron,
		SqlRunnerPrecacheTables:  &sqlRunnerPrecacheTables,
		SqlWritingWithInfoSchema: &sqlWritingWithInfoSchema,
		AfterConnectStatements:   &afterConnectStatements,
		PdtContextOverride:       nil,
		TunnelId:                 &tunnelId,
		PdtConcurrency:           &pdtConcurrency,
		DisableContextComment:    &disable_context_comment,
		OauthApplicationId:       &oauthApplicationId,
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

	return writeDBConnection
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
	if err := d.Set("certificate", connection.Certificate); err != nil {
		return err
	}
	if err := d.Set("file_type", connection.FileType); err != nil {
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
	if err := d.Set("user_attribute_fields", flattenStringListToSet(*connection.UserAttributeFields)); err != nil {
		return err
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
		if err := d.Set("pdt_context_override", []map[string]interface{}{
			{
				"context":                  *connection.PdtContextOverride.Context,
				"host":                     *connection.PdtContextOverride.Host,
				"port":                     *connection.PdtContextOverride.Port,
				"username":                 *connection.PdtContextOverride.Username,
				"password":                 *connection.PdtContextOverride.Password,
				"certitficate":             *connection.PdtContextOverride.Certificate,
				"file_type":                *connection.PdtContextOverride.FileType,
				"database":                 *connection.PdtContextOverride.Database,
				"schema":                   *connection.PdtContextOverride.Schema,
				"jdbc_additional_params":   *connection.PdtContextOverride.JdbcAdditionalParams,
				"after_connect_statements": *connection.PdtContextOverride.AfterConnectStatements,
			},
		}); err != nil {
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

func expandStringListFromSet(set *schema.Set) []string {
	strings := make([]string, 0, set.Len())
	for _, v := range set.List() {
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
