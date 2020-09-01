---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_project"
sidebar_current: "docs-bitbucket-resource-project"
description: |-
  Create and manage a Bitbucket project
---


# bitbucket\_project

This resource allows you to manage your projects in your bitbucket team. 

# Example Usage

```hcl
# Manage your repository
resource "bitbucket_project" "devops" {
  owner = "my-team"
  name  = "devops"
  key = "DEVOPS"
}
```

## Argument Reference

The following arguments are supported:

* `owner` - (Required) The owner of this project. Can be you or any team you have write access to.
* `name` - (Required) The name of the project
* `key` - (Required) The key used for this project
* `description` - (Optional) The description of the project
* `is_private` - (Optional) If you want to keep the project private - defaults to true

## Import

Projects can be imported using their `owner/name` ID, e.g.

```
$ terraform import bitbucket_project.my-repo my-account/my-project
```
