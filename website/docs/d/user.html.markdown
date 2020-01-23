---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_user"
sidebar_current: "docs-bitbucket-data-user"
description: |-
  Provides a data for a Bitbucket user
---

# bitbucket\_user

Provdes a way to fetch data on a current user via there username, uuid or Display Name. 

## Example Usage

```hcl
# Manage your repository
data "bitbucket_user" "reviewer" {
  username = "gob"
}
```

## Argument Reference

The following arguments are supported:

* `username` - (Required)  the username
  have write access to.

## Exports

* `uuid` the uuid that bitbucket users to connect a user to various objects
* `display_name` the display name that the user wants to use for GDPR
* `nickname` typically the username but not always true.