# Manage a proxy host
resource "nginxproxymanager_proxy_hosts" "example" {
  domain_names = ["example.com"]

  forward_scheme = "https"
  forward_host   = "example2.com"
  forward_port   = 443

  caching_enabled         = true
  allow_websocket_upgrade = true
  block_exploits          = true

  access_list_id = 0 # Publicly Accessible

  locations {
    path           = "/admin"
    forward_scheme = "https"
    forward_host   = "example3.com"
    forward_port   = 443

    advanced_config = ""
  }

  locations {
    path           = "/contact"
    forward_scheme = "http"
    forward_host   = "example4.com"
    forward_port   = 80

    advanced_config = ""
  }

  certificate_id  = 0 # No Certificate
  ssl_forced      = false
  hsts_enabled    = false
  hsts_subdomains = false
  http2_support   = false

  advanced_config = ""
}
