//go:generate packer-sdc mapstructure-to-hcl2 -type GalaxyConfig
//go:generate packer-sdc struct-markdown
package common

import (
	"io/ioutil"
	"regexp"
)

type GalaxyConfig struct {
	// A requirements file which provides a way to
	//  install roles or collections with the [ansible-galaxy
	//  cli](https://docs.ansible.com/ansible/latest/galaxy/user_guide.html#the-ansible-galaxy-command-line-tool)
	//  on the local machine before executing `ansible-playbook`. By default, this is empty.
	GalaxyFile string `mapstructure:"galaxy_file"`
	// The command to invoke ansible-galaxy. By default, this is
	// `ansible-galaxy`.
	GalaxyCommand string `mapstructure:"galaxy_command"`
	// Force overwriting an existing role.
	//  Adds `--force` option to `ansible-galaxy` command. By default, this is
	//  `false`.
	GalaxyForceInstall bool `mapstructure:"galaxy_force_install"`
	// The path to the directory on your local system in which to
	//   install the roles. Adds `--roles-path /path/to/your/roles` to
	//   `ansible-galaxy` command. By default, this is empty, and thus `--roles-path`
	//   option is not added to the command.
	GalaxyRolesPath string `mapstructure:"galaxy_roles_path"`
	// The path to the directory on your local system in which to
	//   install the collections. Adds `--collections-path /path/to/your/collections` to
	//   `ansible-galaxy` command. By default, this is empty, and thus `--collections-path`
	//   option is not added to the command.
	GalaxyCollectionsPath string `mapstructure:"galaxy_collections_path"`
	// Force overwriting an existing role and its dependencies.
	//  Adds `--force-with-deps` option to `ansible-galaxy` command. By default,
	//  this is `false`.
	GalaxyForceWithDeps bool `mapstructure:"galaxy_force_with_deps"`
}

func BuildGalaxyArgs(filepath, rolesPath, collectionsPath string, forceInstall, forceWithDeps bool) ([]string, error) {

	// ansible-galaxy install -r requirements.yml
	roleArgs := []string{"install", "-r", filepath, "-p", rolesPath}

	// Instead of modifying args depending on config values and removing or modifying values from
	// the slice between role and collection installs, just use 2 slices and simplify everything
	collectionArgs := []string{"collection", "install", "-r", filepath, "-p", collectionsPath}

	// Add force to arguments
	if forceInstall {
		roleArgs = append(roleArgs, "-f")
		collectionArgs = append(collectionArgs, "-f")
	}

	// Add --force-with-deps to arguments
	if forceWithDeps {
		roleArgs = append(roleArgs, "--force-with-deps")
		collectionArgs = append(collectionArgs, "--force-with-deps")
	}

	// Search galaxy_file for roles and collections keywords
	f, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	hasRoles, _ := regexp.Match(`(?m)^roles:`, f)
	hasCollections, _ := regexp.Match(`(?m)^collections:`, f)

	// If roles keyword present (v2 format), or no collections keyword present (v1), install roles
	if hasRoles || !hasCollections {
		return roleArgs, nil
	}

	// If collections keyword present (v2 format), install collections
	if hasCollections {
		return collectionArgs, nil
	}

	return nil, nil
}
