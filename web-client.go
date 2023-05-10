package main

import (
	"github.com/gin-gonic/gin"
	"hack8-note_rce/cotroller"
	"hack8-note_rce/middlewares"
)

func main() {

	r := gin.Default()
	r.Use(middlewares.Cors())
	r.Use(middlewares.FilterParams())
	r.GET("/down", cotroller.Download)
	r.GET("/getNote", cotroller.GetNotes)
	r.POST("/shell", cotroller.Shell)
	r.GET("/RefreshHost", cotroller.RefreshHost)
	r.GET("/dir", cotroller.ListDirHandler)
	r.POST("/upload", cotroller.Upload)
	r.GET("/fileDownload", cotroller.FileDownload)
	r.POST("/cs", cotroller.Cs)

	r.Run(":8080") // 监听并在 0.0.0.0:8080 上启动服务
}
