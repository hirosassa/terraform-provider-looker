resource "random_id" "this" {
  byte_length = 2
}

resource "looker_project" "this" {
  name = "tf_test_${random_id.this.hex}"
}

resource "looker_project_git_deploy_key" "this" {
  project = looker_project.this.id
}

output "looker_project" {
  value = looker_project.this
}

output "project_git_deploy_key" {
  value = looker_project_git_deploy_key.this
}

