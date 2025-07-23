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

**Note: Update to Packer Plugin Installation**

With the new Packer release starting from version 1.14.0, the `packer init` command will automatically install official plugins from the [HashiCorp release site.](https://releases.hashicorp.com/)

Going forward, to use newer versions of official Packer plugins, you'll need to upgrade to Packer version 1.14.0 or later. If you're using an older version, you can still install plugins, but as a workaround, you'll need to [manually install them using the CLI.](https://developer.hashicorp.com/packer/docs/plugins/install#manually-install-plugins-using-the-cli)

There is no change to the syntax or commands for installing plugins.

### Components

**Note:** Ansible will _not_ be installed automatically by this provisioner. This provisioner expects that Ansible is already installed on the guest/remote machine.
It is common practice to use the [shell provisioner](/packer/docs/provisioners/shell) before the Ansible provisioner to do this.

#### Provisioners:

- [ansible](/packer/integrations/hashicorp/ansible/latest/components/provisioner/ansible) - The Packer provisioner runs Ansible playbooks. It dynamically creates an Ansible inventory file configured to use SSH, runs an SSH server, executes ansible-playbook, and marshals Ansible plays through the SSH server to the machine being provisioned by Packer.

- [ansible-local](/packer/integrations/hashicorp/ansible/latest/components/provisioner/ansible-local) - The Packer provisioner will run ansible in ansible's "local" mode on the remote/guest VM using Playbook and Role files that exist on the guest VM. This means ansible must be installed on the remote/guest VM.
