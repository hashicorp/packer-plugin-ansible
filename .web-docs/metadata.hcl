# For full specification on the configuration of this file visit:
# https://github.com/hashicorp/integration-template#metadata-configuration
integration {
  name = "Ansible"
  description = "The Ansible plugin allows users to execute as a provisioner during a Packer build."
  identifier = "packer/BrandonRomano/ansible"
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
