resource "looker_lookml_model" "lookml_model" {
  name                        = "LookML Model"
  allowed_db_connection_names = ["bigquery-connection"]
  project_name                = "lookml_model"
}
