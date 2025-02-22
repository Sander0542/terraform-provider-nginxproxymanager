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
						tfjsonpath.New("settings"),
						knownvalue.SetExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"id":          knownvalue.StringExact("default-site"),
								"name":        knownvalue.StringExact("Default Site"),
								"description": knownvalue.StringExact("What to show when Nginx is hit with an unknown Host"),
								"value":       knownvalue.StringExact("congratulations"),
								"meta":        knownvalue.Null(),
							}),
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
