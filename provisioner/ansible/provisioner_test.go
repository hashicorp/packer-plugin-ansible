// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build !windows
// +build !windows

package ansible

import (
	"bytes"
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/hashicorp/packer-plugin-sdk/multistep/commonsteps"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	confighelper "github.com/hashicorp/packer-plugin-sdk/template/config"
	"github.com/stretchr/testify/assert"
)

// Be sure to remove the Ansible stub file in each test with:
//
//	defer os.Remove(config["command"].(string))
func testConfig(t *testing.T) map[string]interface{} {
	m := make(map[string]interface{})
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	ansible_stub := path.Join(wd, "packer-ansible-stub.sh")

	err = ioutil.WriteFile(ansible_stub, []byte("#!/usr/bin/env bash\necho ansible 1.6.0"), 0777)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	m["command"] = ansible_stub

	return m
}

func TestProvisioner_Impl(t *testing.T) {
	var raw interface{} = &Provisioner{}
	if _, ok := raw.(packersdk.Provisioner); !ok {
		t.Fatalf("must be a Provisioner")
	}
}

func TestProvisionerPrepare_Defaults(t *testing.T) {
	var p Provisioner
	config := testConfig(t)
	defer os.Remove(config["command"].(string))

	err := p.Prepare(config)
	if err == nil {
		t.Fatalf("should have error")
	}

	hostkey_file, err := ioutil.TempFile("", "hostkey")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.Remove(hostkey_file.Name())

	publickey_file, err := ioutil.TempFile("", "publickey")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.Remove(publickey_file.Name())

	playbook_file, err := ioutil.TempFile("", "playbook")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.Remove(playbook_file.Name())

	config["ssh_host_key_file"] = hostkey_file.Name()
	config["ssh_authorized_key_file"] = publickey_file.Name()
	config["playbook_file"] = playbook_file.Name()
	err = p.Prepare(config)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.Remove(playbook_file.Name())

	err = os.Unsetenv("USER")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	err = p.Prepare(config)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvisionerPrepare_PlaybookFile(t *testing.T) {
	var p Provisioner
	config := testConfig(t)
	defer os.Remove(config["command"].(string))

	hostkey_file, err := ioutil.TempFile("", "hostkey")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.Remove(hostkey_file.Name())

	publickey_file, err := ioutil.TempFile("", "publickey")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.Remove(publickey_file.Name())

	config["ssh_host_key_file"] = hostkey_file.Name()
	config["ssh_authorized_key_file"] = publickey_file.Name()

	err = p.Prepare(config)
	if err == nil {
		t.Fatal("should have error")
	}

	playbook_file, err := ioutil.TempFile("", "playbook")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.Remove(playbook_file.Name())

	config["playbook_file"] = playbook_file.Name()
	err = p.Prepare(config)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvisionerPrepare_HostKeyFile(t *testing.T) {
	var p Provisioner
	config := testConfig(t)
	defer os.Remove(config["command"].(string))

	publickey_file, err := ioutil.TempFile("", "publickey")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.Remove(publickey_file.Name())

	playbook_file, err := ioutil.TempFile("", "playbook")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.Remove(playbook_file.Name())

	filename := make([]byte, 10)
	n, err := io.ReadFull(rand.Reader, filename)
	if n != len(filename) || err != nil {
		t.Fatal("could not create random file name")
	}

	config["ssh_host_key_file"] = fmt.Sprintf("%x", filename)
	config["ssh_authorized_key_file"] = publickey_file.Name()
	config["playbook_file"] = playbook_file.Name()

	err = p.Prepare(config)
	if err == nil {
		t.Fatal("should error if ssh_host_key_file does not exist")
	}

	hostkey_file, err := ioutil.TempFile("", "hostkey")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.Remove(hostkey_file.Name())

	config["ssh_host_key_file"] = hostkey_file.Name()
	err = p.Prepare(config)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvisionerPrepare_AuthorizedKeyFile(t *testing.T) {
	var p Provisioner
	config := testConfig(t)
	defer os.Remove(config["command"].(string))

	hostkey_file, err := ioutil.TempFile("", "hostkey")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.Remove(hostkey_file.Name())

	playbook_file, err := ioutil.TempFile("", "playbook")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.Remove(playbook_file.Name())

	filename := make([]byte, 10)
	n, err := io.ReadFull(rand.Reader, filename)
	if n != len(filename) || err != nil {
		t.Fatal("could not create random file name")
	}

	config["ssh_host_key_file"] = hostkey_file.Name()
	config["playbook_file"] = playbook_file.Name()
	config["ssh_authorized_key_file"] = fmt.Sprintf("%x", filename)

	err = p.Prepare(config)
	if err == nil {
		t.Errorf("should error if ssh_authorized_key_file does not exist")
	}

	publickey_file, err := ioutil.TempFile("", "publickey")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.Remove(publickey_file.Name())

	config["ssh_authorized_key_file"] = publickey_file.Name()
	err = p.Prepare(config)
	if err != nil {
		t.Errorf("err: %s", err)
	}
}

func TestProvisionerPrepare_LocalPort(t *testing.T) {
	var p Provisioner
	config := testConfig(t)
	defer os.Remove(config["command"].(string))

	hostkey_file, err := ioutil.TempFile("", "hostkey")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.Remove(hostkey_file.Name())

	publickey_file, err := ioutil.TempFile("", "publickey")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.Remove(publickey_file.Name())

	playbook_file, err := ioutil.TempFile("", "playbook")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.Remove(playbook_file.Name())

	config["ssh_host_key_file"] = hostkey_file.Name()
	config["ssh_authorized_key_file"] = publickey_file.Name()
	config["playbook_file"] = playbook_file.Name()

	config["local_port"] = 65537
	err = p.Prepare(config)
	if err == nil {
		t.Fatal("should have error")
	}

	config["local_port"] = 22222
	err = p.Prepare(config)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvisionerPrepare_InventoryDirectory(t *testing.T) {
	var p Provisioner
	config := testConfig(t)
	defer os.Remove(config["command"].(string))

	hostkey_file, err := ioutil.TempFile("", "hostkey")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.Remove(hostkey_file.Name())

	publickey_file, err := ioutil.TempFile("", "publickey")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.Remove(publickey_file.Name())

	playbook_file, err := ioutil.TempFile("", "playbook")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.Remove(playbook_file.Name())

	config["ssh_host_key_file"] = hostkey_file.Name()
	config["ssh_authorized_key_file"] = publickey_file.Name()
	config["playbook_file"] = playbook_file.Name()

	config["inventory_directory"] = "doesnotexist"
	err = p.Prepare(config)
	if err == nil {
		t.Errorf("should error if inventory_directory does not exist")
	}

	inventoryDirectory, err := ioutil.TempDir("", "some_inventory_dir")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.Remove(inventoryDirectory)

	config["inventory_directory"] = inventoryDirectory
	err = p.Prepare(config)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestCreateInventoryFile(t *testing.T) {
	type inventoryFileTestCases struct {
		AnsibleVersion uint
		User           string
		Groups         []string
		EmptyGroups    []string
		UseProxy       confighelper.Trilean
		GeneratedData  map[string]interface{}
		Expected       string
	}

	TestCases := []inventoryFileTestCases{
		{
			AnsibleVersion: 1,
			User:           "testuser",
			UseProxy:       confighelper.TriFalse,
			GeneratedData:  basicGenData(nil),
			Expected:       "default ansible_ssh_host=123.45.67.89 ansible_ssh_user=testuser ansible_ssh_port=1234\n",
		},
		{
			AnsibleVersion: 2,
			User:           "testuser",
			UseProxy:       confighelper.TriFalse,
			GeneratedData:  basicGenData(nil),
			Expected:       "default ansible_host=123.45.67.89 ansible_user=testuser ansible_port=1234\n",
		},
		{
			AnsibleVersion: 1,
			User:           "testuser",
			Groups:         []string{"Group1", "Group2"},
			UseProxy:       confighelper.TriFalse,
			GeneratedData:  basicGenData(nil),
			Expected: `default ansible_ssh_host=123.45.67.89 ansible_ssh_user=testuser ansible_ssh_port=1234
[Group1]
default ansible_ssh_host=123.45.67.89 ansible_ssh_user=testuser ansible_ssh_port=1234
[Group2]
default ansible_ssh_host=123.45.67.89 ansible_ssh_user=testuser ansible_ssh_port=1234
`,
		},
		{
			AnsibleVersion: 1,
			User:           "testuser",
			EmptyGroups:    []string{"Group1", "Group2"},
			UseProxy:       confighelper.TriFalse,
			GeneratedData:  basicGenData(nil),
			Expected: `default ansible_ssh_host=123.45.67.89 ansible_ssh_user=testuser ansible_ssh_port=1234
[Group1]
[Group2]
`,
		},
		{
			AnsibleVersion: 1,
			User:           "testuser",
			Groups:         []string{"Group1", "Group2"},
			EmptyGroups:    []string{"Group3"},
			UseProxy:       confighelper.TriFalse,
			GeneratedData:  basicGenData(nil),
			Expected: `default ansible_ssh_host=123.45.67.89 ansible_ssh_user=testuser ansible_ssh_port=1234
[Group1]
default ansible_ssh_host=123.45.67.89 ansible_ssh_user=testuser ansible_ssh_port=1234
[Group2]
default ansible_ssh_host=123.45.67.89 ansible_ssh_user=testuser ansible_ssh_port=1234
[Group3]
`,
		},
		{
			AnsibleVersion: 2,
			User:           "testuser",
			UseProxy:       confighelper.TriFalse,
			GeneratedData: basicGenData(map[string]interface{}{
				"ConnType": "winrm",
				"Password": "12345",
			}),
			Expected: "default ansible_host=123.45.67.89 ansible_connection=winrm ansible_winrm_transport=basic ansible_shell_type=powershell ansible_user=testuser ansible_port=1234\n",
		},
	}

	for _, tc := range TestCases {
		var p Provisioner
		err := p.Prepare(testConfig(t))
		if err == nil {
			t.Fatalf("should have error")
		}
		defer os.Remove(p.config.Command)
		p.ansibleMajVersion = tc.AnsibleVersion
		p.config.User = tc.User
		p.config.Groups = tc.Groups
		p.config.EmptyGroups = tc.EmptyGroups
		p.config.UseProxy = tc.UseProxy
		p.generatedData = tc.GeneratedData

		err = p.createInventoryFile()
		if err != nil {
			t.Fatalf("error creating config using localhost and local port proxy")
		}
		if p.config.InventoryFile == "" {
			t.Fatalf("No inventory file was created")
		}
		defer os.Remove(p.config.InventoryFile)
		f, err := ioutil.ReadFile(p.config.InventoryFile)
		if err != nil {
			t.Fatalf("couldn't read created inventoryfile: %s", err)
		}

		expected := tc.Expected
		if string(f) != expected {
			t.Fatalf("File didn't match expected:\n\n expected: \n%s\n; recieved: \n%s\n", expected, f)
		}
	}
}

func basicGenData(input map[string]interface{}) map[string]interface{} {
	gd := map[string]interface{}{
		"Host":              "123.45.67.89",
		"Port":              int64(1234),
		"ConnType":          "ssh",
		"SSHPrivateKeyFile": "",
		"SSHPrivateKey":     "asdf",
		"SSHAgentAuth":      false,
		"User":              "PartyPacker",
		"PackerHTTPAddr":    commonsteps.HttpAddrNotImplemented,
		"PackerHTTPIP":      commonsteps.HttpIPNotImplemented,
		"PackerHTTPPort":    commonsteps.HttpPortNotImplemented,
	}
	if input == nil {
		return gd
	}
	for k, v := range input {
		gd[k] = v
	}
	return gd
}

func TestCreateCmdArgs(t *testing.T) {
	type testcase struct {
		TestName            string
		PackerBuildName     string
		PackerBuilderType   string
		UseProxy            confighelper.Trilean
		generatedData       map[string]interface{}
		AnsibleSSHExtraArgs []string
		ExtraArguments      []string
		AnsibleEnvVars      []string
		callArgs            []string // httpAddr inventory playbook privKeyFile
		ExpectedArgs        []string
		ExpectedEnvVars     []string
	}
	TestCases := []testcase{
		{
			// SSH with private key and an extra argument.
			TestName:        "SSH with private key and an extra argument",
			PackerBuildName: "packerparty",
			generatedData:   basicGenData(nil),
			ExtraArguments:  []string{"-e", "hello-world"},
			AnsibleEnvVars:  []string{"ENV_1=pancakes", "ENV_2=bananas"},
			callArgs:        []string{commonsteps.HttpAddrNotImplemented, "/var/inventory", "test-playbook.yml", "/path/to/privkey.pem"},
			ExpectedArgs:    []string{"-e", "packer_build_name=\"packerparty\"", "-e", "packer_builder_type=fakebuilder", "-e", "ansible_ssh_private_key_file=/path/to/privkey.pem", "--ssh-extra-args", "'-o IdentitiesOnly=yes'", "-e", "hello-world", "-i", "/var/inventory", "test-playbook.yml"},
			ExpectedEnvVars: []string{"ENV_1=pancakes", "ENV_2=bananas"},
		},
		{
			// SSH with private key and an extra argument.
			TestName:            "SSH with private key and an extra argument and a ssh extra argument",
			PackerBuildName:     "packerparty",
			generatedData:       basicGenData(nil),
			ExtraArguments:      []string{"-e", "hello-world"},
			AnsibleSSHExtraArgs: []string{"-o IdentitiesOnly=no"},
			AnsibleEnvVars:      []string{"ENV_1=pancakes", "ENV_2=bananas"},
			callArgs:            []string{commonsteps.HttpAddrNotImplemented, "/var/inventory", "test-playbook.yml", "/path/to/privkey.pem"},
			ExpectedArgs:        []string{"-e", "packer_build_name=\"packerparty\"", "-e", "packer_builder_type=fakebuilder", "--ssh-extra-args", "'-o IdentitiesOnly=no'", "-e", "ansible_ssh_private_key_file=/path/to/privkey.pem", "-e", "hello-world", "-i", "/var/inventory", "test-playbook.yml"},
			ExpectedEnvVars:     []string{"ENV_1=pancakes", "ENV_2=bananas"},
		},
		{
			// SSH with private key and an extra argument and multiple ssh extra arguments.
			TestName:            "SSH with private key and an extra argument and multiple ssh extra arguments",
			PackerBuildName:     "packerparty",
			generatedData:       basicGenData(nil),
			ExtraArguments:      []string{"-e", "hello-world"},
			AnsibleSSHExtraArgs: []string{"-o HostKeyAlgorithms=+ssh-rsa", "-o PubkeyAcceptedKeyTypes=+ssh-rsa", "-o IdentitiesOnly=no"},
			AnsibleEnvVars:      []string{"ENV_1=pancakes", "ENV_2=bananas"},
			callArgs:            []string{commonsteps.HttpAddrNotImplemented, "/var/inventory", "test-playbook.yml", "/path/to/privkey.pem"},
			ExpectedArgs:        []string{"-e", "packer_build_name=\"packerparty\"", "-e", "packer_builder_type=fakebuilder", "--ssh-extra-args", "'-o HostKeyAlgorithms=+ssh-rsa' '-o PubkeyAcceptedKeyTypes=+ssh-rsa' '-o IdentitiesOnly=no'", "-e", "ansible_ssh_private_key_file=/path/to/privkey.pem", "-e", "hello-world", "-i", "/var/inventory", "test-playbook.yml"},
			ExpectedEnvVars:     []string{"ENV_1=pancakes", "ENV_2=bananas"},
		},
		{
			TestName:        "SSH with private key and an extra argument and UseProxy",
			PackerBuildName: "packerparty",
			UseProxy:        confighelper.TriTrue,
			generatedData:   basicGenData(nil),
			ExtraArguments:  []string{"-e", "hello-world"},
			callArgs:        []string{commonsteps.HttpAddrNotImplemented, "/var/inventory", "test-playbook.yml", "/path/to/privkey.pem"},
			ExpectedArgs:    []string{"-e", "packer_build_name=\"packerparty\"", "-e", "packer_builder_type=fakebuilder", "-e", "ansible_ssh_private_key_file=/path/to/privkey.pem", "--ssh-extra-args", "'-o IdentitiesOnly=yes'", "-e", "hello-world", "-i", "/var/inventory", "test-playbook.yml"},
			ExpectedEnvVars: []string{},
		},
		{
			// Winrm, but no_proxy is unset so we don't do anything with ansible_password.
			TestName:        "Winrm, but no_proxy is unset so we don't do anything with ansible_password",
			PackerBuildName: "packerparty",
			generatedData: basicGenData(map[string]interface{}{
				"ConnType": "winrm",
			}),
			ExtraArguments:  []string{"-e", "hello-world"},
			AnsibleEnvVars:  []string{"ENV_1=pancakes", "ENV_2=bananas"},
			callArgs:        []string{commonsteps.HttpAddrNotImplemented, "/var/inventory", "test-playbook.yml", ""},
			ExpectedArgs:    []string{"-e", "packer_build_name=\"packerparty\"", "-e", "packer_builder_type=fakebuilder", "-e", "hello-world", "-i", "/var/inventory", "test-playbook.yml"},
			ExpectedEnvVars: []string{"ENV_1=pancakes", "ENV_2=bananas"},
		},
		{
			// HTTPAddr should be set. No env vars.
			TestName:        "HTTPAddr should be set. No env vars",
			PackerBuildName: "packerparty",
			ExtraArguments:  []string{"-e", "hello-world"},
			generatedData: basicGenData(map[string]interface{}{
				"PackerHTTPAddr": "123.45.67.89",
			}),
			callArgs:        []string{"123.45.67.89", "/var/inventory", "test-playbook.yml", ""},
			ExpectedArgs:    []string{"-e", "packer_build_name=\"packerparty\"", "-e", "packer_builder_type=fakebuilder", "-e", "packer_http_addr=123.45.67.89", "-e", "hello-world", "-i", "/var/inventory", "test-playbook.yml"},
			ExpectedEnvVars: []string{},
		},
		{
			// Add ansible_password for proxyless winrm connection.
			TestName: "Add ansible_password for proxyless winrm connection.",
			UseProxy: confighelper.TriFalse,
			generatedData: basicGenData(map[string]interface{}{
				"ConnType":       "winrm",
				"Password":       "ilovebananapancakes",
				"PackerHTTPAddr": "123.45.67.89",
			}),
			AnsibleEnvVars:  []string{"ENV_1=pancakes", "ENV_2=bananas"},
			callArgs:        []string{"123.45.67.89", "/var/inventory", "test-playbook.yml", ""},
			ExpectedArgs:    []string{"-e", "packer_builder_type=fakebuilder", "-e", "packer_http_addr=123.45.67.89", "-e", "ansible_password=ilovebananapancakes", "-i", "/var/inventory", "test-playbook.yml"},
			ExpectedEnvVars: []string{"ENV_1=pancakes", "ENV_2=bananas"},
		},
		{
			// Neither special ssh stuff, nor special windows stuff. This is docker!
			TestName:        "Neither special ssh stuff, nor special windows stuff. This is docker!",
			PackerBuildName: "packerparty",
			generatedData: basicGenData(map[string]interface{}{
				"ConnType": "docker",
			}),
			ExtraArguments:  []string{"-e", "hello-world"},
			AnsibleEnvVars:  []string{"ENV_1=pancakes", "ENV_2=bananas"},
			callArgs:        []string{commonsteps.HttpAddrNotImplemented, "/var/inventory", "test-playbook.yml", ""},
			ExpectedArgs:    []string{"-e", "packer_build_name=\"packerparty\"", "-e", "packer_builder_type=fakebuilder", "-e", "hello-world", "-i", "/var/inventory", "test-playbook.yml"},
			ExpectedEnvVars: []string{"ENV_1=pancakes", "ENV_2=bananas"},
		},
		{
			// Windows, no proxy, with extra vars.
			TestName: "Windows, no proxy, with extra vars.",
			UseProxy: confighelper.TriFalse,
			generatedData: basicGenData(map[string]interface{}{
				"ConnType":       "winrm",
				"Password":       "ilovebananapancakes",
				"PackerHTTPAddr": "123.45.67.89",
			}),
			ExtraArguments:  []string{"-e", "hello-world"},
			AnsibleEnvVars:  []string{"ENV_1=pancakes", "ENV_2=bananas"},
			callArgs:        []string{"123.45.67.89", "/var/inventory", "test-playbook.yml", ""},
			ExpectedArgs:    []string{"-e", "packer_builder_type=fakebuilder", "-e", "packer_http_addr=123.45.67.89", "-e", "ansible_password=ilovebananapancakes", "-e", "hello-world", "-i", "/var/inventory", "test-playbook.yml"},
			ExpectedEnvVars: []string{"ENV_1=pancakes", "ENV_2=bananas"},
		},
		{
			// SSH, use Password.
			TestName: "SSH, use ansible_password.",
			generatedData: basicGenData(map[string]interface{}{
				"ConnType":       "ssh",
				"Password":       "ilovebananapancakes",
				"PackerHTTPAddr": "123.45.67.89",
			}),
			ExtraArguments:  []string{"-e", "hello-world", "-e", "ansible_password=ilovebananapancakes"},
			AnsibleEnvVars:  []string{"ENV_1=pancakes", "ENV_2=bananas"},
			callArgs:        []string{"123.45.67.89", "/var/inventory", "test-playbook.yml", ""},
			ExpectedArgs:    []string{"-e", "packer_builder_type=fakebuilder", "-e", "packer_http_addr=123.45.67.89", "-e", "hello-world", "-e", "ansible_password=ilovebananapancakes", "-e", "ansible_host_key_checking=False", "-i", "/var/inventory", "test-playbook.yml"},
			ExpectedEnvVars: []string{"ENV_1=pancakes", "ENV_2=bananas"},
		},
		{
			// SSH, use Password .
			TestName: "SSH, already in ENV ansible_host_key_checking.",
			generatedData: basicGenData(map[string]interface{}{
				"ConnType":       "ssh",
				"Password":       "ilovebananapancakes",
				"PackerHTTPAddr": "123.45.67.89",
			}),
			ExtraArguments:  []string{"-e", "hello-world", "-e", "ansible_password=ilovebananapancakes"},
			AnsibleEnvVars:  []string{"ENV_1=pancakes", "ENV_2=bananas", "ANSIBLE_HOST_KEY_CHECKING=False"},
			callArgs:        []string{"123.45.67.89", "/var/inventory", "test-playbook.yml", ""},
			ExpectedArgs:    []string{"-e", "packer_builder_type=fakebuilder", "-e", "packer_http_addr=123.45.67.89", "-e", "hello-world", "-e", "ansible_password=ilovebananapancakes", "-i", "/var/inventory", "test-playbook.yml"},
			ExpectedEnvVars: []string{"ENV_1=pancakes", "ENV_2=bananas", "ANSIBLE_HOST_KEY_CHECKING=False"},
		},
		{
			TestName:        "Use PrivateKey",
			PackerBuildName: "packerparty",
			UseProxy:        confighelper.TriTrue,
			generatedData:   basicGenData(nil),
			ExtraArguments:  []string{"-e", "hello-world"},
			callArgs:        []string{commonsteps.HttpAddrNotImplemented, "/var/inventory", "test-playbook.yml", "/path/to/privkey.pem"},
			ExpectedArgs:    []string{"-e", "packer_build_name=\"packerparty\"", "-e", "packer_builder_type=fakebuilder", "-e", "ansible_ssh_private_key_file=/path/to/privkey.pem", "--ssh-extra-args", "'-o IdentitiesOnly=yes'", "-e", "hello-world", "-i", "/var/inventory", "test-playbook.yml"},
			ExpectedEnvVars: []string{},
		},
		{
			TestName:            "Use PrivateKey and SSH Extra Arg",
			PackerBuildName:     "packerparty",
			UseProxy:            confighelper.TriTrue,
			generatedData:       basicGenData(nil),
			AnsibleSSHExtraArgs: []string{"-o IdentitiesOnly=no"},
			ExtraArguments:      []string{"-e", "hello-world"},
			callArgs:            []string{commonsteps.HttpAddrNotImplemented, "/var/inventory", "test-playbook.yml", "/path/to/privkey.pem"},
			ExpectedArgs:        []string{"-e", "packer_build_name=\"packerparty\"", "-e", "packer_builder_type=fakebuilder", "-e", "ansible_ssh_private_key_file=/path/to/privkey.pem", "--ssh-extra-args", "'-o IdentitiesOnly=no'", "-e", "hello-world", "-i", "/var/inventory", "test-playbook.yml"},
			ExpectedEnvVars:     []string{},
		},
		{
			TestName:        "Use SSH Agent",
			UseProxy:        confighelper.TriTrue,
			generatedData:   basicGenData(nil),
			callArgs:        []string{commonsteps.HttpAddrNotImplemented, "/var/inventory", "test-playbook.yml", ""},
			ExpectedArgs:    []string{"-e", "packer_builder_type=fakebuilder", "-i", "/var/inventory", "test-playbook.yml"},
			ExpectedEnvVars: []string{},
		},
		{
			// No builder name. This shouldn't cause an error, it just shouldn't be set. HCL, yo.
			TestName:        "No builder name. This shouldn't cause an error, it just shouldn't be set. HCL, yo.",
			generatedData:   basicGenData(nil),
			callArgs:        []string{commonsteps.HttpAddrNotImplemented, "/var/inventory", "test-playbook.yml", ""},
			ExpectedArgs:    []string{"-e", "packer_builder_type=fakebuilder", "-i", "/var/inventory", "test-playbook.yml"},
			ExpectedEnvVars: []string{},
		},
	}

	for _, tc := range TestCases {
		var p Provisioner
		err := p.Prepare(testConfig(t))
		if err == nil {
			t.Fatalf("should have error")
		}
		defer os.Remove(p.config.Command)
		p.config.UseProxy = tc.UseProxy
		p.config.PackerBuilderType = "fakebuilder"
		p.config.PackerBuildName = tc.PackerBuildName
		p.generatedData = tc.generatedData
		p.config.AnsibleSSHExtraArgs = tc.AnsibleSSHExtraArgs
		p.config.ExtraArguments = tc.ExtraArguments
		p.config.AnsibleEnvVars = tc.AnsibleEnvVars

		args, envVars := p.createCmdArgs(tc.callArgs[0], tc.callArgs[1], tc.callArgs[2], tc.callArgs[3])
		assert.ElementsMatch(t, args, tc.ExpectedArgs,
			"TestName: %s\nArgs didn't match expected:\nexpected: \n%s\n; recieved: \n%s\n", tc.TestName, tc.ExpectedArgs, args)
		assert.ElementsMatch(t, envVars, tc.ExpectedEnvVars,
			"TestName: %s\nArgs didn't match expected:\n\nEnvVars didn't match expected:\n\n expected: \n%s\n; recieved: \n%s\n", tc.TestName, tc.ExpectedEnvVars, envVars)
		assert.EqualValues(t, tc.callArgs[2], args[len(args)-1],
			"TestName: %s\nPlayBook File Not Returned as last element: \nexpected: %s\nrecieved: %s\n", tc.TestName, tc.callArgs[2], args[len(args)-1])
	}
}

func TestUseProxy(t *testing.T) {
	type testcase struct {
		UseProxy                   confighelper.Trilean
		generatedData              map[string]interface{}
		expectedSetupAdapterCalled bool
		explanation                string
	}

	tcs := []testcase{
		{
			explanation:                "use_proxy is true; we should set up adapter",
			UseProxy:                   confighelper.TriTrue,
			generatedData:              basicGenData(nil),
			expectedSetupAdapterCalled: true,
		},
		{
			explanation: "use_proxy is false but no IP addr is available; we should set up adapter anyway.",
			UseProxy:    confighelper.TriFalse,
			generatedData: basicGenData(map[string]interface{}{
				"Host": "",
				"Port": nil,
			}),
			expectedSetupAdapterCalled: true,
		},
		{
			explanation:                "use_proxy is false; we shouldn't set up adapter.",
			UseProxy:                   confighelper.TriFalse,
			generatedData:              basicGenData(nil),
			expectedSetupAdapterCalled: false,
		},
		{
			explanation: "use_proxy is false but connType isn't ssh or winrm.",
			UseProxy:    confighelper.TriFalse,
			generatedData: basicGenData(map[string]interface{}{
				"ConnType": "docker",
			}),
			expectedSetupAdapterCalled: true,
		},
		{
			explanation:                "use_proxy is unset; we should default to setting up the adapter (for now).",
			UseProxy:                   confighelper.TriUnset,
			generatedData:              basicGenData(nil),
			expectedSetupAdapterCalled: true,
		},
		{
			explanation: "use_proxy is false and connType is winRM. we should not set up the adapter.",
			UseProxy:    confighelper.TriFalse,
			generatedData: basicGenData(map[string]interface{}{
				"ConnType": "winrm",
			}),
			expectedSetupAdapterCalled: false,
		},
		{
			explanation: "use_proxy is unset and connType is winRM. we should set up the adapter.",
			UseProxy:    confighelper.TriUnset,
			generatedData: basicGenData(map[string]interface{}{
				"ConnType": "winrm",
			}),
			expectedSetupAdapterCalled: true,
		},
	}

	for _, tc := range tcs {
		var p Provisioner
		err := p.Prepare(testConfig(t))
		if err == nil {
			t.Fatalf("%s should have error", tc.explanation)
		}
		p.config.UseProxy = tc.UseProxy
		defer os.Remove(p.config.Command)
		p.ansibleMajVersion = 1

		var l provisionLogicTracker
		l.setupAdapterCalled = false
		p.setupAdapterFunc = l.setupAdapter
		p.executeAnsibleFunc = l.executeAnsible
		ctx := context.TODO()
		comm := new(packersdk.MockCommunicator)
		ui := &packersdk.BasicUi{
			Reader: new(bytes.Buffer),
			Writer: new(bytes.Buffer),
		}
		//nolint:errcheck
		p.Provision(ctx, ui, comm, tc.generatedData)

		if l.setupAdapterCalled != tc.expectedSetupAdapterCalled {
			t.Fatalf("%s", tc.explanation)
		}
		os.Remove(p.config.Command)
	}
}

func TestProvisionerPrepare_WinRMSSL(t *testing.T) {
	type winrmUseHTTPValue int

	const (
		noWinRMUseHTTPOpt winrmUseHTTPValue = iota
		trueWinRMUseHTTPOpt
		falseWinRMUseHTTPOpt
	)

	testcases := []struct {
		name            string
		winRMUseHTTP    winrmUseHTTPValue
		extraArgs       []string
		expectedArgs    []string
		expectSchemeArg bool
	}{
		{
			name:            "ansible_winrm_use_http set to true, should specify -e ansible_winrm_scheme",
			winRMUseHTTP:    trueWinRMUseHTTPOpt,
			extraArgs:       []string{},
			expectedArgs:    []string{},
			expectSchemeArg: true,
		},
		{
			name:            "ansible_winrm_use_http set to true, option defined, should not specify -e ansible_winrm_scheme twice",
			winRMUseHTTP:    trueWinRMUseHTTPOpt,
			extraArgs:       []string{"-e", "ansible_winrm_scheme=http"},
			expectedArgs:    []string{},
			expectSchemeArg: true,
		},
		{
			name:            "ansible_winrm_use_http set to false, option defined, should not specify -e ansible_winrm_scheme twice",
			winRMUseHTTP:    falseWinRMUseHTTPOpt,
			extraArgs:       []string{"-e", "ansible_winrm_scheme=http"},
			expectedArgs:    []string{},
			expectSchemeArg: true,
		},
		{
			name:            "ansible_winrm_use_http set to false, option defined, should not specify -e ansible_winrm_scheme at all",
			winRMUseHTTP:    falseWinRMUseHTTPOpt,
			extraArgs:       []string{},
			expectedArgs:    []string{},
			expectSchemeArg: false,
		},
		{
			name:            "ansible_winrm_use_http set to false, option not defined, should not specify -e ansible_winrm_scheme",
			winRMUseHTTP:    noWinRMUseHTTPOpt,
			extraArgs:       []string{"-e", "ansible_winrm_scheme=http"},
			expectedArgs:    []string{},
			expectSchemeArg: true,
		},
		{
			name:            "ansible_winrm_use_http set to false, option not defined, should not specify -e ansible_winrm_scheme at all",
			winRMUseHTTP:    noWinRMUseHTTPOpt,
			extraArgs:       []string{},
			expectedArgs:    []string{},
			expectSchemeArg: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			config := testConfig(t)
			defer os.Remove(config["command"].(string))

			config["playbook_file"] = "test-fixtures/long-debug-message.yml"
			config["extra_arguments"] = tc.extraArgs

			switch tc.winRMUseHTTP {
			case noWinRMUseHTTPOpt:
			case trueWinRMUseHTTPOpt:
				config["ansible_winrm_use_http"] = true
			case falseWinRMUseHTTPOpt:
				config["ansible_winrm_use_http"] = false
			}

			var p Provisioner
			err := p.Prepare(config)
			if err != nil {
				t.Fatalf("prepare failed: %s", err)
			}

			ansibleArgs := p.config.ExtraArguments
			schemaArgCount := 0
			for _, arg := range ansibleArgs {
				if strings.HasPrefix(arg, "ansible_winrm_scheme") {
					schemaArgCount++
				}
			}

			if tc.expectSchemeArg && schemaArgCount != 1 {
				t.Errorf("expected exactly one ansible_winrm_scheme argument, got %d", schemaArgCount)
			}
			if !tc.expectSchemeArg && schemaArgCount != 0 {
				t.Errorf("expected ansible_winrm_scheme argument not to be defined, got %d", schemaArgCount)
			}
		})
	}
}
