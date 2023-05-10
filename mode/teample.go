package mode

type Teample struct {
	Data    interface{} `json:"data,omitempty"`
	Msg     string      `json:"msg,omitempty"`
	Status  string      `json:"status,omitempty"`
	ModTime string      `json:"modTime,omitempty"`
	Mode    string      `json:"mode,omitempty"`
}
