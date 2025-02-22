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

func TestAccSettingDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized
			{
				Config:      testAccSettingDataSourceConfig + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Failed to authenticate with the Nginx Proxy Manager API"),
			},
			// Read testing
			{
				Config: testAccSettingDataSourceConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.nginxproxymanager_setting.test",
						tfjsonpath.New("id"),
						knownvalue.StringExact("default-site"),
					),
					statecheck.ExpectKnownValue(
						"data.nginxproxymanager_setting.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact("Default Site"),
					),
					statecheck.ExpectKnownValue(
						"data.nginxproxymanager_setting.test",
						tfjsonpath.New("description"),
						knownvalue.StringExact("What to show when Nginx is hit with an unknown Host"),
					),
					statecheck.ExpectKnownValue(
						"data.nginxproxymanager_setting.test",
						tfjsonpath.New("value"),
						knownvalue.StringExact("congratulations"),
					),
					statecheck.ExpectKnownValue(
						"data.nginxproxymanager_setting.test",
						tfjsonpath.New("meta"),
						knownvalue.Null(),
					),
				},
			},
		},
	})
}

const testAccSettingDataSourceConfig = `
data "nginxproxymanager_setting" "test" {
	id = "default-site"
}
`
