package router

import (
	"EasyTier-Monitor/Service"
	"EasyTier-Monitor/Tools"
	"embed"
	"html/template"
	"io/fs"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetRouter(content embed.FS) *gin.Engine {

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// 添加自定义日志中间件
	r.Use(Tools.GinLogger())

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

	//api - 添加认证中间件
	apiGroup := r.Group("/api")
	apiGroup.Use(Tools.AuthMiddleware())

	//获取节点信息
	apiGroup.GET("/peer", func(c *gin.Context) {
		data, err := Service.GetPeerNew()
		if err != nil {
			c.JSON(http.StatusOK, Tools.GetFailMsg(err.Error()))
			return
		}
		c.JSON(http.StatusOK, Tools.GetSuccMsg(len(data), data))
	})

	//获取当前设备信息
	apiGroup.GET("/node", func(c *gin.Context) {
		data, err := Service.GetNodeNew()
		if err != nil {
			c.JSON(http.StatusOK, Tools.GetFailMsg(err.Error()))
			return
		}
		c.JSON(http.StatusOK, Tools.GetSuccMsg(1, data))
	})

	//获取服务器节点信息
	apiGroup.GET("/connector", func(c *gin.Context) {
		data, err := Service.GetConnectorNew()
		if err != nil {
			c.JSON(http.StatusOK, Tools.GetFailMsg(err.Error()))
			return
		}
		c.JSON(http.StatusOK, Tools.GetSuccMsg(1, data))
	})

	//获取状态页聚合数据
	apiGroup.GET("/status", func(c *gin.Context) {
		data, err := Service.GetStatus()
		if err != nil {
			c.JSON(http.StatusOK, Tools.GetFailMsg(err.Error()))
			return
		}
		c.JSON(http.StatusOK, Tools.GetSuccMsg(1, data))
	})

	//获取连接信息聚合数据
	apiGroup.GET("/connection", func(c *gin.Context) {
		data, err := Service.GetConnection()
		if err != nil {
			c.JSON(http.StatusOK, Tools.GetFailMsg(err.Error()))
			return
		}
		c.JSON(http.StatusOK, Tools.GetSuccMsg(1, data))
	})

	// SPA history 模式：未知的非接口、非静态请求统一返回 index.html
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/api/") || strings.HasPrefix(path, "/static/") {
			c.JSON(http.StatusNotFound, Tools.GetFailMsg("资源不存在"))
			return
		}
		c.HTML(http.StatusOK, "index.html", nil)
	})

	return r
}
