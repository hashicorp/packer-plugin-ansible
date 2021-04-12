Ansible Components

Provisioners:
ansible - The Packer provisioner runs Ansible playbooks. It dynamically creates an Ansible inventory file configured to use SSH, runs an SSH server, executes ansible-playbook, and marshals Ansible plays through the SSH server to the machine being provisioned by Packer.

ansible-local - The Packer provisioner will run ansible in ansible's "local" mode on the remote/guest VM using Playbook and Role files that exist on the guest VM. This means ansible must be installed on the remote/guest VM.
