package main

import (
	"errors"
	"github.com/crossfw/Air-Universe/pkg/SSPanelAPI"
	v2rayApi "github.com/crossfw/Air-Universe/pkg/V2RayAPI"
	"github.com/crossfw/Air-Universe/pkg/XrayAPI"
	"github.com/crossfw/Air-Universe/pkg/structures"
	log "github.com/sirupsen/logrus"
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
}

func initPanel(idIndex uint32) (node structures.PanelCommand, err error) {
	switch baseCfg.Panel.Type {
	case "sspanel":
		node = new(SSPanelAPI.SspController)
		err = node.Init(baseCfg, idIndex)
		if err != nil {
			log.Error(err)
		}
		return
	default:
		err := errors.New("unsupported panel type")
		return nil, err
	}

}
