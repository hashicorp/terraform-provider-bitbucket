---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_user_ssh_key"
sidebar_current: "docs-bitbucket-resource-user-ssh-key"
description: |-
  Adds an SSH key to a Bitbucket use
---

# bitbucket\_user\_ssh\_key

Adds an SSH public key to a Bitbucket user.

## Example Usage

```hcl-terraform
resource "tls_private_key" "key1" {
  algorithm = "rsa"
}

resource "bitbucket_user_ssh_key" "key1" {
  key        = tls_private_key.key1.public_key_openssh
  owner      = "myuser"
  label      = "key-1"
}
```

## Argument Reference

The following arguments are supported:

* `key` - (Required) The key in OpenSSH `authorized_keys` format.
* `owner` - (Required) The owner of the key. 
            This can either be the username or the UUID of the account,
            surrounded by curly-braces, for example: {account UUID}.
* `label` - (Optional) The label for the key to be shown in Bitbucket.
