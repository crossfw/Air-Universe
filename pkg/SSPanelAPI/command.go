package SSPanelAPI

import "github.com/crossfw/Air-Universe/pkg/structures"

func (sspCtl *SspController) GetNodeInfo(closeTLS bool) (err error) {
	return getNodeInfo(sspCtl, closeTLS)
}

func (sspCtl *SspController) GetUser() (userList *[]structures.UserInfo, err error) {
	return getUser(sspCtl)
}

func (sspCtl *SspController) PostTraffic(trafficData *[]structures.UserTraffic) (err error) {
	return postTraffic(sspCtl, trafficData)
}

func (sspCtl *SspController) PostSysLoad(load *structures.SysLoad) (err error) {
	return postSysLoad(sspCtl, load)
}

func (sspCtl *SspController) PostAliveIP(baseCfg *structures.BaseConfig, userIP *[]structures.UserIP) (err error) {
	return postUsersIP(baseCfg, userIP)
}

func (sspCtl *SspController) GetNowInfo() (nodeInfo *structures.NodeInfo) {
	return sspCtl.NodeInfo
}
