package V2boardAPI

import (
	"errors"
	"github.com/crossfw/Air-Universe/pkg/structures"
)

func (v2bCtl *V2bController) GetNodeInfo(closeTLS bool) (err error) {
	return getNodeInfo(v2bCtl, closeTLS)
}

func (v2bCtl *V2bController) GetUser() (userList *[]structures.UserInfo, err error) {
	return getUser(v2bCtl)
}

func (v2bCtl *V2bController) PostTraffic(trafficData *[]structures.UserTraffic) (err error) {
	return postTraffic(v2bCtl, trafficData)
}

func (v2bCtl *V2bController) PostSysLoad(load *structures.SysLoad) (err error) {
	return errors.New("unsupported method")
}

func (v2bCtl *V2bController) PostAliveIP(baseCfg *structures.BaseConfig, userIP *[]structures.UserIP) (err error) {
	return errors.New("unsupported method")
}

func (v2bCtl *V2bController) GetNowInfo() (nodeInfo *structures.NodeInfo) {
	return v2bCtl.NodeInfo
}
