<!-- archived-provider -->
Please note: This Terraform provider is archived, per our [provider archiving process](https://terraform.io/docs/internals/archiving.html). What does this mean?
1. The code repository and all commit history will still be available.
1. Existing released binaries will remain available on the releases site.
1. Issues and pull requests are not being monitored.
1. New releases will not be published.

If anyone from the community or an interested third party is willing to maintain it, they can fork the repository and [publish it](https://www.terraform.io/docs/registry/providers/publishing.html) to the Terraform Registry. If you are interested in maintaining this provider, please reach out to the [Terraform Provider Development Program](https://www.terraform.io/guides/terraform-provider-development-program.html) at *terraform-provider-dev@hashicorp.com*.

Terraform Provider
==================

- Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.10.x
-	[Go](https://golang.org/doc/install) 1.11 (to build the provider plugin)

Building The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/terraform-providers/terraform-provider-bitbucket`

```sh
$ mkdir -p $GOPATH/src/github.com/terraform-providers; cd $GOPATH/src/github.com/terraform-providers
$ git clone git@github.com:terraform-providers/terraform-provider-bitbucket
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/terraform-providers/terraform-provider-bitbucket
$ make build
```

Using the provider
----------------------

```hcl
# Configure the Bitbucket Provider
provider "bitbucket" {
  username = "GobBluthe"
  password = "idoillusions" # you can also use app passwords
}

# Manage your repository
resource "bitbucket_repository" "infrastructure" {
  owner = "myteam"
  name  = "terraform-code"
}

# Manage your project
resource "bitbucket_project" "infrastructure" {
  owner = "myteam" # must be a team
  name  = "terraform-project"
  key   = "TERRAFORMPROJ"
}
```

If `username` or `password` are not given, they'll be fetched from environment
variables `BITBUCKET_USERNAME` and `BITBUCKET_PASSWORD`.

Developing the Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.11+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-bitbucket
...
```

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Terraform needs TF_ACC env variable set to run acceptance tests

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
$ make testacc
```

About V1 APIs
------------------

This provider will not take any PRs about the v1 apis that don't have v2
equivalents. Please only focus on v2 apis when adding new features to this
provider.
