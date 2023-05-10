package mode

type Host struct {
	HostName string `json:"hostName,omitempty"`
	Noteaddr string `json:"noteaddr,omitempty"` //note所在地址 : id  随机产生
}
