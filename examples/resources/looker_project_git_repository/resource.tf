resource "random_id" "this" {
  byte_length = 2
}

resource "looker_project" "this" {
  name = "tf_test_${random_id.this.hex}"
}

resource "looker_project_git_deploy_key" "this" {
  project = looker_project.this.id
}

resource "looker_project_git_repository" "this" {
  project = looker_project.this.name
  git_service_name = "github"
  git_remote_url = "git@github.com:puc-business-intelligence/looker_project_tf_test_123.git"

  depends_on = [looker_project_git_deploy_key.this]
}

output "looker_project_git_repository" {
  value = looker_project_git_repository.this
}

