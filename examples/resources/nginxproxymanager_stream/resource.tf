resource "nginxproxymanager_stream" "stream" {
  incoming_port   = 22
  forwarding_host = "github.com"
  forwarding_port = 22

  tcp_forwarding = true
  udp_forwarding = false
}
