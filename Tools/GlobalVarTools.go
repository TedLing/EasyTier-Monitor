package Tools

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

// Config 应用配置结构体
type Config struct {
	// 服务器配置
	Server struct {
		HttpPort int    `yaml:"http_port"`
		Host     string `yaml:"host"`
	} `yaml:"server"`

	// CLI配置
	CLI struct {
		Path       string `yaml:"path"`
		Timeout    int    `yaml:"timeout"`
		MaxRetries int    `yaml:"max_retries"`
	} `yaml:"cli"`

	// 缓存配置
	Cache struct {
		Enabled      bool `yaml:"enabled"`
		DefaultTTL   int  `yaml:"default_ttl"`
		PeerTTL      int  `yaml:"peer_ttl"`
		NodeTTL      int  `yaml:"node_ttl"`
		ConnectorTTL int  `yaml:"connector_ttl"`
	} `yaml:"cache"`

	// 日志配置
	Log struct {
		Level      string `yaml:"level"`
		File       string `yaml:"file"`
		MaxSize    int    `yaml:"max_size"`
		MaxAge     int    `yaml:"max_age"`
		MaxBackups int    `yaml:"max_backups"`
		Compress   bool   `yaml:"compress"`
	} `yaml:"log"`

	// 安全配置
	Security struct {
		Enabled    bool   `yaml:"enabled"`
		ApiKey     string `yaml:"api_key"`
		AllowedIPs string `yaml:"allowed_ips"`
	} `yaml:"security"`
}

// AppConfig 全局配置实例
var AppConfig Config

// LoadConfig 加载配置文件
func LoadConfig(configPath string) error {
	// 设置默认配置
	SetDefaultConfig()

	// 如果没有指定配置文件路径，尝试从默认位置加载
	if configPath == "" {
		configPath = "config.yaml"
		// 检查是否存在
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			configPath = "config.yml"
			// 再次检查
			if _, err := os.Stat(configPath); os.IsNotExist(err) {
				// 没有找到配置文件，使用默认配置
				return nil
			}
		}
	}

	// 读取配置文件
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("读取配置文件失败: %v", err)
	}

	// 解析YAML配置
	if err := yaml.Unmarshal(data, &AppConfig); err != nil {
		return fmt.Errorf("解析配置文件失败: %v", err)
	}

	// 加载环境变量覆盖
	loadEnvOverrides()

	return nil
}

// SetDefaultConfig 设置默认配置
func SetDefaultConfig() {
	// 服务器默认配置
	AppConfig.Server.HttpPort = 2010
	AppConfig.Server.Host = "0.0.0.0"

	// CLI默认配置
	AppConfig.CLI.Path = "easytier"
	AppConfig.CLI.Timeout = 5
	AppConfig.CLI.MaxRetries = 3

	// 缓存默认配置
	AppConfig.Cache.Enabled = true
	AppConfig.Cache.DefaultTTL = 10
	AppConfig.Cache.PeerTTL = 10
	AppConfig.Cache.NodeTTL = 30
	AppConfig.Cache.ConnectorTTL = 30

	// 日志默认配置
	AppConfig.Log.Level = "info"
	AppConfig.Log.File = ""
	AppConfig.Log.MaxSize = 100
	AppConfig.Log.MaxAge = 7
	AppConfig.Log.MaxBackups = 3
	AppConfig.Log.Compress = false

	// 安全默认配置
	AppConfig.Security.Enabled = false
	AppConfig.Security.ApiKey = ""
	AppConfig.Security.AllowedIPs = ""
}

// loadEnvOverrides 从环境变量加载配置覆盖
func loadEnvOverrides() {
	// 服务器配置
	if port, exists := os.LookupEnv("EASYTIER_MONITOR_HTTP_PORT"); exists {
		if p, err := strconv.Atoi(port); err == nil {
			AppConfig.Server.HttpPort = p
		}
	}
	if host, exists := os.LookupEnv("EASYTIER_MONITOR_HOST"); exists {
		AppConfig.Server.Host = host
	}

	// CLI配置
	if path, exists := os.LookupEnv("EASYTIER_MONITOR_CLI_PATH"); exists {
		AppConfig.CLI.Path = path
	}
	if timeout, exists := os.LookupEnv("EASYTIER_MONITOR_CLI_TIMEOUT"); exists {
		if t, err := strconv.Atoi(timeout); err == nil {
			AppConfig.CLI.Timeout = t
		}
	}

	// 安全配置
	if enabled, exists := os.LookupEnv("EASYTIER_MONITOR_SECURITY_ENABLED"); exists {
		if b, err := strconv.ParseBool(enabled); err == nil {
			AppConfig.Security.Enabled = b
		}
	}
	if apiKey, exists := os.LookupEnv("EASYTIER_MONITOR_API_KEY"); exists {
		AppConfig.Security.ApiKey = apiKey
	}
}
