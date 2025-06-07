resource "random_id" "this" {
  byte_length = 2
}

resource "looker_project" "this" {
  name = "tf_test_${random_id.this.hex}"
}
