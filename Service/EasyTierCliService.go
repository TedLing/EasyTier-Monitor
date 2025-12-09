package Service

import (
	"EasyTier-Monitor/Model"
	"EasyTier-Monitor/Tools"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

func GetPeerNew() ([]Model.Peer, error) {
	// 先从缓存获取数据
	if cachedData, found := Tools.AppCache.Get(Tools.CacheKeyPeer); found {
		if peerData, ok := cachedData.([]Model.Peer); ok {
			Tools.AppLogger.Debug("从缓存获取Peer信息成功")
			return peerData, nil
		}
	}

	// 缓存不存在或已过期，执行命令获取新数据
	res, err := Tools.RunCmd(Tools.AppConfig.CLI.Path, "-o", "json", "peer")
	if err != nil {
		Tools.AppLogger.Error("执行peer命令失败: %v", err)
		return nil, fmt.Errorf("获取peer信息失败: %w", err)
	}

	var rows []Model.Peer
	err = json.Unmarshal([]byte(res), &rows)
	if err != nil {
		Tools.AppLogger.Error("解析peer JSON数据失败: %v, 原始数据: %s", err, res)
		return nil, fmt.Errorf("解析peer JSON数据失败: %w, 原始数据: %s", err, res)
	}

	// 将新数据存入缓存
	Tools.AppCache.Set(Tools.CacheKeyPeer, rows, Tools.GetCacheDuration(Tools.CacheKeyPeer))
	Tools.AppLogger.Info("获取Peer信息成功，共 %d 条记录", len(rows))

	return rows, nil

}

// GetNodeNew 新方式获取Node信息  直接返回JSON 目前支持 2.3.0
func GetNodeNew() (Model.Node, error) {
	// 先从缓存获取数据
	if cachedData, found := Tools.AppCache.Get(Tools.CacheKeyNode); found {
		if nodeData, ok := cachedData.(Model.Node); ok {
			Tools.AppLogger.Debug("从缓存获取Node信息成功")
			return nodeData, nil
		}
	}

	// 缓存不存在或已过期，执行命令获取新数据
	res, err := Tools.RunCmd(Tools.AppConfig.CLI.Path, "-o", "json", "node")
	if err != nil {
		Tools.AppLogger.Error("执行node命令失败: %v", err)
		return Model.Node{}, fmt.Errorf("获取node信息失败: %w", err)
	}

	nodeNew := Model.NodeNew{}
	err = json.Unmarshal([]byte(res), &nodeNew)
	if err != nil {
		Tools.AppLogger.Error("解析node JSON数据失败: %v, 原始数据: %s", err, res)
		return Model.Node{}, fmt.Errorf("解析node JSON数据失败: %w, 原始数据: %s", err, res)
	}

	//重新拼装处理 兼容前端逻辑
	var nodeInfo Model.Node
	nodeInfo.VirtualIP = nodeNew.IPv4Addr
	nodeInfo.Hostname = nodeNew.Hostname
	nodeInfo.ProxyCIDRs = strings.Join(nodeNew.ProxyCidrs, ",")
	nodeInfo.PeerID = fmt.Sprintf("%d", nodeNew.PeerID)
	nodeInfo.PublicIPv4 = Tools.ToIPv4(nodeNew.IPList.PublicIPv4) //ipv4 地址
	nodeInfo.PublicIPv6 = Tools.ToIPv6(nodeNew.IPList.PublicIPv6)
	nodeInfo.InterfaceIPv4 = Tools.ToIPv4List(nodeNew.IPList.InterfaceIPv4s)
	nodeInfo.InterfaceIPv6 = Tools.ToIPv6List(nodeNew.IPList.InterfaceIPv6s)
	nodeInfo.Listeners = nodeNew.Listeners

	// 将新数据存入缓存
	Tools.AppCache.Set(Tools.CacheKeyNode, nodeInfo, Tools.GetCacheDuration(Tools.CacheKeyNode))
	Tools.AppLogger.Info("获取Node信息成功")

	return nodeInfo, nil
}

func GetConnectorNew() ([]Model.ConnectorApi, error) {
	// 先从缓存获取数据
	if cachedData, found := Tools.AppCache.Get(Tools.CacheKeyConnector); found {
		if connectorData, ok := cachedData.([]Model.ConnectorApi); ok {
			Tools.AppLogger.Debug("从缓存获取Connector信息成功")
			return connectorData, nil
		}
	}

	// 缓存不存在或已过期，执行命令获取新数据
	res, err := Tools.RunCmd(Tools.AppConfig.CLI.Path, "-o", "json", "connector")
	if err != nil {
		Tools.AppLogger.Error("执行connector命令失败: %v", err)
		return nil, fmt.Errorf("获取connector信息失败: %w", err)
	}

	var ConnectorInfo []Model.Connector
	err = json.Unmarshal([]byte(res), &ConnectorInfo)
	if err != nil {
		Tools.AppLogger.Error("解析connector JSON数据失败: %v, 原始数据: %s", err, res)
		return nil, fmt.Errorf("解析connector JSON数据失败: %w, 原始数据: %s", err, res)
	}

	//转为兼容前端的输出
	var connectorApis []Model.ConnectorApi
	for _, connector := range ConnectorInfo {
		//定义局部变量和赋值
		var connectorApi Model.ConnectorApi
		connectorApi.Url = connector.Url.Url

		if connector.Status == 2 {
			connectorApi.Status = "连接中"
		} else if connector.Status == 1 {
			connectorApi.Status = "连接失败"
		} else if connector.Status == 0 {
			connectorApi.Status = "连接成功"
		} else {
			connectorApi.Status = strconv.Itoa(connector.Status)
		}

		//拼装到数组
		connectorApis = append(connectorApis, connectorApi)
	}

	// 将新数据存入缓存
	Tools.AppCache.Set(Tools.CacheKeyConnector, connectorApis, Tools.GetCacheDuration(Tools.CacheKeyConnector))
	Tools.AppLogger.Info("获取Connector信息成功")

	return connectorApis, nil

}
