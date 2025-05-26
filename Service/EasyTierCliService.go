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
	//直接执行命令
	res, err := Tools.RunCmd(Tools.CliPath, "-o", "json", "peer")
	if err != nil {
		fmt.Println("执行cmd命令失败：", err)
		return nil, err
	}

	var rows []Model.Peer
	err = json.Unmarshal([]byte(res), &rows)
	if err != nil {
		fmt.Println("解析json失败：", err)
		return nil, err
	}

	return rows, nil

}

// GetPeer 获取客户端连接列表
func GetPeer() ([]Model.Peer, error) {

	//直接执行命令
	res, err := Tools.RunCmd(Tools.CliPath, "peer")
	if err != nil {
		fmt.Println("执行cmd命令失败：", err)
		return nil, err
	}

	// 解析表格并转换为 JSON
	jsonOutput, err := Tools.ParsePeerTableToJSON(res)
	if err != nil {
		fmt.Println("解析返回数据失败：", err)
		return nil, err
	}

	var rows []Model.Peer
	err = json.Unmarshal([]byte(jsonOutput), &rows)
	if err != nil {
		fmt.Println("解析json失败：", err)
		return nil, err
	}

	return rows, nil

}

// GetNodeNew 新方式获取Node信息  直接返回JSON 目前支持 2.3.0
func GetNodeNew() (Model.Node, error) {

	res, err := Tools.RunCmd(Tools.CliPath, "-o", "json", "node")
	if err != nil {
		fmt.Println("执行cmd命令失败：", err)
		return Model.Node{}, err
	}

	nodeNew := Model.NodeNew{}
	err = json.Unmarshal([]byte(res), &nodeNew)
	if err != nil {
		return Model.Node{}, err
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

	return nodeInfo, nil
}

// GetNode 获取当前节点信息
func GetNode() (Model.Node, error) {

	res, err := Tools.RunCmd(Tools.CliPath, "node")
	if err != nil {
		fmt.Println("执行cmd命令失败：", err)
		return Model.Node{}, err
	}

	// 解析表格并转换为 JSON
	nodeinfo, err := Tools.ParseNodeToModel(res)
	if err != nil {
		fmt.Println("解析返回数据失败：", err)
		return Model.Node{}, err
	}

	return nodeinfo, nil

}

func GetConnectorNew() ([]Model.ConnectorApi, error) {

	res, err := Tools.RunCmd(Tools.CliPath, "-o", "json", "connector")
	if err != nil {
		fmt.Println("执行cmd命令失败：", err)
		return nil, err
	}

	var ConnectorInfo []Model.Connector
	err = json.Unmarshal([]byte(res), &ConnectorInfo)
	if err != nil {
		fmt.Println("解析json失败：", err)
		return nil, err
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

	return connectorApis, nil

}

// GetConnector 获取服务器连接信息
func GetConnector() ([]Model.Connector, error) {

	res, err := Tools.RunCmd(Tools.CliPath, "connector")
	if err != nil {
		fmt.Println("执行cmd命令失败：", err)
		return nil, err
	}

	// 解析表格并转换为 JSON
	ConnectorInfo, err := Tools.ParseConnectorToModel(res)
	if err != nil {
		fmt.Println("解析返回数据失败：", err)
		return nil, err
	}

	return ConnectorInfo, nil

}
