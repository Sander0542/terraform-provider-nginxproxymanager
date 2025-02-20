resource "nginxproxymanager_certificate_custom" "certificate" {
  name = "Certificate"

  certificate     = file("certificate.pem")
  certificate_key = file("certificate.key")
}
