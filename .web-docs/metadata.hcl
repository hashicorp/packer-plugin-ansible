# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

# For full specification on the configuration of this file visit:
# https://github.com/hashicorp/integration-template#metadata-configuration
integration {
  name = "Ansible"
  description = "The Ansible plugin allows users to execute as a provisioner during a Packer build."
  identifier = "packer/hashicorp/ansible"
  component {
    type = "provisioner"
    name = "Ansible Local"
    slug = "ansible-local"
  }
  component {
    type = "provisioner"
    name = "Ansible"
    slug = "ansible"
  }
}
