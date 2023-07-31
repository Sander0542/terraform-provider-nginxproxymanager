# Fetch a dead host by id (404 host)
data "nginxproxymanager_dead_host" "dead_host" {
  id = 1
}
