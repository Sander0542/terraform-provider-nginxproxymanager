package resources

type Api struct {
	Status  string  `json:"status"`
	Version Version `json:"version"`
}
