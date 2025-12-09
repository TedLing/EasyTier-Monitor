package main

import (
	router "EasyTier-Monitor/Router"
	"EasyTier-Monitor/Tools"
	"embed"
	"fmt"

	"github.com/gin-gonic/gin"
)

//go:embed  static/*
var content embed.FS

func main() {
	// 加载配置
	if err := Tools.LoadConfig(""); err != nil {
		fmt.Printf("加载配置失败: %v, 使用默认配置\n", err)
	}

	// 前置参数校验
	if !checkCli() {
		return
	}

	// 初始化日志
	Tools.InitLogger()

	// 设置gin模式
	if Tools.AppConfig.Log.Level == Tools.LogLevelDebug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// 引入路由 传入打包的文件
	r := router.GetRouter(content)

	// 启动服务
	addr := fmt.Sprintf("%s:%d", Tools.AppConfig.Server.Host, Tools.AppConfig.Server.HttpPort)
	Tools.AppLogger.Info("尝试开启服务，监听地址: %s", addr)

	// 显示启动成功信息
	Tools.AppLogger.Info("******程序启动成功******")

	err := r.Run(addr)
	if err != nil {
		Tools.AppLogger.Error("启动程序错误：%v", err)
		return
	}
}

// 校验参数
func checkCli() bool {
	if len(Tools.AppConfig.CLI.Path) == 0 {
		fmt.Println("配置文件未配置cli.path，请检查！")
		return false
	}

	if !Tools.FileExist(Tools.AppConfig.CLI.Path) {
		fmt.Printf("配置文件配置的cli.path未找到文件: %s，请检查！\n", Tools.AppConfig.CLI.Path)
		return false
	}

	// 初始化临时日志以便记录
	if Tools.AppLogger == nil {
		Tools.InitLogger()
	}
	Tools.AppLogger.Info("CLI路径检查通过: %s", Tools.AppConfig.CLI.Path)

	return true
}

// 初始化函数
func init() {
	// 设置默认配置
	Tools.SetDefaultConfig()
}
