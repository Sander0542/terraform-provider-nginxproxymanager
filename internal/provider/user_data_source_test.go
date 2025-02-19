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

func TestAccUserDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized
			{
				Config:      testAccUserDataSourceConfig + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Failed to authenticate with the Nginx Proxy Manager API"),
			},
			// Read testing
			{
				Config: testAccUserDataSourceConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.nginxproxymanager_user.test",
						tfjsonpath.New("id"),
						knownvalue.Int64Exact(1),
					),
					statecheck.ExpectKnownValue(
						"data.nginxproxymanager_user.test",
						tfjsonpath.New("created_on"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.nginxproxymanager_user.test",
						tfjsonpath.New("modified_on"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.nginxproxymanager_user.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact("Administrator"),
					),
					statecheck.ExpectKnownValue(
						"data.nginxproxymanager_user.test",
						tfjsonpath.New("nickname"),
						knownvalue.StringExact("Admin"),
					),
					statecheck.ExpectKnownValue(
						"data.nginxproxymanager_user.test",
						tfjsonpath.New("email"),
						knownvalue.StringExact("admin@example.com"),
					),
					statecheck.ExpectKnownValue(
						"data.nginxproxymanager_user.test",
						tfjsonpath.New("avatar"),
						knownvalue.StringExact(""),
					),
					statecheck.ExpectKnownValue(
						"data.nginxproxymanager_user.test",
						tfjsonpath.New("is_disabled"),
						knownvalue.Bool(false),
					),
					statecheck.ExpectKnownValue(
						"data.nginxproxymanager_user.test",
						tfjsonpath.New("roles"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("admin"),
						}),
					),
					statecheck.ExpectKnownValue(
						"data.nginxproxymanager_user.test",
						tfjsonpath.New("permissions"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"access_lists":      knownvalue.StringExact("manage"),
							"certificates":      knownvalue.StringExact("manage"),
							"dead_hosts":        knownvalue.StringExact("manage"),
							"proxy_hosts":       knownvalue.StringExact("manage"),
							"redirection_hosts": knownvalue.StringExact("manage"),
							"streams":           knownvalue.StringExact("manage"),
							"visibility":        knownvalue.StringExact("all"),
						}),
					),
				},
			},
		},
	})
}

const testAccUserDataSourceConfig = `
data "nginxproxymanager_user" "test" {
	id = 1
}
`
