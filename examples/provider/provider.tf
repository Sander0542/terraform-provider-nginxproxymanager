# Configuration-based authentication
provider "nginxproxymanager" {
  host     = "http://localhost:81"
  username = "admin"
  password = "changeme"
}

# Environment variable-based authentication
provider "nginxproxymanager" {}
