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

func TestAccUsersDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized
			{
				Config:      testAccUsersDataSourceConfig + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Failed to authenticate with the Nginx Proxy Manager API"),
			},
			// Read testing
			{
				Config: testAccUsersDataSourceConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.nginxproxymanager_users.test",
						tfjsonpath.New("users"),
						knownvalue.SetExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"id":          knownvalue.Int64Exact(1),
								"created_on":  knownvalue.NotNull(),
								"modified_on": knownvalue.NotNull(),
								"name":        knownvalue.StringExact("Administrator"),
								"nickname":    knownvalue.StringExact("Admin"),
								"email":       knownvalue.StringExact("admin@example.com"),
								"avatar":      knownvalue.StringExact(""),
								"is_disabled": knownvalue.Bool(false),
								"roles": knownvalue.SetExact([]knownvalue.Check{
									knownvalue.StringExact("admin"),
								}),
								"permissions": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"access_lists":      knownvalue.StringExact("manage"),
									"certificates":      knownvalue.StringExact("manage"),
									"dead_hosts":        knownvalue.StringExact("manage"),
									"proxy_hosts":       knownvalue.StringExact("manage"),
									"redirection_hosts": knownvalue.StringExact("manage"),
									"streams":           knownvalue.StringExact("manage"),
									"visibility":        knownvalue.StringExact("all"),
								}),
							}),
						}),
					),
				},
			},
		},
	})
}

const testAccUsersDataSourceConfig = `
data "nginxproxymanager_users" "test" {}
`
