The Ansible plugin allows users to execute as a provisioner during a Packer build.


### Installation
To install this plugin add this code into your Packer configuration and run [packer init](/packer/docs/commands/init)

```hcl
packer {
  required_plugins {
    ansible = {
      version = "~> 1"
      source = "github.com/hashicorp/ansible"
    }
  }
}
```

Alternatively, you can use `packer plugins install` to manage installation of this plugin.

```sh
packer plugins install github.com/hashicorp/ansible
```

### Components

**Note:** Ansible will _not_ be installed automatically by this provisioner. This provisioner expects that Ansible is already installed on the guest/remote machine.
It is common practice to use the [shell provisioner](/packer/docs/provisioners/shell) before the Ansible provisioner to do this.

#### Provisioners:

- [ansible](/packer/integrations/hashicorp/ansible/latest/components/provisioner/ansible) - The Packer provisioner runs Ansible playbooks. It dynamically creates an Ansible inventory file configured to use SSH, runs an SSH server, executes ansible-playbook, and marshals Ansible plays through the SSH server to the machine being provisioned by Packer.

- [ansible-local](/packer/integrations/hashicorp/ansible/latest/components/provisioner/ansibl-local) - The Packer provisioner will run ansible in ansible's "local" mode on the remote/guest VM using Playbook and Role files that exist on the guest VM. This means ansible must be installed on the remote/guest VM.
