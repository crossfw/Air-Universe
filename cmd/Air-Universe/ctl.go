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
		err = errors.New("node_ids length isn't equal to in_tags length")
	}

	return
}

func initProxyCore() (apiClient structures.ProxyCommand, err error) {
	switch baseCfg.Proxy.Type {
	case "v2ray":
		apiClient = new(v2rayApi.V2rayController)
		for {
			err = apiClient.Init(baseCfg)
			if err != nil {
				log.Error(err)
			} else {
				break
			}
		}
		return
	case "xray":
		apiClient = new(XrayAPI.XrayController)
		for {
			err = apiClient.Init(baseCfg)
			if err != nil {
				log.Error(err)
			} else {
				break
			}
		}
		return
	default:
		err = errors.New("unsupported proxy core")
		return
	}
	return
}

//func initNode() (node structures.PanelCmd, err error) {
//	switch baseCfg.Panel.Type {
//	case "sspanel":
//		node = new(SSPanelAPI.NodeInfo)
//		return
//	default:
//		err := errors.New("unsupported panel type")
//		return nil, err
//	}
//
//}

func getNodeInfo(idIndex uint32) (node *structures.NodeInfo, err error) {
	node = new(structures.NodeInfo)
	switch baseCfg.Panel.Type {
	case "sspanel":
		return sspanelGetNodeInfo(idIndex)
	}
	return nil, errors.New("unsupported panel type ")
}

func getUsers(node *structures.NodeInfo) (*[]structures.UserInfo, error) {
	switch baseCfg.Panel.Type {
	case "sspanel":
		return sspanelGetUsers(node)
	default:
		err := errors.New("unsupported panel type")
		return nil, err
	}
}

func postUsersTraffic(node *structures.NodeInfo, traffic *[]structures.UserTraffic) (ret int, err error) {
	switch baseCfg.Panel.Type {
	case "sspanel":
		return sspanelPostTraffic(node, traffic)
	default:
		err := errors.New("unsupported panel type")
		return -1, err
	}
}

func sspanelGetUsers(node *structures.NodeInfo) (users *[]structures.UserInfo, err error) {
	for {
		users, err = SSPanelAPI.GetUser(baseCfg, node)
		if err != nil {
			log.Warnf("Failed to get users - %s", err)
			time.Sleep(time.Duration(baseCfg.Sync.FailDelay) * time.Second)
		} else {
			break
		}
	}
	return
}

func sspanelPostTraffic(node *structures.NodeInfo, traffic *[]structures.UserTraffic) (ret int, err error) {
	for {
		ret, err = SSPanelAPI.PostTraffic(baseCfg, node, traffic)
		if err != nil {
			log.Warnf("Failed to post traffic - %s", err)
			time.Sleep(time.Duration(baseCfg.Sync.FailDelay) * time.Second)
		} else {
			break
		}
	}
	return
}

func sspanelGetNodeInfo(idIndex uint32) (node *structures.NodeInfo, err error) {
	node = new(structures.NodeInfo)
	for {
		node, err = SSPanelAPI.GetNodeInfo(baseCfg, idIndex)
		if err != nil {
			log.Warnf("Failed to Get NodeInfo - %s", err)
			time.Sleep(time.Duration(baseCfg.Sync.FailDelay) * time.Second)
		} else {
			return
		}
	}
}
