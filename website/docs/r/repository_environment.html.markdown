---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_repository_environment"
sidebar_current: "docs-bitbucket-resource-repository-environment"
description: |-
  Manage your pipelines repository environments
---


# bitbucket\_repository\_environment

This resource allows you to setup pipelines environments.

# Example Usage

```hcl
resource "bitbucket_repository" "monorepo" {
    owner = "gob"
    name = "illusions"
    pipelines_enabled = true
}

resource "bitbucket_environment" "test" {
  repository = bitbucket_repository.monorepo.id
  name = "test"
  stage = "Test"
}
```

# Argument Reference

* `name` - (Required) The name of the environment
* `stage` - (Required) The stage (Test, Staging, Production)
* `repository` - (Required) The repository ID you want to put this variable onto.
* `uuid` - (Computed) The UUID of the environment
