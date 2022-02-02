resource "looker_user_attribute_group_value" "my_user_attribute_group_value" {
  group_id          = 1
  user_attribute_id = 1
  value             = "foo"
}
