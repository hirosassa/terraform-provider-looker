package looker

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apiclient "github.com/looker-open-source/sdk-codegen/go/sdk/v4"
)

func resourceConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceConnectionCreate,
		Read: resourceConnectionRead,
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
					return strings.EqualFold(old, new)  // case-insensive comparing
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
			"user_name": {
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
			"ssl": {
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
			"dialect_name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceConnectionCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.LookerSDK)

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
	ssl := d.Get("ssl").(bool)
	tmpDbName := d.Get("tmp_db_name").(string)
	jdbcAdditionalParams := d.Get("jdbc_additional_params").(string)
	dialectName := d.Get("dialect_name").(string)

	body := apiclient.WriteDBConnection{
		Name:                 &name,
		Host:                 &host,
		Port:                 &port,
		Username:             &username,
		Password:             &password,
		Certificate:          &certificate,
		FileType:             &fileType,
		Database:             &database,
		DbTimezone:           &dbTimezone,
		QueryTimezone:        &queryTimezone,
		Schema:               &schema,
		Ssl:                  &ssl,
		TmpDbName:            &tmpDbName,
		JdbcAdditionalParams: &jdbcAdditionalParams,
		DialectName:          &dialectName,
	}

	result, err := client.CreateConnection(body, nil)
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

	if err = d.Set("name", connection.Name); err != nil {
		return err
	}
	if err = d.Set("host", connection.Host); err != nil {
		return err
	}
	if err = d.Set("port", connection.Port); err != nil {
		return err
	}
	if err = d.Set("username", connection.Username); err != nil {
		return err
	}
	if err = d.Set("password", connection.Password); err != nil {
		return err
	}
	if err = d.Set("certificate", connection.Certificate); err != nil {
		return err
	}
	if err = d.Set("file_type", connection.FileType); err != nil {
		return err
	}
	if err = d.Set("database", connection.Database); err != nil {
		return err
	}
	if err = d.Set("db_timezone", connection.DbTimezone); err != nil {
		return err
	}
	if err = d.Set("query_timezone", connection.QueryTimezone); err != nil {
		return err
	}
	if err = d.Set("schema", connection.Schema); err != nil {
		return err
	}
	if err = d.Set("ssl", connection.Ssl); err != nil {
		return err
	}
	if err = d.Set("tmp_db_name", connection.TmpDbName); err != nil {
		return err
	}
	if err = d.Set("jdbc_additional_params", connection.JdbcAdditionalParams); err != nil {
		return err
	}
	if err = d.Set("dialect_name", connection.DialectName); err != nil {
		return err
	}

	return nil
}

func resourceConnectionUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.LookerSDK)

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
	ssl := d.Get("ssl").(bool)
	tmpDbName := d.Get("tmp_db_name").(string)
	jdbcAdditionalParams := d.Get("jdbc_additional_params").(string)
	dialectName := d.Get("dialect_name").(string)

	body := apiclient.WriteDBConnection{
		Name:                 &name,
		Host:                 &host,
		Port:                 &port,
		Username:             &username,
		Password:             &password,
		Certificate:          &certificate,
		FileType:             &fileType,
		Database:             &database,
		DbTimezone:           &dbTimezone,
		QueryTimezone:        &queryTimezone,
		Schema:               &schema,
		Ssl:                  &ssl,
		TmpDbName:            &tmpDbName,
		JdbcAdditionalParams: &jdbcAdditionalParams,
		DialectName:          &dialectName,
	}

	_, err := client.UpdateConnection(name, body, nil)
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
