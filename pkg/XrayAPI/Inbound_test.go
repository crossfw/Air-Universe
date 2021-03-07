package XrayAPI

import (
	"fmt"
	sspApi "github.com/crossfw/Air-Universe/pkg/SSPanelAPI"
	"github.com/crossfw/Air-Universe/pkg/structures"
	"log"
	"testing"
)

var (
	baseCfg = &structures.BaseConfig{
		Panel: structures.Panel{
			Type:    "sspanel",
			URL:     "http://",
			Key:     "",
			NodeIDs: []uint32{37},
		},
		Proxy: structures.Proxy{
			Type:       "v2ray",
			AlertID:    1,
			InTags:     []string{"p0"},
			APIAddress: "127.0.0.1",
			APIPort:    10085,
			LogPath:    "./v2.log",
			//Cert: structures.Cert{
			//	CertPath: "",
			//	KeyPath:  "",
			//},
		},
		Sync: structures.Sync{
			Interval:  60,
			FailDelay: 5,
			Timeout:   5,
		},
	}
)

func TestAutoAddInbound(t *testing.T) {
	var (
		xrayCtl *XrayController
		sspCtl  *sspApi.SspController
	)

	xrayCtl = new(XrayController)
	sspCtl = new(sspApi.SspController)
	_ = sspCtl.Init(baseCfg, 0)
	err := sspCtl.GetNodeInfo()

	fmt.Println(sspCtl.NodeInfo)
	_ = xrayCtl.Init(baseCfg)
	err = addInbound(*xrayCtl.HsClient, sspCtl.NodeInfo)
	users, _ := sspCtl.GetUser()
	log.Println(users)
	err = xrayCtl.AddUsers(users)

	log.Println(err)
	if err != nil {
		t.Errorf("Failed")
	}
}

func TestAddSSInbound(t *testing.T) {
	var (
		xrayCtl *XrayController
		ssp     = &structures.NodeInfo{
			Tag:                 "p0",
			Protocol:            "ss",
			TransportMode:       "tcp",
			EnableTLS:           false,
			EnableProxyProtocol: false,
			ListenPort:          31856,
		}
		users = &[]structures.UserInfo{
			{
				Id:         1,
				Level:      0,
				InTag:      "p0",
				Tag:        "1-p0",
				Protocol:   "ss",
				CipherType: "aes-256-gcm",
				Password:   "1234567",
			},
		}
	)
	xrayCtl = new(XrayController)
	_ = xrayCtl.Init(baseCfg)
	err := addInbound(*xrayCtl.HsClient, ssp)
	if err != nil {
		t.Errorf("Failed%s", err)
	}
	err = xrayCtl.AddUsers(users)
	if err != nil {
		t.Errorf("Failed%s", err)
	}
}

func TestRemoveInbound(t *testing.T) {
	var (
		xrayCtl *XrayController
		ssp     *structures.NodeInfo
	)

	ssp = new(structures.NodeInfo)
	//sspApi.GetNodeInfo(baseCfg, ssp, 0)
	xrayCtl = new(XrayController)
	_ = xrayCtl.Init(baseCfg)

	err := removeInbound(*xrayCtl.HsClient, ssp)
	//err := removeInboundManual(*xrayCtl.HsClient)
	_ = xrayCtl.CmdConn.Close()
	if err != nil {
		t.Errorf("Failed")
	}
}
