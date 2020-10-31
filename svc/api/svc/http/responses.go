package http

type upperCaseResp struct {
	V string `json:"v"`
	Err string `json:",omitempty"`
}

type countResp struct {
	V int `json:"v"`
}
