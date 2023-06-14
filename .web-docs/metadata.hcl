# For full specification on the configuration of this file visit:
# https://github.com/hashicorp/integration-template#metadata-configuration
integration {
  name = "Ansible"
  description = "TODO"
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
