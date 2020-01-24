---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_repository_deploy_key"
sidebar_current: "docs-bitbucket-resource-repository-deploy-key"
description: |-
  Adds a deploy key to a Bitbucket repository
---

# bitbucket\_repository\_deploy\_key

Adds a deploy key to a Bitbucket repository.

## Example Usage

```hcl-terraform
resource "tls_private_key" "key1" {
  algorithm = "rsa"
}

resource "bitbucket_repository_deploy_key" "key1" {
  key       = tls_private_key.key1.public_key_openssh
  owner     = "myteam"
  repo_slug = "my-repo"
  label     = "key-1"
}
```

## Argument Reference

The following arguments are supported:

* `key` - (Required) The key in OpenSSH `authorized_keys` format.
* `owner` - (Required) The owner of the repository. Can be a user or a team.
* `repo_slug` - (Required) The repository slug.
* `label` - (Optional) The label for the key to be shown in Bitbucket.

## Computed Arguments

None.

## Import

Repository deploy keys cannot be imported.
