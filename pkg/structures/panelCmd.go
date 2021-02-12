package structures

type PanelCmd interface {
	GetNodeInfo(cfg *BaseConfig, idIndex uint32) (err error)
	GetUser(cfg *BaseConfig) (userList *[]UserInfo, err error)
	PostTraffic(cfg *BaseConfig, trafficData *[]UserTraffic) (ret int, err error)
}
