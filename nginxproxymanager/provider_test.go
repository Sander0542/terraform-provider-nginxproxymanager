package nginxproxymanager

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"os"
	"testing"
)

const (
	providerConfig = `provider "nginxproxymanager" {}`
)

var (
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"nginxproxymanager": providerserver.NewProtocol6WithError(New("test")()),
	}
)

func testAccPreCheck(t *testing.T) {
	requiredEnvs := []string{
		"NGINX_PROXY_MANAGER_HOST",
		"NGINX_PROXY_MANAGER_USERNAME",
		"NGINX_PROXY_MANAGER_PASSWORD",
	}

	for _, env := range requiredEnvs {
		if _, result := os.LookupEnv(env); !result {
			t.Fatalf("Environment variable %s must be set for acceptance tests", env)
		}
	}
}
