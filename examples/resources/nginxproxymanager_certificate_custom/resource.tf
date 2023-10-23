# Manage a custom certificate
resource "nginxproxymanager_certificate_custom" "example" {
  name = "example.com"

  certificate     = file("example.pem")
  certificate_key = file("example.key")
}
