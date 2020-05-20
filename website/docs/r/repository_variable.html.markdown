---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_repository_variable"
sidebar_current: "docs-bitbucket-resource-repository-variable"
description: |-
  Manage your pipelines repository variables and configuration
---


# bitbucket\_repository\_variable

This resource allows you to setup pipelines variables to manage your builds with. Once you have enabled pipelines on your repository you can then further setup variables here to use.

# Example Usage

```hcl
resource "bitbucket_repository" "monorepo" {
    owner = "gob"
    name = "illusions"
    pipelines_enabled = true
}

resource "bitbucket_repository_variable" "debug" {
    key = "DEBUG"
    value = "true"
    repository = "${bitbucket_repository.monorepo.id}"
    secured = false
}
```

# Argument Reference

* `key` - (Required) The key of the key value pair
* `value` - (Required) The value of the key
* `repository` - (Required) The repository ID you want to put this variable onto.
* `secured` - (Optional) If you want to make this viewable in the UI.

* `uuid` - (Computed) The UUID of the variable
