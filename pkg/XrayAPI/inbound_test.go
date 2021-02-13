package XrayAPI

import (
	"fmt"
	sspApi "github.com/crossfw/Air-Universe/pkg/SSPanelAPI"
	"github.com/crossfw/Air-Universe/pkg/structures"
	"log"
	"testing"
	"time"
)

var (
	baseCfg = &structures.BaseConfig{
		Panel: structures.Panel{
			Type:    "sspanel",
			URL:     "http://10.50.1.1:10080",
			Key:     "crossfw666",
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
	)

	xrayCtl = new(XrayController)
	ssp := new(sspApi.NodeInfo)
	ssp.GetNodeInfo(baseCfg, 0)

	fmt.Println(ssp)
	_ = xrayCtl.Init(baseCfg)
	err := addInbound(*xrayCtl.HsClient, ssp)
	//err := addInboundManual(*xrayCtl.HsClient)
	users, _ := ssp.GetUser(baseCfg)
	log.Println(users)
	err = xrayCtl.AddUsers(users)

	time.Sleep(10 * time.Second)

	err = removeInbound(*xrayCtl.HsClient, ssp)
	_ = xrayCtl.CmdConn.Close()
	log.Println(err)
	if err != nil {
		t.Errorf("Failed")
	}
}
func TestRemoveInbound(t *testing.T) {
	var (
		xrayCtl *XrayController
		ssp     *sspApi.NodeInfo
	)

	ssp = new(sspApi.NodeInfo)
	ssp.GetNodeInfo(baseCfg, 0)
	xrayCtl = new(XrayController)
	_ = xrayCtl.Init(baseCfg)

	err := removeInbound(*xrayCtl.HsClient, ssp)
	//err := removeInboundManual(*xrayCtl.HsClient)
	_ = xrayCtl.CmdConn.Close()
	if err != nil {
		t.Errorf("Failed")
	}
}
