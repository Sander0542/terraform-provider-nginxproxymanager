// Copyright (c) Sander Jochems
// SPDX-License-Identifier: MIT

package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccVersionDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized
			{
				Config:      testAccVersionDataSourceConfig + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Failed to authenticate with the Nginx Proxy Manager API"),
			},
			// Read testing
			{
				Config: testAccVersionDataSourceConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.nginxproxymanager_version.test",
						tfjsonpath.New("major"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.nginxproxymanager_version.test",
						tfjsonpath.New("minor"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.nginxproxymanager_version.test",
						tfjsonpath.New("revision"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.nginxproxymanager_version.test",
						tfjsonpath.New("version"),
						knownvalue.NotNull(),
					),
				},
			},
		},
	})
}

const testAccVersionDataSourceConfig = `
data "nginxproxymanager_version" "test" {}
`
