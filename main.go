package main

import (
	"context"

	"github.com/sander0542/terraform-provider-nginxproxymanager/nginxproxymanager"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

// Provider documentation generation.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-name nginxproxymanager

func main() {
	providerserver.Serve(context.Background(), nginxproxymanager.New, providerserver.ServeOpts{
		// NOTE: This is not a typical Terraform Registry provider address,
		// such as registry.terraform.io/Sander0542/nginxproxymanager. This specific
		// provider address is used in these tutorials in conjunction with a
		// specific Terraform CLI configuration for manual development testing
		// of this provider.
		Address: "hashicorp.com/edu/nginxproxymanager",
	})
}
