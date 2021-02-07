package main

import (
	"errors"
	"github.com/crossfw/Air-Universe/pkg/SSPanelApi"
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

	return
}

func GetUserSelector() (*[]structures.UserInfo, error) {
	switch baseCfg.Panel.Type {
	case "sspanel":
		return sspanelGetUsers()
	default:
		err := errors.New("unsupported panel type")
		return nil, err
	}
}

func PostUserSelector(traffic *[]structures.UserTraffic) (ret int, err error) {
	switch baseCfg.Panel.Type {
	case "sspanel":
		return sspanelPostTraffic(traffic)
	default:
		err := errors.New("unsupported panel type")
		return -1, err
	}
}

func sspanelGetUsers() (users *[]structures.UserInfo, err error) {
	for {
		users, err = sspApi.GetUser(baseCfg)
		if err != nil {
			log.Warnf("Failed to get users - %s", err)
			time.Sleep(time.Duration(baseCfg.Sync.FailDelay) * time.Second)
		} else {
			break
		}
	}
	return
}

func sspanelPostTraffic(traffic *[]structures.UserTraffic) (ret int, err error) {
	for {
		ret, err = sspApi.PostTraffic(baseCfg, traffic)
		if err != nil {
			log.Warnf("Failed to post traffic - %s", err)
			time.Sleep(time.Duration(baseCfg.Sync.FailDelay) * time.Second)
		} else {
			break
		}
	}
	return
}
