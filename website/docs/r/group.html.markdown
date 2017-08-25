---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_group"
sidebar_current: "docs-bitbucket-resource-group"
description: |-
  Provides a Bitbucket Group
---

# bitbucket\_group

Provides a Bitbucket group resource.

This resource allows you manage your groups' settings, such as their name,
if they are automatically added to new repositories, their permissions,
and their members.

## Example Usage

```hcl
# Manage your group
resource "bitbucket_group" "devs" {
  accountname = "my-account"
  name        = "Developers group"
  auto_add    = false
  permission  = "write"

  members = [
    "dev1",
    "dev2",
  ]
}
```


* `accountname` - (Required) The team or individual account name. You can supply an account name or the primary email address for the account.
* `name` - (Required) The name of the group.
* `auto_add` - (Optional) True to automatically add the group to new repositories.
* `permissions` - (Optional) One of `read`, `write`, or `admin`.
* `members` (Required) - The list of the group's members. Note that it can be empty.

## Computed Arguments

The following arguments are computed:

 * `slug` - The slug of the group.
 * `owner` - The name of the group's owner.

## Import

Groups can be imported using their `owner/slug` ID, e.g.

```
$ terraform import bitbucket_group.my-group my-account/my-group
```
