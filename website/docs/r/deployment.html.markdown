---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_deployment"
sidebar_current: "docs-bitbucket-resource-deployment"
description: |-
  Manage your pipelines repository deployment environments
---


# bitbucket\_deployment

This resource allows you to setup pipelines deployment environments.

# Example Usage

```hcl
resource "bitbucket_repository" "monorepo" {
    owner = "gob"
    name = "illusions"
    pipelines_enabled = true
}
resource "bitbucket_deployment" "test" {
  repository = bitbucket_repository.monorepo.id
  name = "test"
  stage = "Test"
}
```

# Argument Reference

* `name` - (Required) The name of the deployment environment
* `stage` - (Required) The stage (Test, Staging, Production)
* `repository` - (Required) The repository ID to which you want to assign this deployment environment to
* `uuid` - (Computed) The UUID of the deployment environment
