# Configuration-based authentication
provider "nginxproxymanager" {
  host     = "http://localhost:81"
  username = "terraform"
  password = "password"
}

# Environment variable-based authentication
provider "nginxproxymanager" {}
