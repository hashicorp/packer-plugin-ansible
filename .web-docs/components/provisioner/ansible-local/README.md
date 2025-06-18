Type: `ansible-local`

The `ansible-local` Packer provisioner will execute `ansible` in Ansible's "local"
mode on the remote/guest VM using Playbook and Role files that exist on the
guest VM. This means Ansible must be installed on the remote/guest VM.
Playbooks and Roles can be uploaded from your build machine (the one running
Packer) to the vm. Ansible is then run on the guest machine in [local
mode](https://docs.ansible.com/ansible/latest/playbooks_delegation.html#local-playbooks)
via the `ansible-playbook` command.

-> **Note:** Ansible will _not_ be installed automatically by this
provisioner. This provisioner expects that Ansible is already installed on the
guest/remote machine. It is common practice to use the [shell
provisioner](/packer/docs/provisioner/shell) before the Ansible provisioner to
do this.

## Basic Example

The example below is fully functional.

**HCL2**

```hcl
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
    playbook_file   = "./playbook.yml"
    extra_arguments = ["--extra-vars", "\"pizza_toppings=${var.topping}\""]
  }
}
```

**JSON**

```json
{
  "builders": [
    {
      "type": "docker",
      "image": "williamyeh/ansible:ubuntu14.04",
      "export_path": "packer_example",
      "run_command": ["-d", "-i", "-t", "--entrypoint=/bin/bash", "{{.Image}}"]
    }
  ],
  "variables": {
    "topping": "mushroom"
  },
  "provisioners": [
    {
      "type": "ansible-local",
      "playbook_file": "./playbook.yml",
      "extra_arguments": [
        "--extra-vars",
        "\"pizza_toppings={{ user `topping`}}\""
      ]
    }
  ]
}
```


where ./playbook.yml contains

```
---
- name: hello world
  hosts: 127.0.0.1
  connection: local

  tasks:
    - command: echo {{ pizza_toppings }}
    - debug: msg="{{ pizza_toppings }}"

```

## Configuration Reference

The reference of available configuration options is listed below.

Required:

- `playbook_file` (string) - The playbook file to be executed by ansible.
  This file must exist on your local system and will be uploaded to the
  remote machine. This option is exclusive with `playbook_files`.

- `playbook_files` (array of strings) - The playbook files to be executed by
  ansible. These files must exist on your local system. If the files don't
  exist in the `playbook_dir` or you don't set `playbook_dir` they will be
  uploaded to the remote machine. This option is exclusive with
  `playbook_file`.

Optional:

<!-- Code generated from the comments of the Config struct in provisioner/ansible-local/provisioner.go; DO NOT EDIT MANUALLY -->

- `command` (string) - The command to invoke ansible. Defaults to
   `ansible-playbook`. If you would like to provide a more complex command,
   for example, something that sets up a virtual environment before calling
   ansible, take a look at the ansible wrapper guide [here](/packer/integrations/hashicorp/ansible/latest/components/provisioner/ansible#using-a-wrapping-script-for-your-ansible-call) for inspiration.
   Please note that Packer expects Command to be a path to an executable.
   Arbitrary bash scripting will not work and needs to go inside an
   executable script.

- `extra_arguments` ([]string) - Extra arguments to pass to Ansible.
  These arguments _will not_ be passed through a shell and arguments should
  not be quoted. Usage example:
  
  ```json
     "extra_arguments": [ "--extra-vars", "Region={{user `Region`}} Stage={{user `Stage`}}" ]
  ```
  
  In certain scenarios where you want to pass ansible command line arguments
  that include parameter and value (for example `--vault-password-file pwfile`),
  from ansible documentation this is correct format but that is NOT accepted here.
  Instead you need to do it like `--vault-password-file=pwfile`.
  
  If you are running a Windows build on AWS, Azure, Google Compute, or OpenStack
  and would like to access the auto-generated password that Packer uses to
  connect to a Windows instance via WinRM, you can use the template variable
  
  ```build.Password``` in HCL templates or ```{{ build `Password`}}``` in
  legacy JSON templates. For example:
  
  in JSON templates:
  
  ```json
  "extra_arguments": [
     "--extra-vars", "winrm_password={{ build `Password`}}"
  ]
  ```
  
  in HCL templates:
  ```hcl
  extra_arguments = [
     "--extra-vars", "winrm_password=${build.Password}"
  ]
  ```

- `group_vars` (string) - A path to the directory containing ansible group
  variables on your local system to be copied to the remote machine. By
  default, this is empty.

- `host_vars` (string) - A path to the directory containing ansible host variables on your local
  system to be copied to the remote machine. By default, this is empty.

- `playbook_dir` (string) - A path to the complete ansible directory structure on your local system
  to be copied to the remote machine as the `staging_directory` before all
  other files and directories.

- `playbook_file` (string) - The playbook file to be executed by ansible. This file must exist on your
  local system and will be uploaded to the remote machine. This option is
  exclusive with `playbook_files`.

- `playbook_files` ([]string) - The playbook files to be executed by ansible. These files must exist on
  your local system. If the files don't exist in the `playbook_dir` or you
  don't set `playbook_dir` they will be uploaded to the remote machine. This
  option is exclusive with `playbook_file`.

- `playbook_paths` ([]string) - An array of directories of playbook files on your local system. These
  will be uploaded to the remote machine under `staging_directory`/playbooks.
  By default, this is empty.

- `role_paths` ([]string) - An array of paths to role directories on your local system. These will be
  uploaded to the remote machine under `staging_directory`/roles. By default,
  this is empty.

- `collection_paths` ([]string) - An array of local paths of collections to upload.

- `staging_directory` (string) - The directory where files will be uploaded. Packer requires write
  permissions in this directory.

- `clean_staging_directory` (bool) - If set to `true`, the content of the `staging_directory` will be removed after
  executing ansible. By default this is set to `false`.

- `inventory_file` (string) - The inventory file to be used by ansible. This
  file must exist on your local system and will be uploaded to the remote
  machine.
  
  When using an inventory file, it's also required to `--limit` the hosts to the
  specified host you're building. The `--limit` argument can be provided in the
  `extra_arguments` option.
  
  An example inventory file may look like:
  
  ```text
  [chi-dbservers]
  db-01 ansible_connection=local
  db-02 ansible_connection=local
  
  [chi-appservers]
  app-01 ansible_connection=local
  app-02 ansible_connection=local
  
  [chi:children]
  chi-dbservers
  chi-appservers
  
  [dbservers:children]
  chi-dbservers
  
  [appservers:children]
  chi-appservers
  ```

- `inventory_groups` ([]string) - `inventory_groups` (string) - A comma-separated list of groups to which
  packer will assign the host `127.0.0.1`. A value of `my_group_1,my_group_2`
  will generate an Ansible inventory like:
  
  ```text
  [my_group_1]
  127.0.0.1
  [my_group_2]
  127.0.0.1
  ```

- `galaxy_file` (string) - A requirements file which provides a way to
   install roles or collections with the [ansible-galaxy
   cli](https://docs.ansible.com/ansible/latest/galaxy/user_guide.html#the-ansible-galaxy-command-line-tool)
   on the local machine before executing `ansible-playbook`. By default, this is empty.

- `galaxy_command` (string) - The command to invoke ansible-galaxy. By default, this is
  `ansible-galaxy`.

- `galaxy_force_install` (bool) - Force overwriting an existing role.
   Adds `--force` option to `ansible-galaxy` command. By default, this is
   `false`.

- `galaxy_roles_path` (string) - The path to the directory on the remote system in which to
    install the roles. Adds `--roles-path /path/to/your/roles` to
    `ansible-galaxy` command. By default, this will install to a 'galaxy_roles' subfolder in the
    staging/roles directory.

- `galaxy_collections_path` (string) - The path to the directory on the remote system in which to
    install the collections. Adds `--collections-path /path/to/your/collections` to
    `ansible-galaxy` command. By default, this will install to a 'galaxy_collections' subfolder in the
    staging/collections directory.

<!-- End of code generated from the comments of the Config struct in provisioner/ansible-local/provisioner.go; -->


Parameters common to all provisioners:

- `pause_before` (duration) - Sleep for duration before execution.

- `max_retries` (int) - Max times the provisioner will retry in case of failure. Defaults to zero (0). Zero means an error will not be retried.

- `only` (array of string) - Only run the provisioner for listed builder(s)
  by name.

- `override` (object) - Override the builder with different settings for a
  specific builder, eg :

  In HCL2:

  ```hcl
  source "null" "example1" {
    communicator = "none"
  }

  source "null" "example2" {
    communicator = "none"
  }

  build {
    sources = ["source.null.example1", "source.null.example2"]
    provisioner "shell-local" {
      inline = ["echo not overridden"]
      override = {
        example1 = {
          inline = ["echo yes overridden"]
        }
      }
    }
  }
  ```

  In JSON:

  ```json
  {
    "builders": [
      {
        "type": "null",
        "name": "example1",
        "communicator": "none"
      },
      {
        "type": "null",
        "name": "example2",
        "communicator": "none"
      }
    ],
    "provisioners": [
      {
        "type": "shell-local",
        "inline": ["echo not overridden"],
        "override": {
          "example1": {
            "inline": ["echo yes overridden"]
          }
        }
      }
    ]
  }
  ```

- `timeout` (duration) - If the provisioner takes more than for example
  `1h10m1s` or `10m` to finish, the provisioner will timeout and fail.


## Default Extra Variables

In addition to being able to specify extra arguments using the
`extra_arguments` configuration, the provisioner automatically defines certain
commonly useful Ansible variables:

- `packer_build_name` is set to the name of the build that Packer is running.
  This is most useful when Packer is making multiple builds and you want to
  distinguish them slightly when using a common playbook.

- `packer_builder_type` is the type of the builder that was used to create
  the machine that the script is running on. This is useful if you want to
  run only certain parts of the playbook on systems built with certain
  builders.

- `packer_http_addr` If using a builder that provides an HTTP server for file
  transfer (such as `hyperv`, `parallels`, `qemu`, `virtualbox`, and `vmware`), this
  will be set to the address. You can use this address in your provisioner to
  download large files over HTTP. This may be useful if you're experiencing
  slower speeds using the default file provisioner. A file provisioner using
  the `winrm` communicator may experience these types of difficulties.
