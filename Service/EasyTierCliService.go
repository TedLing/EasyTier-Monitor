package Service

import (
	"EasyWeb/Model"
	"EasyWeb/Tools"
	"encoding/json"
	"fmt"
)

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
