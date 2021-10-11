package main

import (
	"errors"
	"github.com/crossfw/Air-Universe/pkg/DjangoSSPanelAPI"
	"github.com/crossfw/Air-Universe/pkg/SSPanelAPI"
	"github.com/crossfw/Air-Universe/pkg/V2boardAPI"
	"github.com/crossfw/Air-Universe/pkg/XrayAPI"
	"github.com/crossfw/Air-Universe/pkg/structures"
	log "github.com/sirupsen/logrus"
)

func checkCfg() (err error) {
	switch baseCfg.Panel.Type {
	case "sspanel":
		break
	case "v2board":
		break
	case "django-sspanel":
		break
	default:
		err = errors.New("unsupported panel type")
		return
	}

	switch baseCfg.Proxy.Type {
	case "v2ray":
		break
	case "xray":
		break
	default:
		err = errors.New("unsupported proxy type")
		return
	}

	if len(baseCfg.Panel.NodeIDs) != len(baseCfg.Proxy.InTags) {
		err = errors.New("node_ids length isn't equal to in_tags length")
	}

	if len(baseCfg.Panel.NodeIDs) != len(baseCfg.Panel.NodesType) && baseCfg.Panel.Type == "v2board" {
		err = errors.New("node_ids length isn't equal to nodes_type length")
	}

	return
}

func initProxyCore() (apiClient structures.ProxyCommand, err error) {
	switch baseCfg.Proxy.Type {
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
	case "v2board":
		node = new(V2boardAPI.V2bController)
		err = node.Init(baseCfg, idIndex)
		if err != nil {
			log.Error(err)
		}
		return
	case "django-sspanel":
		node = new(DjangoSSPanelAPI.DjSspController)
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
