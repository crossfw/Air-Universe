package XrayAPI

import (
	"fmt"
	"github.com/crossfw/Air-Universe/pkg/SSPanelAPI"
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
			AlterID:    1,
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
		sspCtl  *SSPanelAPI.SspController
	)

	xrayCtl = new(XrayController)
	sspCtl = new(SSPanelAPI.SspController)
	_ = sspCtl.Init(baseCfg, 0)
	err := sspCtl.GetNodeInfo(false)

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

func TestAddTLSInbound(t *testing.T) {
	var (
		xrayCtl *XrayController
		ssp     = &structures.NodeInfo{
			Id:                  0,
			IdIndex:             0,
			Tag:                 "p0",
			SpeedLimit:          0,
			Sort:                0,
			RawInfo:             "",
			Url:                 "",
			Protocol:            "vmess",
			CipherType:          "",
			ListenPort:          31856,
			AlterID:             0,
			EnableTLS:           true,
			EnableProxyProtocol: false,
			TransportMode:       "tcp",
			Path:                "/a",
			Host:                "xxx.com",
			Cert: structures.Cert{
				CertPath: "cert\\f.crt",
				KeyPath:  "cert\\f.key",
			},
		}
		users = &[]structures.UserInfo{
			{
				Id:         1,
				Level:      0,
				InTag:      "p0",
				Tag:        "1-p0",
				Uuid:       "23ad6b10-8d1a-40f7-8ad0-e3e35cd38297",
				Protocol:   "vmess",
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
		sspCtl  structures.PanelCommand
	)
	xrayCtl = new(XrayController)
	sspCtl = new(SSPanelAPI.SspController)
	_ = sspCtl.Init(baseCfg, 0)
	err := sspCtl.GetNodeInfo(false)

	fmt.Println(sspCtl.GetNowInfo())
	_ = xrayCtl.Init(baseCfg)
	err = xrayCtl.RemoveInbound(sspCtl.GetNowInfo())

	if err != nil {
		t.Errorf("Failed %s", err)
	}
}
