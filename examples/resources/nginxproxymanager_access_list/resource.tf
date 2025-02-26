resource "nginxproxymanager_access_list" "access_list" {
  name = "RFC1918 with Auth"

  authorizations = [
    {
      username = "username"
      password = "password"
    }
  ]

  access = [
    {
      directive = "allow"
      address   = "10.0.0.0/8"
    },
    {
      directive = "allow"
      address   = "192.168.0.0/16"
    },
    {
      directive = "allow"
      address   = "172.16.0.0/12"
    }
  ]

  pass_auth   = false
  satisfy_any = true
}
