resource "looker_group_membership" "group_membership" {
  group_id = looker_group.group.id
  user_id  = looker_user.user.id
}
