# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

variable "topping" {
  type    = string
  default = "mushroom"
}

source "docker" "example" {
  image       = "williamyeh/ansible:ubuntu14.04"
  export_path = "packer_example"
  run_command = ["-d", "-i", "-t", "--entrypoint=/bin/bash", "{{.Image}}"]
}

build {
  sources = [
    "source.docker.example"
  ]

  provisioner "ansible-local" {
    playbook_file   = "./local-playbook.yml"
    extra_arguments = ["--extra-vars", "\"pizza_toppings=${var.topping}\""]
  }
}
