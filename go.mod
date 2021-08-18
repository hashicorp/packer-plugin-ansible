module github.com/hashicorp/packer-plugin-ansible

go 1.16

require (
	github.com/hashicorp/hcl/v2 v2.10.1
	github.com/hashicorp/packer-plugin-sdk v0.2.3
	github.com/stretchr/testify v1.7.0
	github.com/zclconf/go-cty v1.9.1
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad
)

// Incorrect plugin registration for ansible-local; see packer-plugin-ansible/pull/44
retract [v0.0.1, v0.0.2]
