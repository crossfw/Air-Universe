package structures

import sspApi "github.com/crossfw/Air-Universe/pkg/SSPanelAPI"

type ProxyCommand interface {
	Init(cfg *BaseConfig) error
	AddUsers(user *[]UserInfo) error
	RemoveUsers(user *[]UserInfo) error
	QueryUsersTraffic(user *[]UserInfo) (*[]UserTraffic, error)
	AddInbound(node *sspApi.NodeInfo) (err error)
	//AddTrojanInbound()
	//AddSSInbound()
	RemoveInbound(node *sspApi.NodeInfo) (err error)
}
