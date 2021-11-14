resource "looker_user_attribute_user_value" "my_user_attribute_user_value" {
  user_id           = 1
  user_attribute_id = 1
  value             = "foo"
}
