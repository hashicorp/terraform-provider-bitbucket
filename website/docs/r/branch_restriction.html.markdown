---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_branch_restriction"
sidebar_current: "docs-bitbucket-resource-branch-restriction"
description: |-
  Provides a Bitbucket Branch Restriction
---

# bitbucket\_branch\_restriction

Provides a Bitbucket branch restriction resource.

This allows you for setting up branch restrictions for your repository.

## Example Usage

```hcl
# Manage your repositories branch restrictions
resource "bitbucket_branch_restriction" "master" {
  owner      = "myteam"
  repository = "terraform-code"

  kind = "push"
  pattern = "master"
}
```

## Argument Reference

The following arguments are supported:

* `owner` - (Required) The owner of this repository. Can be you or any team you
  have write access to.
* `repository` - (Required) The name of the repository.
* `kind` - (Required) The type of restriction that is being applied. List of possible stages is [here](https://developer.atlassian.com/bitbucket/api/2/reference/resource/repositories/%7Bworkspace%7D/%7Brepo_slug%7D/branch-restrictions/%7Bid%7Da).
* `pattern` - (Required) The pattern to determine which branches will be restricted.
* `users` - (Optional) A list of users to use.
* `groups` - (Optional) A list of groups to use.
