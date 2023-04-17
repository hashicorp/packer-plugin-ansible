# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

locals { timestamp = regex_replace(timestamp(), "[- TZ:]", "") }

source "googlecompute" "scp-to-sftp" {
  account_file             = "${var.account_file}"
  image_name               = "packerbats-scp-to-sftp-${local.timestamp}"
  project_id               = "${var.project_id}"
  source_image             = "debian-8-jessie-v20161027"
  ssh_file_transfer_method = "sftp"
  ssh_username             = "debian"
  zone                     = "us-central1-a"
}

build {
  sources = ["source.googlecompute.scp-to-sftp"]

  provisioner "ansible" {
    playbook_file = "./playbook.yml"
    sftp_command  = "/usr/bin/false"
    use_sftp      = false
  }

}
