---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "nginxproxymanager Provider"
description: |-
  Use the Nginx Proxy Manager (NPM) provider to interact with resources from Nginx Proxy Manager.
---

# nginxproxymanager Provider

Use the Nginx Proxy Manager (NPM) provider to interact with resources from Nginx Proxy Manager.

## Example Usage

```terraform
# Configuration-based authentication
provider "nginxproxymanager" {
  url      = "http://localhost:81"
  username = "admin@example.com"
  password = "changeme"
}

# Environment variable-based authentication
provider "nginxproxymanager" {}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `password` (String, Sensitive) Password for Nginx Proxy Manager authentication. Can be specified via the `NGINXPROXYMANAGER_PASSWORD` environment variable.
- `url` (String) Full Nginx Proxy Manager URL with protocol and port (e.g. `http://localhost:81`). You should **NOT** supply any path (`/api`), the SDK will use the appropriate paths. Can be specified via the `NGINXPROXYMANAGER_URL` environment variable.
- `username` (String) Username for Nginx Proxy Manager authentication. Can be specified via the `NGINXPROXYMANAGER_USERNAME` environment variable.
