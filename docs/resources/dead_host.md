---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "nginxproxymanager_dead_host Resource - nginxproxymanager"
subcategory: "Hosts"
description: |-
  This resource can be used to manage a dead host.
---

# nginxproxymanager_dead_host (Resource)

This resource can be used to manage a dead host.


## Example Usage

```terraform
resource "nginxproxymanager_dead_host" "host" {
  domain_names = ["example.com"]

  certificate_id  = 1
  ssl_forced      = true
  hsts_enabled    = false
  hsts_subdomains = false
  http2_support   = true
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `domain_names` (Set of String) The domain names associated with the dead host.

### Optional

- `advanced_config` (String) The advanced configuration used by the dead host.
- `certificate_id` (Number) The Id of the certificate used by the dead host.
- `enabled` (Boolean) Whether the dead host is enabled.
- `hsts_enabled` (Boolean) Whether HSTS is enabled for the dead host.
- `hsts_subdomains` (Boolean) Whether HSTS is enabled for subdomains of the dead host.
- `http2_support` (Boolean) Whether HTTP/2 is supported for the dead host.
- `ssl_forced` (Boolean) Whether SSL is forced for the dead host.

### Read-Only

- `created_on` (String) The date and time the dead host was created.
- `id` (Number) The Id of the dead host.
- `meta` (Map of String) The meta data associated with the dead host.
- `modified_on` (String) The date and time the dead host was last modified.
- `owner_user_id` (Number) The Id of the user that owns the dead host.

## Import

Import is supported using the following syntax:

```shell
# 404 hosts can be imported by specifying the numeric identifier of the 404 host.
terraform import nginxproxymanager_dead_host.host 1
```
