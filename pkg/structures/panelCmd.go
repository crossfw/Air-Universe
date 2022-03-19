package structures

type PanelCommand interface {
	Init(cfg *BaseConfig, idIndex uint32) error
	GetNodeInfo(closeTLS bool) (err error)
	GetUser() (userList *[]UserInfo, err error)
	PostTraffic(trafficData *[]UserTraffic) (err error)
	PostSysLoad(load *SysLoad) (err error)
	PostAliveIP(baseCfg *BaseConfig, userIP *[]UserIP) (err error)
	GetNowInfo() (nodeInfo *NodeInfo)
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
	CipherType          string
	ListenPort          uint32
	AlterID             uint32
	EnableSniffing      bool
	EnableTLS           bool
	EnableProxyProtocol bool
	TransportMode       string
	Path                string
	Host                string
	Cert                Cert
}

type SysLoad struct {
	Uptime uint64
	Load1  float64
	Load5  float64
	Load15 float64
}
