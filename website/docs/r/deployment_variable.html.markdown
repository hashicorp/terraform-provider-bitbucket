---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_deployment_variable"
sidebar_current: "docs-bitbucket-resource-deployment-variable"
description: |-
  Manage your pipelines deployment variables
---


# bitbucket\_deployment\_variable

This resource allows you to setup pipelines deployment variables to manage your builds with. Once you have enabled pipelines on your repository you can then further setup deployment variables here to use.

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

resource "bitbucket_deployment_variable" "country" {
    key = "COUNTRY"
    value = "ke"
    deployment = bitbucket_deployment.monorepo.id
    secured = false
}
```

# Argument Reference

* `key` - (Required) The key of the key value pair
* `value` - (Required) The value of the key
* `repository` - (Required) The repository ID you want to put this variable onto.
* `secured` - (Optional) If you want to make this viewable in the UI.

* `uuid` - (Computed) The UUID of the variable
