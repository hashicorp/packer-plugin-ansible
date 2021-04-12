Ansible Components

The Ansible plugin allows users to execute as a provisioner during a Packer build.

-> **Note:** Ansible will _not_ be installed automatically by this
provisioner. This provisioner expects that Ansible is already installed on the
guest/remote machine. It is common practice to use the [shell
provisioner](/docs/provisioners/shell) before the Ansible provisioner to
do this.

Provisioners:

- [ansible](provisioner/ansible.mdx) - The Packer provisioner runs Ansible playbooks. It dynamically creates an Ansible inventory file configured to use SSH, runs an SSH server, executes ansible-playbook, and marshals Ansible plays through the SSH server to the machine being provisioned by Packer.

- [ansible-local](provisioner/ansible-local.mdx) - The Packer provisioner will run ansible in ansible's "local" mode on the remote/guest VM using Playbook and Role files that exist on the guest VM. This means ansible must be installed on the remote/guest VM.
