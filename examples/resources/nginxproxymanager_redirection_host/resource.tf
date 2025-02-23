resource "nginxproxymanager_redirection_host" "host" {
  domain_names = ["example.com"]

  forward_http_code   = 300
  forward_scheme      = "auto"
  forward_domain_name = "example2.com"

  preserve_path  = true
  block_exploits = false

  certificate_id  = 1
  ssl_forced      = true
  hsts_enabled    = false
  hsts_subdomains = false
  http2_support   = true
}
