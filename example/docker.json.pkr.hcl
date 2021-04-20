
source "docker" "debian" {
  discard = true
  image   = "debian:jessie"
}

build {
  sources = ["source.docker.debian"]

  provisioner "shell" {
    inline = ["apt-get update", "apt-get -y install python"]
  }

  provisioner "ansible" {
    playbook_file = "./playbook.yml"
    sftp_command  = "/usr/bin/false"
    use_sftp      = false
  }

}
