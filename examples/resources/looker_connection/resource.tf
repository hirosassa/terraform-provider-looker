resource "looker_connection" "bigquery_connection" {
  name         = "bigquery_connection"
  host         = "gcp_project_id"
  user         = var.gcp_service_account_email
  certificate  = var.gcp_service_account_json
  file_type    = ".json"
  database     = "dataset_name"
  tmp_db_name  = "tmp_dataset_name"
  dialect_name = "bigquery_standard_sql"
}

resource "looker_connection" "snowflake_connection" {
  name                   = "snowflake_connection"
  host                   = var.snowflake_host
  port                   = 443
  user                   = var.snowflake_username
  password               = var.snowflake_password
  database               = "DATABASE"
  db_timezone            = "UTC"
  query_timezone         = "UTC"
  schema                 = "SAMPLE"
  ssl                    = true
  tmp_db_name            = "tmp_dataset_name"
  jdbc_additional_params = "account=${var.snowflake_account}&warehouse=WHARE_HOUSE"
  dialect_name           = "snowflake"
}
