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

func TestAccSettingsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized
			{
				Config:      testAccSettingsDataSourceConfig + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Failed to authenticate with the Nginx Proxy Manager API"),
			},
			// Read testing
			{
				Config: testAccSettingsDataSourceConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.nginxproxymanager_settings.test",
						tfjsonpath.New("default_site"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"page":     knownvalue.StringExact("congratulations"),
							"html":     knownvalue.Null(),
							"redirect": knownvalue.Null(),
						}),
					),
				},
			},
		},
	})
}

const testAccSettingsDataSourceConfig = `
data "nginxproxymanager_settings" "test" {}
`
