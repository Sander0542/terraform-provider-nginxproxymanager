resource "nginxproxymanager_proxy_host" "host" {
  domain_names = ["example.com"]

  forward_scheme = "https"
  forward_host   = "example2.com"
  forward_port   = 443

  caching_enabled         = true
  allow_websocket_upgrade = true
  block_exploits          = true

  access_list_id = 1

  locations = [
    {
      path           = "/admin"
      forward_scheme = "https"
      forward_host   = "example3.com"
      forward_port   = 443
    },
    {
      path           = "/contact"
      forward_scheme = "http"
      forward_host   = "example4.com"
      forward_port   = 80
    }
  ]

  certificate_id  = 1
  ssl_forced      = true
  hsts_enabled    = false
  hsts_subdomains = false
  http2_support   = true
}
