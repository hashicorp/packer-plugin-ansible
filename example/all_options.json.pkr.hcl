# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

locals { timestamp = regex_replace(timestamp(), "[- TZ:]", "") }

source "googlecompute" "alloptions" {
  account_file = "${var.account_file}"
  image_name   = "packerbats-alloptions-${local.timestamp}"
  project_id   = "${var.project_id}"
  source_image = "debian-8-jessie-v20161027"
  ssh_username = "debian"
  zone         = "us-central1-a"
}

build {
  sources = ["source.googlecompute.alloptions"]

  provisioner "shell-local" {
    command = "echo 'TODO(bhcleek): write the public key to $HOME/.ssh/known_hosts and stop using ANSIBLE_HOST_KEY_CHECKING=False'"
  }

  provisioner "ansible" {
    ansible_env_vars        = ["PACKER_ANSIBLE_TEST=1", "ANSIBLE_HOST_KEY_CHECKING=False"]
    empty_groups            = ["PACKER_EMPTY_GROUP"]
    extra_arguments         = ["--private-key", "ansible-test-id"]
    groups                  = ["PACKER_TEST"]
    host_alias              = "packer-test"
    local_port              = 2222
    playbook_file           = "./playbook.yml"
    sftp_command            = "/usr/lib/sftp-server -e -l INFO"
    ssh_authorized_key_file = "ansible-test-id.pub"
    ssh_host_key_file       = "ansible-server.key"
    use_sftp                = true
    user                    = "packer"
  }

}
