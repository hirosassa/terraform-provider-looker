resource "looker_folder" "test" {
  name      = "Test Folder"
  parent_id = "1"
}

resource "looker_group" "test" {
  name = "Test Group"
}

resource "looker_content_metadata_access" "test" {
  content_metadata_id = looker_folder.test.content_metadata_id
  group_id            = looker_group.test.id
  permission_type     = "view"
}
