---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_repository_deploy_key"
sidebar_current: "docs-bitbucket-resource-deploy-key"
description: |-
  Provides a Bitbucket repository deploy key
---

# bitbucket\_repository\_deploy\_key

Provides a Bitbucket repository deploy key.

Deployment keys allow users to clone/pull source from a repository over SSH using Git and Hg. Deployment keys are similar to adding SSH keys to your account, but they are done on a per-repository basis.

## Example Usage

```hcl
resource "bitbucket_repository_deploy_key" "example" {
  owner = "myteam"
  repository = "terraform-code"
  public_key_contents = "${trimspace(file("my_bitbucket_deploy_key.pub"))}"
  label = "Terraform deploy key"
}
```

## Argument Reference

The following arguments are supported:

* `owner` - (Required) The owner of this repository. Can be you or any team you
  have write access to.
* `repository` - (Required) The name of the repository.
* `public_key_contents` - (Required) String beginning with "ssh-rsa AAAA..."
* `label` - (Optional) Name for this key in the Bitbucket user interface
