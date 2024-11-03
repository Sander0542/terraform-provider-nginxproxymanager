package resources

type Stream struct {
	resourceWithOwner
	IncomingPort   uint16 `json:"incoming_port"`
	ForwardingHost string `json:"forwarding_host"`
	ForwardingPort uint16 `json:"forwarding_port"`
	TCPForwarding  bool   `json:"tcp_forwarding"`
	UDPForwarding  bool   `json:"udp_forwarding"`
	Enabled        bool   `json:"enabled"`
}

type StreamCollection []Stream
