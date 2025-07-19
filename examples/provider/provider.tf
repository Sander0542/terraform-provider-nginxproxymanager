# Configuration-based authentication
provider "nginxproxymanager" {
  url      = "http://localhost:81"
  username = "admin@example.com"
  password = "changeme"
}

# Environment variable-based authentication
provider "nginxproxymanager" {}
