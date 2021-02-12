package structures

type ProxyCommand interface {
	Init(cfg *BaseConfig) error
	AddUsers(user *[]UserInfo) error
	RemoveUsers(user *[]UserInfo) error
	QueryUsersTraffic(user *[]UserInfo) (*[]UserTraffic, error)
	//AddVmessInbound()
	//AddTrojanInbound()
	//AddSSInbound()
	//RemoveInbound()
}
