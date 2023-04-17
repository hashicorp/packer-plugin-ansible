# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0


locals { timestamp = regex_replace(timestamp(), "[- TZ:]", "") }

source "googlecompute" "sftp" {
  account_file = "${var.account_file}"
  image_name   = "packerbats-sftp-${local.timestamp}"
  project_id   = "${var.project_id}"
  source_image = "debian-8-jessie-v20161027"
  ssh_username = "debian"
  zone         = "us-central1-a"
}

build {
  sources = ["source.googlecompute.sftp"]

  provisioner "ansible" {
    playbook_file = "./playbook.yml"
    sftp_command  = "/usr/lib/sftp-server -e -l INFO"
    use_sftp      = true
  }

}
