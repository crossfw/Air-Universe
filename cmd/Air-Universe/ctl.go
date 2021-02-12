package main

import (
	"errors"
	"github.com/crossfw/Air-Universe/pkg/SSPanelAPI"
	v2rayApi "github.com/crossfw/Air-Universe/pkg/V2RayAPI"
	"github.com/crossfw/Air-Universe/pkg/XrayApi"
	"github.com/crossfw/Air-Universe/pkg/structures"
	log "github.com/sirupsen/logrus"
	"time"
)

type v2rayController structures.V2rayController
type xrayController structures.XrayController

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

func GetUserSelector(idIndex uint32) (*[]structures.UserInfo, error) {
	switch baseCfg.Panel.Type {
	case "sspanel":
		return sspanelGetUsers(idIndex)
	default:
		err := errors.New("unsupported panel type")
		return nil, err
	}
}

func PostUserSelector(idIndex uint32, traffic *[]structures.UserTraffic) (ret int, err error) {
	switch baseCfg.Panel.Type {
	case "sspanel":
		return sspanelPostTraffic(idIndex, traffic)
	default:
		err := errors.New("unsupported panel type")
		return -1, err
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

func (v2rayCtl v2rayController) v2rayAddUsers(users *[]structures.UserInfo) (err error) {
	err = v2rayApi.V2AddUsers(v2rayCtl.HsClient, users)
	if err != nil {
		log.Warnf("An error caused when adding users to V2ray-core - %s", err)
	}
	return
}

func (v2rayCtl v2rayController) v2rayRemoveUsers(users *[]structures.UserInfo) (err error) {
	err = v2rayApi.V2RemoveUsers(v2rayCtl.HsClient, users)
	if err != nil {
		log.Warnf("An error caused when removing users from V2ray-core - %s", err)
	}
	return
}

func (v2rayCtl v2rayController) v2rayQueryTraffic(users *[]structures.UserInfo) (usersTraffic *[]structures.UserTraffic, err error) {
	usersTraffic, err = v2rayApi.V2QueryUsersTraffic(v2rayCtl.SsClient, users)
	if err != nil {
		log.Warnf("An error caused when querying traffic from V2ray-core - %s", err)
	}
	return
}

func (xrayCtl xrayController) xrayAddUsers(users *[]structures.UserInfo) (err error) {
	err = XrayApi.XrayAddVmessUsers(xrayCtl.HsClient, users)
	if err != nil {
		log.Warnf("An error caused when adding users to V2ray-core - %s", err)
	}
	return
}
