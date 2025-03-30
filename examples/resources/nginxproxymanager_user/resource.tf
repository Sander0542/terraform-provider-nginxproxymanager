resource "nginxproxymanager_user" "user" {
  name     = "John Doe"
  nickname = "john.doe"
  email    = "john.doe@example.com"

  is_admin = false

  permissions = {
    visibility = "all"

    proxy_hosts       = "manage"
    dead_hosts        = "manage"
    redirection_hosts = "manage"
    streams           = "manage"

    access_lists = "view"
    certificates = "view"
  }
}
