## 1.2.0 (Unreleased)
## 1.1.0 (June 19, 2019)

### Features

* add `skip_cert_verification` for hooks [#19]

### Bug fixes

* handle missing hooks [#24]
* fix default reviewer pagination bug [#28]

### Dev updates

* add `website` and `website-test` targets [#16]
* add `website-test` target to Travis [#17]
* upgrade to go 1.11 [#25]
* switch to go modules [#27]
* upgrade to `hashicorp/terraform` v0.12.2 [#34]

### Documentation

* add note about v1 APIs [#21]

## 1.0.0 (December 08, 2017)

* resource/bitbucket_repository: Add the ability to define a seperate slug for a repository ([#5](https://github.com/terraform-providers/terraform-provider-bitbucket/issues/5))

## 0.1.0 (June 20, 2017)

NOTES:

* Same functionality as that of Terraform 0.9.8. Repacked as part of [Provider Splitout](https://www.hashicorp.com/blog/upcoming-provider-changes-in-terraform-0-10/)
