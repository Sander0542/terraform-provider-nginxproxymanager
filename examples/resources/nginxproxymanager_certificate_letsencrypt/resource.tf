resource "nginxproxymanager_certificate_letsencrypt" "certificate" {
  domain_names = ["example.com", "*.example.com"]

  letsencrypt_email = "admin@example.com"
  letsencrypt_agree = true

  dns_challenge            = true
  dns_provider             = "cloudflare"
  dns_provider_credentials = "dns_cloudflare_api_token=0123456789abcdef0123456789abcdef01234567"
}
