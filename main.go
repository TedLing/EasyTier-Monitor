package main

import (
	"EasyTier-Monitor/Router"
	"EasyTier-Monitor/Tools"
	"embed"
	"fmt"
	"strconv"

	viper "github.com/spf13/viper"
)

//go:embed  static/*
var content embed.FS

func main() {

	//前置参数校验
	if !checkCli() {
		return
	}

	//引入路由 传入打包的文件
	r := router.GetRouter(content)

	// 启动服务
	strPort := strconv.Itoa(Tools.HttpPort)
	fmt.Println("尝试开启服务 端口号:" + strPort)

	//协程代码 显示当前启动成功的逻辑
	go func() {
		fmt.Println("******程序启动成功******")
	}()

	err := r.Run("0.0.0.0:" + strPort)
	if err != nil {
		fmt.Println("启动程序错误：" + err.Error())
		return
	}
}

// 校验参数 是不是等于空
func checkCli() bool {

	if len(Tools.CliPath) == 0 {
		fmt.Println("配置文件未配置[cli_path]请检查！")
		return false
	}

	if !Tools.FileExist(Tools.CliPath) {
		fmt.Println("配置文件配置的[cli_path]未找到文件，请检查！")
		return false
	}

	return true
}

// 初始化获取 配置文件
func init() {

	config := viper.New()
	config.AddConfigPath("./")
	config.SetConfigName("EasyTier-Monitor")
	config.SetConfigType("json")

	if err := config.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("未找到配置文件，端口号默认8899..")
			Tools.HttpPort = 8899
			return
		} else {
			fmt.Println("获取配置文件出错，端口号默认8899")
			Tools.HttpPort = 8899
			return
		}
	}

	Tools.HttpPort = config.GetInt("port")
	Tools.CliPath = config.GetString("cli_path")

}
