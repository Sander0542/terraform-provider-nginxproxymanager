package resources

type UserPermissions struct {
	resource
	AccessLists      string `json:"access_lists"`
	Certificates     string `json:"certificates"`
	DeadHosts        string `json:"dead_hosts"`
	ProxyHosts       string `json:"proxy_hosts"`
	RedirectionHosts string `json:"redirection_hosts"`
	Streams          string `json:"streams"`
	Visibility       string `json:"visibility"`
	UserId           int64  `json:"user_id"`
}
