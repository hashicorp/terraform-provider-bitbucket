---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_deploy_key"
sidebar_current: "docs-bitbucket-resource-deploy-key"
description: |-
  Provides a Bitbucket Webhook
---

# bitbucket\_deploy\_key

Provides a Bitbucket deploy key.

This allows you read only access to your repo.

## Example Usage

```hcl
# Manage your repositories deploy keys
resource "bitbucket_deploy_key" "sample" {
  owner       = "myteam"
  repository  = "terraform-code"
  label       = "CI/CD key"
  key         = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDEZBW8HOt+CKpEoFtF5q2NJNUyb+Z3wFcmZBX1UtW6CROJDz8AZfQTWQCi/7pz5+K1iqkZ7VvEg153MxbJXXa2sbpzeqLTuZdk8dGumhGxOGua6oLWLqO51k3H/dgK/tF4IQJqTe8p7XaolL4dnz87MU9GdDL1JV+ctdWH96+lX+9XGyC3momWNGCUxtGWwJAeyU0PSwcNmUjqqAryKMCrPtajKRjcjKS2WMpG1RML9nlkV4JLljof4wDo9aDxMhYSMyV1FQryUMcrOBaVbmP8AKru2AipHY89gReRG3pLgJrCe4Fi+d+BTqmMoJ2Sa8+RPPZA72sKg91+0KigIsl7 test@test"
}
```

## Argument Reference

The following arguments are supported:

* `owner` - (Required) The owner of this repository. Can be you or any team you
  have write access to.
* `repository` - (Required) The name of the repository.
* `label` - (Required) A description of the key
* `key` - (Required) A base64 encoded public RSA key (at least 2048 bits in length)
