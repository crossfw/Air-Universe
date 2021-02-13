package structures

type PanelCmd interface {
	GetNodeInfo(cfg *BaseConfig, idIndex uint32) (changed bool, err error)
	GetUser(cfg *BaseConfig) (userList *[]UserInfo, err error)
	PostTraffic(cfg *BaseConfig, trafficData *[]UserTraffic) (ret int, err error)

	AddInbound(apiClient *interface{}) (ret int, err error)
}
