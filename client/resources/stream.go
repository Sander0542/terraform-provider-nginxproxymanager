package resources

type Stream struct {
	resourceWithOwner
	IncomingPort   uint16    `json:"incoming_port"`
	ForwardingHost string    `json:"forwarding_host"`
	ForwardingPort uint16    `json:"forwarding_port"`
	TCPForwarding  boolAsInt `json:"tcp_forwarding"`
	UDPForwarding  boolAsInt `json:"udp_forwarding"`
	Enabled        boolAsInt `json:"enabled"`
}

type StreamCollection []Stream
