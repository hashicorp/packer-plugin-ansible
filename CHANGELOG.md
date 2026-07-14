## 1.1.6 (July 14, 2026)
### IMPROVEMENTS:

* deps: bump packer-plugin-sdk (0.6.9 -> 0.6.10)
  [GH-266](https://github.com/hashicorp/packer-plugin-ansible/pull/266)

* security/deps: upgrade golang.org/x/crypto and golang.org/x/net; refresh related golang.org/x indirect modules
  [GH-262](https://github.com/hashicorp/packer-plugin-ansible/pull/262)
  [GH-267](https://github.com/hashicorp/packer-plugin-ansible/pull/267)

## 1.1.5 (June 16, 2026)
### IMPROVEMENTS:

* deps: bump packer-plugin-sdk (0.6.2 -> 0.6.4 -> 0.6.9)
  [GH-242](https://github.com/hashicorp/packer-plugin-ansible/pull/242)
  [GH-258](https://github.com/hashicorp/packer-plugin-ansible/pull/258)

* security/deps: upgrade Azure NTLM, go-jose, and golang.org/x/crypto; include authorization-bypass go.mod updates
  [GH-260](https://github.com/hashicorp/packer-plugin-ansible/pull/260)
  [GH-245](https://github.com/hashicorp/packer-plugin-ansible/pull/245)
  [GH-250](https://github.com/hashicorp/packer-plugin-ansible/pull/250)

* ci/release: add backport workflow and prepare 1.1.5-dev release metadata
  [GH-241](https://github.com/hashicorp/packer-plugin-ansible/pull/241)
  [GH-238](https://github.com/hashicorp/packer-plugin-ansible/pull/238)

* build/compliance: remove illumos builds and update copyright/license headers
  [GH-235](https://github.com/hashicorp/packer-plugin-ansible/pull/235)
  [GH-244](https://github.com/hashicorp/packer-plugin-ansible/pull/244)

### BUG FIXES:

* core: fix golangci-lint staticcheck and errcheck findings in ansible provisioner


## 1.1.4 (July 30, 2025)
### IMPROVEMENTS:

* core: added environment vars to ansible-galaxy execution
  [GH-210](https://github.com/hashicorp/packer-plugin-ansible/pull/210)
  
* docs: update links to ansible wrapper guide for clarity and fix broken links
  [GH-212](https://github.com/hashicorp/packer-plugin-ansible/pull/212)

* docs: Update ansible script link to configure remoting for ansible
  [GH-205](https://github.com/hashicorp/packer-plugin-ansible/pull/205)

* Updated plugin release process: Plugin binaries are now published on the HashiCorp official [release site](https://releases.hashicorp.com/packer-plugin-ansible), ensuring a secure and standardized delivery pipeline.


### BUG FIXES:

* handle missing or invalid Host IP gracefully
  [GH-213](https://github.com/hashicorp/packer-plugin-ansible/pull/213)


## 1.0.0 (June 14, 2021)
The code base for this plugin has been stable since the Packer core split.
We are marking this plugin as v1.0.0 to indicate that it is stable and ready for consumption via `packer init`.

* Update packer-plugin-sdk to v0.2.3 [GH-48]
* Add module retraction for v0.0.1 as it was a bad release. [GH-46]


## 0.0.3 (May 11, 2021)

### BUG FIXES:
* Fix registration bug that externally caused plugin not to load properly [GH-44]

## 0.0.2 (April 15, 2021)

### BUG FIXES:
* core: Update module name in go.mod to fix plugin import path issue

## 0.0.1 (April 14, 2021)

* Ansible Plugin break out from Packer core. Changes prior to break out can be found in [Packer's CHANGELOG](https://github.com/hashicorp/packer/blob/master/CHANGELOG.md)

