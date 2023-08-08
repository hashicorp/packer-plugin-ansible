Type: `ansible`

The `ansible` Packer provisioner runs Ansible playbooks. It dynamically creates
an Ansible inventory file configured to use SSH, runs an SSH server, executes
`ansible-playbook`, and marshals Ansible plays through the SSH server to the
machine being provisioned by Packer.

-> **Note:** Any `remote_user` defined in tasks will be ignored. Packer
will always connect with the user given in the json config for this
provisioner.

-> **Note:** Options below that use the Packer template engine won't be able to
accept jinja2 `{{ function }}` macro syntax in a way that can be preserved to
the Ansible run. If you need to set variables using Ansible macros, you need to
do so inside your playbooks or inventory files.

Please see the [Debugging](#debugging), [Limitations](#limitations), or [Troubleshooting](#troubleshooting) if you are having trouble
getting started.

## Basic Example

This is a fully functional template that will provision an image on
DigitalOcean. Replace the mock `api_token` value with your own.

Example Packer template:

**HCL2**

```hcl
source "digitalocean" "example"{
    api_token = "6a561151587389c7cf8faa2d83e94150a4202da0e2bad34dd2bf236018ffaeeb"
    image     = "ubuntu-20-04-x64"
    region    = "sfo1"
}

build {
    sources = [
        "source.digitalocean.example"
    ]

    provisioner "ansible" {
      playbook_file = "./playbook.yml"
    }
}
```

**JSON**

```json
{
  "builders": [
    {
      "type": "digitalocean",
      "api_token": "6a561151587389c7cf8faa2d83e94150a4202da0e2bad34dd2bf236018ffaeeb",
      "image": "ubuntu-20-04-x64",
      "region": "sfo1"
    }
  ],
  "provisioners": [
    {
      "type": "ansible",
      "playbook_file": "./playbook.yml"
    }
  ]
}
```


Example playbook:

```yaml
---
# playbook.yml
- name: 'Provision Image'
  hosts: default
  become: true

  tasks:
    - name: install Apache
      package:
        name: 'httpd'
        state: present
```

## Configuration Reference

Required Parameters:

<!-- Code generated from the comments of the Config struct in provisioner/ansible/provisioner.go; DO NOT EDIT MANUALLY -->

- `playbook_file` (string) - The playbook to be run by Ansible.

<!-- End of code generated from the comments of the Config struct in provisioner/ansible/provisioner.go; -->


Optional Parameters:

<!-- Code generated from the comments of the Config struct in provisioner/ansible/provisioner.go; DO NOT EDIT MANUALLY -->

- `command` (string) - The command to invoke ansible. Defaults to
   `ansible-playbook`. If you would like to provide a more complex command,
   for example, something that sets up a virtual environment before calling
   ansible, take a look at the ansible wrapper guide below for inspiration.
   Please note that Packer expects Command to be a path to an executable.
   Arbitrary bash scripting will not work and needs to go inside an
   executable script.

- `extra_arguments` ([]string) - Extra arguments to pass to Ansible. These arguments _will not_ be passed
  through a shell and arguments should not be quoted. Usage example:
  
  ```json
     "extra_arguments": [ "--extra-vars", "Region={{user `Region`}} Stage={{user `Stage`}}" ]
  ```
  
  In certain scenarios where you want to pass ansible command line
  arguments that include parameter and value (for example
  `--vault-password-file pwfile`), from ansible documentation this is
  correct format but that is NOT accepted here. Instead you need to do it
  like `--vault-password-file=pwfile`.
  
  If you are running a Windows build on AWS, Azure, Google Compute, or
  OpenStack and would like to access the auto-generated password that
  Packer uses to connect to a Windows instance via WinRM, you can use the
  template variable
  
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
  
  If the lefthand side of a value contains 'secret' or 'password' (case
  insensitive) it will be hidden from output. For example, passing
  "my_password=secr3t" will hide "secr3t" from output.

- `ansible_env_vars` ([]string) - Environment variables to set before
    running Ansible. Usage example:
  
    ```json
      "ansible_env_vars": [ "ANSIBLE_HOST_KEY_CHECKING=False", "ANSIBLE_SSH_ARGS='-o ForwardAgent=yes -o ControlMaster=auto -o ControlPersist=60s'", "ANSIBLE_NOCOLOR=True" ]
    ```
  
    This is a [template engine](/packer/docs/templates/legacy_json_templates/engine). Therefore, you
    may use user variables and template functions in this field.
  
    For example, if you are running a Windows build on AWS, Azure,
    Google Compute, or OpenStack and would like to access the auto-generated
    password that Packer uses to connect to a Windows instance via WinRM, you
    can use the template variable `{{.WinRMPassword}}` in this option. Example:
  
    ```json
    "ansible_env_vars": [ "WINRM_PASSWORD={{.WinRMPassword}}" ],
    ```

- `ansible_ssh_extra_args` ([]string) - Specifies --ssh-extra-args on command line defaults to -o IdentitiesOnly=yes

- `groups` ([]string) - The groups into which the Ansible host should
   be placed. When unspecified, the host is not associated with any groups.

- `empty_groups` ([]string) - The groups which should be present in
   inventory file but remain empty.

- `host_alias` (string) - The alias by which the Ansible host should be
  known. Defaults to `default`. This setting is ignored when using a custom
  inventory file.

- `user` (string) - The `ansible_user` to use. Defaults to the user running
   packer, NOT the user set for your communicator. If you want to use the same
   user as the communicator, you will need to manually set it again in this
   field.

- `local_port` (int) - The port on which to attempt to listen for SSH
   connections. This value is a starting point. The provisioner will attempt
   listen for SSH connections on the first available of ten ports, starting at
   `local_port`. A system-chosen port is used when `local_port` is missing or
   empty.

- `ssh_host_key_file` (string) - The SSH key that will be used to run the SSH
   server on the host machine to forward commands to the target machine.
   Ansible connects to this server and will validate the identity of the
   server using the system known_hosts. The default behavior is to generate
   and use a onetime key. Host key checking is disabled via the
   `ANSIBLE_HOST_KEY_CHECKING` environment variable if the key is generated.

- `ssh_authorized_key_file` (string) - The SSH public key of the Ansible
   `ssh_user`. The default behavior is to generate and use a onetime key. If
   this key is generated, the corresponding private key is passed to
   `ansible-playbook` with the `-e ansible_ssh_private_key_file` option.

- `ansible_proxy_key_type` (string) - Change the key type used for the adapter.
  
  Supported values:
  
  * ECDSA (default)
  * RSA
  
  NOTE: using RSA may cause problems if the key is used to authenticate with rsa-sha1
  as modern OpenSSH versions reject this by default as it is unsafe.

- `sftp_command` (string) - The command to run on the machine being
   provisioned by Packer to handle the SFTP protocol that Ansible will use to
   transfer files. The command should read and write on stdin and stdout,
   respectively. Defaults to `/usr/lib/sftp-server -e`.

- `skip_version_check` (bool) - Check if ansible is installed prior to
   running. Set this to `true`, for example, if you're going to install
   ansible during the packer run.

- `use_sftp` (bool) - Use SFTP

- `inventory_directory` (string) - The directory in which to place the
   temporary generated Ansible inventory file. By default, this is the
   system-specific temporary file location. The fully-qualified name of this
   temporary file will be passed to the `-i` argument of the `ansible` command
   when this provisioner runs ansible. Specify this if you have an existing
   inventory directory with `host_vars` `group_vars` that you would like to
   use in the playbook that this provisioner will run.

- `inventory_file_template` (string) - This template represents the format for the lines added to the temporary
  inventory file that Packer will create to run Ansible against your image.
  The default for recent versions of Ansible is:
  "{{ .HostAlias }} ansible_host={{ .Host }} ansible_user={{ .User }} ansible_port={{ .Port }}\n"
  Available template engines are: This option is a template engine;
  variables available to you include the examples in the default (Host,
  HostAlias, User, Port) as well as any variables available to you via the
  "build" template engine.

- `inventory_file` (string) - The inventory file to use during provisioning.
   When unspecified, Packer will create a temporary inventory file and will
   use the `host_alias`.

- `keep_inventory_file` (bool) - If `true`, the Ansible provisioner will
   not delete the temporary inventory file it creates in order to connect to
   the instance. This is useful if you are trying to debug your ansible run
   and using "--on-error=ask" in order to leave your instance running while you
   test your playbook. this option is not used if you set an `inventory_file`.

- `galaxy_file` (string) - A requirements file which provides a way to
   install roles or collections with the [ansible-galaxy
   cli](https://docs.ansible.com/ansible/latest/galaxy/user_guide.html#the-ansible-galaxy-command-line-tool)
   on the local machine before executing `ansible-playbook`. By default, this is empty.

- `galaxy_command` (string) - The command to invoke ansible-galaxy. By default, this is
  `ansible-galaxy`.

- `galaxy_force_install` (bool) - Force overwriting an existing role.
   Adds `--force` option to `ansible-galaxy` command. By default, this is
   `false`.

- `galaxy_force_with_deps` (bool) - Force overwriting an existing role and its dependencies.
   Adds `--force-with-deps` option to `ansible-galaxy` command. By default,
   this is `false`.

- `roles_path` (string) - The path to the directory on your local system in which to
    install the roles. Adds `--roles-path /path/to/your/roles` to
    `ansible-galaxy` command. By default, this is empty, and thus `--roles-path`
    option is not added to the command.

- `collections_path` (string) - The path to the directory on your local system in which to
    install the collections. Adds `--collections-path /path/to/your/collections` to
    `ansible-galaxy` command. By default, this is empty, and thus `--collections-path`
    option is not added to the command.

- `use_proxy` (boolean) - When `true`, set up a localhost proxy adapter
  so that Ansible has an IP address to connect to, even if your guest does not
  have an IP address. For example, the adapter is necessary for Docker builds
  to use the Ansible provisioner. If you set this option to `false`, but
  Packer cannot find an IP address to connect Ansible to, it will
  automatically set up the adapter anyway.
  
   In order for Ansible to connect properly even when use_proxy is false, you
  need to make sure that you are either providing a valid username and ssh key
  to the ansible provisioner directly, or that the username and ssh key
  being used by the ssh communicator will work for your needs. If you do not
  provide a user to ansible, it will use the user associated with your
  builder, not the user running Packer.
   use_proxy=false is currently only supported for SSH and WinRM.
  
  Currently, this defaults to `true` for all connection types. In the future,
  this option will be changed to default to `false` for SSH and WinRM
  connections where the provisioner has access to a host IP.

- `ansible_winrm_use_http` (bool) - Force WinRM to use HTTP instead of HTTPS.
  
  Set this to true to force Ansible to use HTTP instead of HTTPS to communicate
  over WinRM to the destination host.
  
  Ansible uses the port as a heuristic to determine whether to use HTTP
  or not. In the current state, Packer assigns a random port for connecting
  to WinRM and Ansible's heuristic fails to determine that it should be
  using HTTP, even when the communicator is setup to use it.
  
  Alternatively, you may also directly add the following arguments to the
  `extra_arguments` section for ansible: `"-e", "ansible_winrm_scheme=http"`.
  
  Default: `false`

<!-- End of code generated from the comments of the Config struct in provisioner/ansible/provisioner.go; -->


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

## Debugging

To debug underlying issues with Ansible, add `"-vvvv"` to `"extra_arguments"`
to enable verbose logging.

**HCL2**

```hcl
  extra_arguments = [ "-vvvv" ]
```

**JSON**

```json
  "extra_arguments": [ "-vvvv" ]
```


## Limitations

### Redhat / CentOS

Redhat / CentOS builds have been known to fail with the following error due to
`sftp_command`, which should be set to `/usr/libexec/openssh/sftp-server -e`:

```text
==> virtualbox-ovf: starting sftp subsystem
    virtualbox-ovf: fatal: [default]: UNREACHABLE! => {"changed": false, "msg": "SSH Error: data could not be sent to the remote host. Make sure this host can be reached over ssh", "unreachable": true}
```

### chroot communicator

Building within a chroot (e.g. `amazon-chroot`) requires changing the Ansible
connection to chroot and running Ansible as root/sudo.

**HCL2**

```hcl
source "amazon-chroot" "example" {
  mount_path = "/mnt/packer-amazon-chroot"
  region     = "us-east-1"
  source_ami = "ami-123456"
}

build {
  sources = [
    "source.amazon-chroot.example"
  ]

  provisioner "ansible" {
    extra_arguments = [
        "--connection=chroot",
        "--inventory-file=/mnt/packer-amazon-chroot"
      ]

    playbook_file = "main.yml"
  }
}
```

**JSON**

```json
{
  "builders": [
    {
      "type": "amazon-chroot",
      "mount_path": "/mnt/packer-amazon-chroot",
      "region": "us-east-1",
      "source_ami": "ami-123456"
    }
  ],
  "provisioners": [
    {
      "type": "ansible",
      "extra_arguments": [
        "--connection=chroot",
        "--inventory-file=/mnt/packer-amazon-chroot"
      ],
      "playbook_file": "main.yml"
    }
  ]
}
```


### WinRM Communicator

There are two possible methods for using Ansible with the WinRM communicator.

Please note that if you're having trouble getting Ansible to connect, you may
want to take a look at the script that the Ansible project provides to help
configure remoting for Ansible:
https://github.com/ansible/ansible/blob/devel/examples/scripts/ConfigureRemotingForAnsible.ps1

#### Method 1 (recommended)

The recommended way to use the WinRM communicator is to set `"use_proxy": false`
and let the Ansible provisioner handle the rest for you. If you
are using WinRM with HTTPS, and you are using a self-signed certificate you
will also have to set `ansible_winrm_server_cert_validation=ignore` in your
extra_arguments.

Below is a fully functioning Ansible example for amazon-ebs using WinRM:

```hcl

variable "aws_region"{
  type    = string
  default = "us-east-1"
}

variable "aws_access_key"{
  type = string
}

variable "aws_secret_key"{
  type = string
}

data "amazon-ami" "windows" {

  filters = {
    name                = "*Windows_Server-2012*English-64Bit-Base*"
    root-device-type    = "ebs"
    virtualization-type = "hvm"
  }

  most_recent = true
  owners      = ["099720109477"]

  region      = var.aws_region
  access_key  = var.aws_access_key
  secret_key  = var.aws_secret_key
}

source "amazon-ebs" "example" {
    region      = var.aws_region
    instance_type = "t2.micro"

    source_ami       = data.amazon-ami.windows.id

    ami_name         = "test-ansible-packer"
    user_data_file   = "windows_bootstrap.txt"
    communicator     = "winrm"
    force_deregister = true
    winrm_username   = "Administrator"
    winrm_insecure   = true
    winrm_use_ssl    = true
}

build {
    sources = [
        "source.amazon-ebs.example",
    ]

    provisioner "ansible" {
      playbook_file   = "./playbooks/playbook-windows.yml"
      user            = "Administrator"
      use_proxy       = false
      extra_arguments = [
        "-e",
        "ansible_winrm_server_cert_validation=ignore"
      ]
    }
}
```

Below is a fully functioning Ansible example for azure-arm using WinRM.
Note: pywinrm needs to be installed into the python environment on your local build machine if it's not already installed.
Note: The ConfigureRemotingForAnsible.ps1 script can be found here https://github.com/ansible/ansible/blob/devel/examples/scripts/ConfigureRemotingForAnsible.ps1.

**HCL2**

```hcl
source "azure-arm" "server_2019" {
  use_azure_cli_auth                               = true
  build_resource_group_name                        = "ManagedImages-RGP"
  build_key_vault_name                             = "Example-Packer-Keyvault"
  os_type                                          = "Windows"
  image_publisher                                  = "MicrosoftWindowsServer"
  image_offer                                      = "WindowsServer"
  image_sku                                        = "2019-Datacenter"
  vm_size                                          = "Standard_D2as_v5"
  os_disk_size_gb                                  = 130
  shared_gallery_image_version_exclude_from_latest = false
  virtual_network_resource_group_name              = "VNET-Resource-Group"
  virtual_network_name                             = "My-VNET"
  virtual_network_subnet_name                      = "My-Subnet"
  private_virtual_network_with_public_ip           = false
  communicator                                     = "winrm"
  winrm_use_ssl                                    = true
  winrm_insecure                                   = true
  winrm_timeout                                    = "3m"
  winrm_username                                   = "Packer"
  managed_image_name                               = "Managed-Image-Name"
  managed_image_resource_group_name                = "ManagedImages-RGP"
  managed_image_storage_account_type               = "Standard_LRS"

  shared_image_gallery_destination {
    resource_group       = "ManagedImages-RGP"
    gallery_name         = "MyGallery"
    image_name           = "Server2019"
    storage_account_type = "Standard_LRS"
  }
}

build {
  sources = [
    "sources.azure-arm.server_2019",
  ]

  provisioner "shell-local" {
    inline_shebang = "/bin/bash -e"
    inline = [
      "pipx inject python-env-name \"pywinrm\"",
    ]
  }

  provisioner "powershell" {
    script = "../../scripts/ConfigureRemotingForAnsible.ps1"
  }

  provisioner "ansible" {
    skip_version_check  = false
    user                = "Packer"
    use_proxy           = false
    playbook_file       = "windows2019.yml"
    extra_arguments = [
      "-e",
      "ansible_winrm_server_cert_validation=ignore",
      "-e",
      "ansible_winrm_transport=ntlm",
    ]
  }
```

**JSON**

```json
{
  "builders": [
    {
      "type": "amazon-ebs",
      "region": "us-east-1",
      "instance_type": "t2.micro",
      "source_ami_filter": {
        "filters": {
          "virtualization-type": "hvm",
          "name": "*Windows_Server-2012*English-64Bit-Base*",
          "root-device-type": "ebs"
        },
        "most_recent": true,
        "owners": "amazon"
      },
      "ami_name": "test-ansible-packer",
      "user_data_file": "windows_bootstrap.txt",
      "communicator": "winrm",
      "force_deregister": true,
      "winrm_insecure": true,
      "winrm_username": "Administrator",
      "winrm_use_ssl": true
    }
  ],
  "provisioners": [
    {
      "type": "ansible",
      "playbook_file": "./playbook.yml",
      "user": "Administrator",
      "use_proxy": false,
      "extra_arguments": ["-e", "ansible_winrm_server_cert_validation=ignore"]
    }
  ]
}
```


Note that you do have to set the "Administrator" user, because otherwise Ansible
will default to using the user that is calling Packer, rather than the user
configured inside of the Packer communicator. For the contents of
windows_bootstrap.txt, see the WinRM docs for the amazon-ebs communicator.

When running from OSX, you may see an error like:

```text
    amazon-ebs: objc[9752]: +[__NSCFConstantString initialize] may have been in progress in another thread when fork() was called.
    amazon-ebs: objc[9752]: +[__NSCFConstantString initialize] may have been in progress in another thread when fork() was called. We cannot safely call it or ignore it in the fork() child process. Crashing instead. Set a breakpoint on objc_initializeAfterForkError to debug.
    amazon-ebs: ERROR! A worker was found in a dead state
```

If you see this, you may be able to work around the issue by telling Ansible to
explicitly not use any proxying; you can do this by setting the template option

**HCL2**

```hcl
ansible_env_vars = ["no_proxy=\"*\""]
```

**JSON**

```json
"ansible_env_vars": ["no_proxy=\"*\""],
```


in the above Ansible template.

#### Method 2 (Not recommended)

If you want to use the Packer SSH proxy, then you need a custom Ansible
connection plugin and a particular configuration. You need a directory named
`connection_plugins` next to the playbook which contains a file named
packer.py` which implements the connection plugin. On versions of Ansible
before 2.4.x, the following works as the connection plugin:

```python
from __future__ import (absolute_import, division, print_function)
__metaclass__ = type

from ansible.plugins.connection.ssh import Connection as SSHConnection

class Connection(SSHConnection):
    ''' ssh based connections for powershell via packer'''

    transport = 'packer'
    has_pipelining = True
    become_methods = []
    allow_executable = False
    module_implementation_preferences = ('.ps1', '')

    def __init__(self, *args, **kwargs):
        super(Connection, self).__init__(*args, **kwargs)
```

Newer versions of Ansible require all plugins to have a documentation string.
You can see if there is a plugin available for the version of Ansible you are
using
[here](https://github.com/hashicorp/packer/tree/master/provisioner/ansible/examples/connection-plugin).

To create the plugin yourself, you will need to copy all of the `options` from
the `DOCUMENTATION` string from the [ssh.py Ansible connection
plugin](https://github.com/ansible/ansible/blob/devel/lib/ansible/plugins/connection/ssh.py)
of the Ansible version you are using and add it to a packer.py file similar to
as follows

```python
from __future__ import (absolute_import, division, print_function)
__metaclass__ = type

from ansible.plugins.connection.ssh import Connection as SSHConnection

DOCUMENTATION = '''
    connection: packer
    short_description: ssh based connections for powershell via packer
    description:
        - This connection plugin allows ansible to communicate to the target packer machines via ssh based connections for powershell.
    author: Packer
    version_added: na
    options:
      **** Copy ALL the options from
      https://github.com/ansible/ansible/blob/devel/lib/ansible/plugins/connection/ssh.py
      for the version of Ansible you are using ****
'''

class Connection(SSHConnection):
    ''' ssh based connections for powershell via packer'''

    transport = 'packer'
    has_pipelining = True
    become_methods = []
    allow_executable = False
    module_implementation_preferences = ('.ps1', '')

    def __init__(self, *args, **kwargs):
        super(Connection, self).__init__(*args, **kwargs)
```

This template should build a Windows Server 2012 image on Google Cloud
Platform:

```json
{
  "variables": {},
  "provisioners": [
    {
      "type": "ansible",
      "playbook_file": "./win-playbook.yml",
      "extra_arguments": [
        "--connection",
        "packer",
        "--extra-vars",
        "\"ansible_shell_type=powershell ansible_shell_executable=None\""
      ]
    }
  ],
  "builders": [
    {
      "type": "googlecompute",
      "account_file": "{{ user `account_file`}}",
      "project_id": "{{user `project_id`}}",
      "source_image": "windows-server-2012-r2-dc-v20160916",
      "communicator": "winrm",
      "zone": "us-central1-a",
      "disk_size": 50,
      "winrm_username": "packer",
      "winrm_use_ssl": true,
      "winrm_insecure": true,
      "metadata": {
        "sysprep-specialize-script-cmd": "winrm set winrm/config/service/auth @{Basic=\"true\"}"
      }
    }
  ]
}
```

-> **Warning:** Please note that if you're setting up WinRM for provisioning, you'll probably want to turn it off or restrict its permissions as part of a shutdown script at the end of Packer's provisioning process. For more details on the why/how, check out this useful blog post and the associated code:
https://cloudywindows.io/post/winrm-for-provisioning-close-the-door-on-the-way-out-eh/

### Post i/o timeout errors

If you see
`unknown error: Post http://<ip>:<port>/wsman:dial tcp <ip>:<port>: i/o timeout`
errors while provisioning a Windows machine, try setting Ansible to copy files
over [ssh instead of
sftp](https://docs.ansible.com/ansible/latest/reference_appendices/config.html#envvar-ANSIBLE_SCP_IF_SSH).

### Too many SSH keys

SSH servers only allow you to attempt to authenticate a certain number of
times. All of your loaded keys will be tried before the dynamically generated
key. If you have too many SSH keys loaded in your `ssh-agent`, the Ansible
provisioner may fail authentication with a message similar to this:

```text
    googlecompute: fatal: [default]: UNREACHABLE! => {"changed": false, "msg": "Failed to connect to the host via ssh: Warning: Permanently added '[127.0.0.1]:62684' (RSA) to the list of known hosts.\r\nReceived disconnect from 127.0.0.1 port 62684:2: too many authentication failures\r\nAuthentication failed.\r\n", "unreachable": true}
```

To unload all keys from your `ssh-agent`, run:

```shell-session
$ ssh-add -D
```

### Become: yes

We recommend against running Packer as root; if you do then you won't be able
to successfully run your Ansible playbook as root; `become: yes` will fail.

### Using a wrapping script for your Ansible call

Sometimes, you may have extra setup that needs to be called as part of your
ansible run. The easiest way to do this is by writing a small bash script and
using that bash script in your `command` in place of the default
`ansible-playbook`. For example, you may need to launch a Python `virtualenv`
before calling Ansible. To do this, you'd want to create a bash script like

```shell
#!/bin/bash
source /tmp/venv/bin/activate && ANSIBLE_FORCE_COLOR=1 PYTHONUNBUFFERED=1 /tmp/venv/bin/ansible-playbook "$@"
```

The Ansible provisioner template remains very simple. For example:

**HCL2**

```hcl
provisioner "ansible" {
  command       = "/Path/To/call_ansible.sh"
  playbook_file = "./playbook.yml"
}
```

**JSON**

```json
{
  "type": "ansible",
  "command": "/Path/To/call_ansible.sh",
  "playbook_file": "./playbook.yml"
}
```


Note that we're calling ansible-playbook at the end of this command and passing
all command line arguments through into this call; this is necessary for
making sure that --extra-vars and other important Ansible arguments get set.
Note the quoting around the bash array, too; if you don't use quotes, any
arguments with spaces will not be read properly.

## Docker

When trying to use Ansible with Docker, it should "just work" but if it doesn't
you may need to tweak a few options.

- Change the ansible_connection from "ssh" to "docker"
- Set a Docker container name via the --name option.

On a CI server you probably want to overwrite ansible_host with a random name.

Example Packer template:

**HCL2**

```hcl
variable "ansible_host" {
  default = "default"
}

variable "ansible_connection" {
  default = "docker"
}

source "docker" "example" {
      image       = "centos:7"
      commit      = true
      run_command = [ "-d", "-i", "-t", "--name", var.ansible_host, "{{.Image}}", "/bin/bash" ]
}

build {
  sources = [
    "source.docker.example"
  ]

  provisioner "ansible" {
      groups          = [ "webserver" ]
      playbook_file   = "./webserver.yml"
      extra_arguments = [
          "--extra-vars",
          "ansible_host=${var.ansible_host} ansible_connection=${var.ansible_connection}"
      ]
  }
}
```

**JSON**

```json
{
  "variables": {
    "ansible_host": "default",
    "ansible_connection": "docker"
  },
  "builders": [
    {
      "type": "docker",
      "image": "centos:7",
      "commit": true,
      "run_command": [
        "-d",
        "-i",
        "-t",
        "--name",
        "{{user `ansible_host`}}",
        "{{.Image}}",
        "/bin/bash"
      ]
    }
  ],
  "provisioners": [
    {
      "type": "ansible",
      "groups": ["webserver"],
      "playbook_file": "./webserver.yml",
      "extra_arguments": [
        "--extra-vars",
        "ansible_host={{user `ansible_host`}} ansible_connection={{user `ansible_connection`}}"
      ]
    }
  ]
}
```


Example playbook:

```yaml
- name: configure webserver
  hosts: webserver
  tasks:
    - name: install Apache
      yum:
        name: httpd
```

## Amazon Session Manager

When trying to use Ansible with Amazon's Session Manager, you may run into an error where Ansible
is unable to connect to the remote Amazon instance if the local proxy adapter for Ansible [use_proxy](#use_proxy) is false.

The error may look something like the following:

```
amazon-ebs: fatal: [default]: UNREACHABLE! => {"changed": false, "msg": "Failed to connect to the host via ssh: ssh: connect to host 127.0.0.1 port 8362: Connection timed out", "unreachable": true}
```

The error is caused by a limitation on using Amazon's SSM default Port Forwarding session which only allows for one
remote connection on the forwarded port. Since Ansible's SSH communication is not using the local proxy adapter
it will try to make a new SSH connection to the same forwarded localhost port and fail.

In order to workaround this issue Ansible can be configured via a custom inventory file to use the AWS session-manager-plugin
directly to create a new session, separate from the one created by Packer, at runtime to connect and remotely provision the instance.

-> **Warning:** Please note that the default region configured for the `aws` cli must match the build region where the instance is being
provisioned otherwise you may run into a TargetNotConnected error. Users can use `AWS_DEFAULT_REGION` to temporarily override
their configured region.

**HCL2**

```hcl
  provisioner "ansible" {
      use_proxy               =  false
      playbook_file           =  "./playbooks/playbook_remote.yml"
      ansible_env_vars        =  ["PACKER_BUILD_NAME={{ build_name }}"]
      inventory_file_template =  "{{ .HostAlias }} ansible_host={{ .ID }} ansible_user={{ .User }} ansible_ssh_common_args='-o StrictHostKeyChecking=no -o ProxyCommand=\"sh -c \\\"aws ssm start-session --target %h --document-name AWS-StartSSHSession --parameters portNumber=%p\\\"\"'\n"
  }
```

**JSON**

```json
  "provisioners": [
    {
      "type": "ansible",
      "use_proxy": false,
      "ansible_env_vars": ["PACKER_BUILD_NAME={{ build_name }}"],
      "playbook_file":   "./playbooks/playbook_remote.yml",
      "inventory_file_template": "{{ .HostAlias }} ansible_host={{ .ID }} ansible_user={{ .User }} ansible_ssh_common_args='-o StrictHostKeyChecking=no -o ProxyCommand=\"sh -c \\\"aws ssm start-session --target %h --document-name AWS-StartSSHSession --parameters portNumber=%p\\\"\"'\n"
    }
  ]
```


Full Packer template example:

**HCL2**

```hcl

variables {
  instance_role = "SSMInstanceProfile"
}

source "amazon-ebs" "ansible-example" {
  region        = "us-east-1"
  ami_name      = "packer-ami-ansible"
  instance_type = "t2.micro"

  source_ami_filter {
      filters = {
          name                = "ubuntu/images/*ubuntu-xenial-16.04-amd64-server-*"
          virtualization-type = "hvm"
          root-device-type    = "ebs"
      }
      owners      = [ "099720109477" ]
      most_recent =  true
  }
  communicator         = "ssh"
  ssh_username         = "ubuntu"
  ssh_interface        = "session_manager"
  iam_instance_profile = var.instance_role
}

build {
  sources = ["source.amazon-ebs.ansible-example"]

  provisioner "ansible" {
      use_proxy               =  false
      playbook_file           =    "./playbooks/playbook_remote.yml"
      ansible_env_vars        =  ["PACKER_BUILD_NAME={{ build_name }}"]
      inventory_file_template =  "{{ .HostAlias }} ansible_host={{ .ID }} ansible_user={{ .User }} ansible_ssh_common_args='-o StrictHostKeyChecking=no -o ProxyCommand=\"sh -c \\\"aws ssm start-session --target %h --document-name AWS-StartSSHSession --parameters portNumber=%p\\\"\"'\n"
    }
}
```

**JSON**

```json
{
  "variables": {
    "instance_role": "SSMInstanceProfile"
  },

  "builders": [
    {
      "type": "amazon-ebs",
      "region": "us-east-1",
      "ami_name": "packer-ami-ansible",
      "instance_type": "t2.micro",
      "source_ami_filter": {
        "filters": {
          "virtualization-type": "hvm",
          "name": "ubuntu/images/*ubuntu-xenial-16.04-amd64-server-*",
          "root-device-type": "ebs"
        },
        "owners": ["099720109477"],
        "most_recent": true
      },
      "communicator": "ssh",
      "ssh_username": "ubuntu",
      "ssh_interface": "session_manager",
      "iam_instance_profile": "{{user `instance_role`}}"
    }
  ],
  "provisioners": [
    {
      "type": "ansible",
      "use_proxy": false,
      "ansible_env_vars": ["PACKER_BUILD_NAME={{ build_name }}"],
      "playbook_file": "./playbooks/playbook_remote.yml",
      "inventory_file_template": "{{ .HostAlias }} ansible_host={{ .ID }} ansible_user={{ .User }} ansible_ssh_common_args='-o StrictHostKeyChecking=no -o ProxyCommand=\"sh -c \\\"aws ssm start-session --target %h --document-name AWS-StartSSHSession --parameters portNumber=%p\\\"\"'\n"
    }
  ]
}
```


## Troubleshooting

If you are using an Ansible version >= 2.8 and Packer hangs in the
"Gathering Facts" stage, this could be the result of a pipelineing issue with
the proxy adapter that Packer uses. Setting `use_proxy` to `false` in the ansible
provisioner block of your Packer config should resolve the issue. In the future
we will default to setting this, so you won't have to but for now it is a manual
change you must make.
