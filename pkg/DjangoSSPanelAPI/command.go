package DjangoSSPanelAPI

import (
	"errors"
	"github.com/crossfw/Air-Universe/pkg/structures"
)

func (sspCtl *DjSspController) GetNodeInfo(closeTLS bool) (err error) {
	return getNodeInfo(sspCtl, closeTLS)
}

func (sspCtl *DjSspController) GetUser() (userList *[]structures.UserInfo, err error) {
	return getUser(sspCtl)
}

func (sspCtl *DjSspController) PostTraffic(trafficData *[]structures.UserTraffic) (err error) {
	return postTraffic(sspCtl, trafficData)
}

func (sspCtl *DjSspController) PostSysLoad(load *structures.SysLoad) (err error) {
	return errors.New("unsupported")
}

func (sspCtl *DjSspController) PostAliveIP(baseCfg *structures.BaseConfig, userIP *[]structures.UserIP) (err error) {
	return errors.New("unsupported")
}

func (sspCtl *DjSspController) GetNowInfo() (nodeInfo *structures.NodeInfo) {
	return sspCtl.NodeInfo
}
