package Tools

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware 创建认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 如果认证未启用，直接通过
		if !AppConfig.Security.Enabled {
			c.Next()
			return
		}

		// IP白名单检查
		if !checkIPWhitelist(c.ClientIP()) {
			AppLogger.Warn("IP %s 不在白名单中，拒绝访问", c.ClientIP())
			c.JSON(http.StatusForbidden, GetFailMsg("IP地址不在允许列表中"))
			c.Abort()
			return
		}

		// API密钥验证
		if !validateAPIKey(c) {
			AppLogger.Warn("无效的API密钥，来自IP: %s", c.ClientIP())
			c.JSON(http.StatusUnauthorized, GetFailMsg("无效的API密钥"))
			c.Abort()
			return
		}

		// 认证通过
		AppLogger.Debug("API认证通过，来自IP: %s", c.ClientIP())
		c.Next()
	}
}

// 检查IP是否在白名单中
func checkIPWhitelist(clientIP string) bool {
	allowedIPs := AppConfig.Security.AllowedIPs
	// 如果白名单为空，允许所有IP
	if allowedIPs == "" {
		return true
	}

	// 分割白名单IP
	ips := strings.Split(allowedIPs, ",")
	for _, ip := range ips {
		ip = strings.TrimSpace(ip)
		// 支持CIDR范围检查（简单实现，仅支持精确匹配）
		if ip == clientIP {
			return true
		}
	}

	return false
}

// 验证API密钥
func validateAPIKey(c *gin.Context) bool {
	expectedKey := AppConfig.Security.ApiKey
	// 如果API密钥为空，不需要验证
	if expectedKey == "" {
		return true
	}

	// 从请求头获取API密钥
	apiKey := c.GetHeader("X-API-Key")
	if apiKey == "" {
		// 从查询参数获取API密钥
		apiKey = c.Query("api_key")
	}

	return apiKey == expectedKey
}
