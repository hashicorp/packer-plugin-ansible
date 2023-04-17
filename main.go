// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/packer-plugin-sdk/plugin"

	ansible "github.com/hashicorp/packer-plugin-ansible/provisioner/ansible"
	ansibleLocal "github.com/hashicorp/packer-plugin-ansible/provisioner/ansible-local"
	"github.com/hashicorp/packer-plugin-ansible/version"
)

func main() {
	pps := plugin.NewSet()
	pps.RegisterProvisioner(plugin.DEFAULT_NAME, new(ansible.Provisioner))
	pps.RegisterProvisioner("local", new(ansibleLocal.Provisioner))
	pps.SetVersion(version.PluginVersion)
	err := pps.Run()

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
