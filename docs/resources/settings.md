---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "nginxproxymanager_settings Resource - nginxproxymanager"
subcategory: "Settings"
description: |-
  This resource can be used to manage settings.
---

# nginxproxymanager_settings (Resource)

This resource can be used to manage settings.


## Example Usage

```terraform
resource "nginxproxymanager_settings" "settings" {
  default_site = {
    page = "congratulations"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `default_site` (Attributes) What to show when Nginx is hit with an unknown Host. (see [below for nested schema](#nestedatt--default_site))

<a id="nestedatt--default_site"></a>
### Nested Schema for `default_site`

Required:

- `page` (String) What to show when Nginx is hit with an unknown Host.

Optional:

- `html` (String) HTML Content.
- `redirect` (String) Redirect to.
