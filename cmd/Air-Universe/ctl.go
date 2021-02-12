package main

import (
	"errors"
	"github.com/crossfw/Air-Universe/pkg/SSPanelAPI"
	v2rayApi "github.com/crossfw/Air-Universe/pkg/V2RayAPI"
	"github.com/crossfw/Air-Universe/pkg/XrayAPI"
	"github.com/crossfw/Air-Universe/pkg/structures"
	log "github.com/sirupsen/logrus"
	"time"
)

var (
	v2rayCtl *structures.V2rayController
	xrayCtl  *structures.XrayController
)

func checkCfg() (err error) {
	switch baseCfg.Panel.Type {
	case "sspanel":
		return
	default:
		err = errors.New("unsupported panel type")
	}

	switch baseCfg.Proxy.Type {
	case "v2ray":
		return
	default:
		err = errors.New("unsupported proxy type")
	}

	if len(baseCfg.Panel.NodeIDs) != len(baseCfg.Proxy.InTags) {
		err = errors.New("node_ids isn't equal to in_tags")
	}

	return
}

func initAPI() (err error) {
	switch baseCfg.Proxy.Type {
	case "v2ray":
		for {
			v2rayCtl = new(structures.V2rayController)
			err = v2rayApi.InitApi(baseCfg, v2rayCtl)
			if err != nil {
				log.Error(err)
			} else {
				break
			}
		}
	case "xray":
		for {
			xrayCtl = new(structures.XrayController)
			XrayAPI.InitApi(baseCfg, xrayCtl)
			if err != nil {
				log.Error(err)
			} else {
				break
			}
		}
	default:
		err := errors.New("unsupported proxy core")
		return err
	}
	return
}

func getUser(idIndex uint32) (*[]structures.UserInfo, error) {
	switch baseCfg.Panel.Type {
	case "sspanel":
		return sspanelGetUsers(idIndex)
	default:
		err := errors.New("unsupported panel type")
		return nil, err
	}
}

func postUser(idIndex uint32, traffic *[]structures.UserTraffic) (ret int, err error) {
	switch baseCfg.Panel.Type {
	case "sspanel":
		return sspanelPostTraffic(idIndex, traffic)
	default:
		err := errors.New("unsupported panel type")
		return -1, err
	}
}

func addUser(users *[]structures.UserInfo) (err error) {
	switch baseCfg.Proxy.Type {
	case "v2ray":
		return v2rayAddUsers(users)
	case "xray":
		switch (*users)[0].Protocol {
		case "vmess":
			return xrayAddVmessUsers(users)
		case "trojan":
			return xrayAddTrojanUsers(users)
		}
	default:
		err := errors.New("unsupported proxy core")
		return err
	}
	return
}

func removeUser(users *[]structures.UserInfo) (err error) {
	switch baseCfg.Proxy.Type {
	case "v2ray":
		return v2rayRemoveUsers(users)
	case "xray":
		return xrayRemoveUsers(users)
	default:
		err := errors.New("unsupported proxy core")
		return err
	}
}

func queryTraffic(users *[]structures.UserInfo) (*[]structures.UserTraffic, error) {
	switch baseCfg.Proxy.Type {
	case "v2ray":
		return v2rayQueryTraffic(users)
	case "xray":
		return xrayQueryTraffic(users)
	default:
		err := errors.New("unsupported proxy core")
		return nil, err
	}
}

func sspanelGetUsers(idIndex uint32) (users *[]structures.UserInfo, err error) {
	for {
		users, err = SSPanelAPI.GetUser(baseCfg, idIndex)
		if err != nil {
			log.Warnf("Failed to get users - %s", err)
			time.Sleep(time.Duration(baseCfg.Sync.FailDelay) * time.Second)
		} else {
			break
		}
	}
	return
}

func sspanelPostTraffic(idIndex uint32, traffic *[]structures.UserTraffic) (ret int, err error) {
	for {
		ret, err = SSPanelAPI.PostTraffic(baseCfg, idIndex, traffic)
		if err != nil {
			log.Warnf("Failed to post traffic - %s", err)
			time.Sleep(time.Duration(baseCfg.Sync.FailDelay) * time.Second)
		} else {
			break
		}
	}
	return
}

func v2rayAddUsers(users *[]structures.UserInfo) (err error) {
	err = v2rayApi.V2AddUsers(v2rayCtl.HsClient, users)
	if err != nil {
		log.Warnf("An error caused when adding users to V2ray-core - %s", err)
	}
	return
}

func v2rayRemoveUsers(users *[]structures.UserInfo) (err error) {
	err = v2rayApi.V2RemoveUsers(v2rayCtl.HsClient, users)
	if err != nil {
		log.Warnf("An error caused when removing users from V2ray-core - %s", err)
	}
	return
}

func v2rayQueryTraffic(users *[]structures.UserInfo) (usersTraffic *[]structures.UserTraffic, err error) {
	usersTraffic, err = v2rayApi.V2QueryUsersTraffic(v2rayCtl.SsClient, users)
	if err != nil {
		log.Warnf("An error caused when querying traffic from V2ray-core - %s", err)
	}
	return
}

func xrayAddVmessUsers(users *[]structures.UserInfo) (err error) {
	err = XrayAPI.XrayAddVmessUsers(xrayCtl.HsClient, users)
	if err != nil {
		log.Warnf("An error caused when adding users to Xray-core - %s", err)
	}
	return
}

func xrayAddTrojanUsers(users *[]structures.UserInfo) (err error) {
	err = XrayAPI.XrayAddTrojanUsers(xrayCtl.HsClient, users)
	if err != nil {
		log.Warnf("An error caused when adding users to Xray-core - %s", err)
	}
	return
}

func xrayRemoveUsers(users *[]structures.UserInfo) (err error) {
	err = XrayAPI.XrayRemoveUsers(xrayCtl.HsClient, users)
	if err != nil {
		log.Warnf("An error caused when removing users from Xray-core - %s", err)
	}
	return
}

func xrayQueryTraffic(users *[]structures.UserInfo) (usersTraffic *[]structures.UserTraffic, err error) {
	usersTraffic, err = XrayAPI.XrayQueryUsersTraffic(xrayCtl.SsClient, users)
	if err != nil {
		log.Warnf("An error caused when querying traffic from Xray-core - %s", err)
	}
	return
}
