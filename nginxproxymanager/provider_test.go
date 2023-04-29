package nginxproxymanager

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

const (
	providerConfig = `
provider "nginxproxymanager" {
	host = "http://localhost:81"
	username = "terraform"
	password = "terraform"
}
`
)

var (
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"nginxproxymanager": providerserver.NewProtocol6WithError(New()),
	}
)
