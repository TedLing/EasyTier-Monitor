package Model

// Connector 定义输出结构
type Connector struct {
	Url    UrlObj `json:"url"`
	Status int    `json:"status"`
}

type UrlObj struct {
	Url string `json:"url"`
}

type ConnectorApi struct {
	Url    string `json:"url"`
	Status string `json:"status"`
}
