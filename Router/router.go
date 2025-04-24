package router

import (
	"EasyWeb/Service"
	"EasyWeb/Tools"
	_ "EasyWeb/Tools"
	"embed"
	_ "encoding/json"
	_ "fmt"
	"html/template"
	"io/fs"
	_ "net"
	"net/http"
	_ "strings"

	"github.com/gin-gonic/gin"
)

func GetRouter(content embed.FS) *gin.Engine {

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// 浏览界面处理
	r.SetHTMLTemplate(template.Must(template.New("").ParseFS(content, "static/*.html")))
	//r.LoadHTMLFiles("./static/index.html")
	// 注册静态文件:参数1：别名、参数2：当前static文件目录，
	//r.Static("static", "./static")
	fp, _ := fs.Sub(content, "static")
	r.StaticFS("static", http.FS(fp))

	// 注册路由
	r.GET("/", func(context *gin.Context) {
		context.HTML(http.StatusOK, "index.html", nil)
	})

	//api
	apiGroup := r.Group("/api")

	//获取节点信息
	apiGroup.GET("/peer", func(c *gin.Context) {
		data, err := Service.GetPeer()
		if err != nil {
			c.JSON(http.StatusOK, Tools.GetFailMsg(err.Error()))
			return
		}
		c.JSON(http.StatusOK, Tools.GetSuccMsg(len(data), data))
	})

	//获取当前设备信息
	apiGroup.GET("/node", func(c *gin.Context) {
		data, err := Service.GetNode()
		if err != nil {
			c.JSON(http.StatusOK, Tools.GetFailMsg(err.Error()))
			return
		}
		c.JSON(http.StatusOK, Tools.GetSuccMsg(1, data))

	})

	//获取服务器节点信息
	apiGroup.GET("/connector", func(c *gin.Context) {
		data, err := Service.GetConnector()
		if err != nil {
			c.JSON(http.StatusOK, Tools.GetFailMsg(err.Error()))
			return
		}
		c.JSON(http.StatusOK, Tools.GetSuccMsg(1, data))

	})

	return r
}
