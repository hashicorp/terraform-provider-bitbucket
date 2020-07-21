---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_deployment_variable"
sidebar_current: "docs-bitbucket-resource-deployment-variable"
description: |-
  Manage variables for your pipelines deployment environments
---


# bitbucket\_deployment\_variable

This resource allows you to configure deployment variables.

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
  deployment = bitbucket_deployment.test.id
  key = "COUNTRY"
  value = "Kenya"
  secured = false
}
```

# Argument Reference

* `deployment` - (Required) The deployment ID you want to assign this variable to.
* `key` - (Required) The key of the variable
* `value` - (Required) The stage (Test, Staging, Production)
* `secured` - (Optional) Boolean indicating whether the variable contains sensitive data
* `uuid` - (Computed) The UUID of the variable
