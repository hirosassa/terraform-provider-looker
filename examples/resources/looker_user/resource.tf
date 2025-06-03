resource "looker_user" "user" {
  email                     = "user@email.com"
  first_name                = "Reporting"
  last_name                 = "API User"
  send_setup_link_on_create = true
}
