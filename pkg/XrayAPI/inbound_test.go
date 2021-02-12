package XrayAPI

import (
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
			NodeIDs: []uint32{24},
		},
		Proxy: structures.Proxy{
			Type:       "v2ray",
			AlertID:    1,
			InTags:     []string{"p0"},
			APIAddress: "127.0.0.1",
			APIPort:    10085,
			LogPath:    "./v2.log",
		},
		Sync: structures.Sync{
			Interval:  60,
			FailDelay: 5,
			Timeout:   5,
		},
	}
)

func TestAddInbound(t *testing.T) {
	var xrayCtl *XrayController
	xrayCtl = new(XrayController)
	InitApi(baseCfg, xrayCtl)
	err := addInbound(*xrayCtl.HsClient)
	_ = xrayCtl.CmdConn.Close()
	log.Println(err)
	if err != nil {
		t.Errorf("Failed")
	}
}
