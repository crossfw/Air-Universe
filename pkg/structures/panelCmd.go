package structures

type PanelCmd interface {
	GetNodeInfo(cfg *BaseConfig, idIndex uint32) (changed bool, err error)
	GetUser(cfg *BaseConfig) (userList *[]UserInfo, err error)
	PostTraffic(cfg *BaseConfig, trafficData *[]UserTraffic) (ret int, err error)

	AddInbound(apiClient *interface{}) (ret int, err error)
}

type NodeInfo struct {
	Id                  uint32
	IdIndex             uint32
	Tag                 string
	SpeedLimit          uint32 `json:"node_speedlimit"`
	Sort                uint32 `json:"sort"`
	RawInfo             string `json:"server"`
	Url                 string
	Protocol            string
	ListenPort          uint32
	AlertID             uint32
	EnableTLS           bool
	EnableProxyProtocol bool
	TransportMode       string
	Path                string
	Host                string
	Cert                Cert
}
