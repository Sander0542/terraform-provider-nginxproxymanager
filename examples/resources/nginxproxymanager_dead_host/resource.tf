resource "nginxproxymanager_dead_host" "host" {
  domain_names = ["example.com"]

  certificate_id  = 1
  ssl_forced      = true
  hsts_enabled    = false
  hsts_subdomains = false
  http2_support   = true
}
