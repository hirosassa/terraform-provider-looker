---
page_title: "looker_folder Resource - terraform-provider-looker"
subcategory: ""
description: |-
  Manages a Looker Folder.
---

# looker_folder (Resource)

Manages a Looker Folder.

## Example Usage

```terraform
resource "looker_folder" "my_folder" {
  name      = "My Custom Folder"
  parent_id = "1"
}
```

## Schema

### Required

- `name` (String) The name of the folder.

### Optional

- `parent_id` (String) The ID of the parent folder. If not provided, it may default to a root-level entry depending on permissions.
