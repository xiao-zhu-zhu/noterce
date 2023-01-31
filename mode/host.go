package mode

type Host struct {
	HostName string `json:"hostName,omitempty"`
	Id       string `json:"id,omitempty"`      //note所在地址 : id  随机产生
	Notekey  string `json:"notekey,omitempty"` //note所需执行命令的密码 : notekey
}
