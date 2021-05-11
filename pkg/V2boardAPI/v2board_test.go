package V2boardAPI

import (
	"fmt"
	"github.com/crossfw/Air-Universe/pkg/structures"
	"testing"
)

var (
	testCfg = &structures.BaseConfig{
		Panel: structures.Panel{
			Type:      "v2board",
			URL:       "http://",
			Key:       "",
			NodeIDs:   []uint32{1},
			NodesType: []string{"v2ray"},
		},
		Proxy: structures.Proxy{
			Type:         "xray",
			AlterID:      1,
			AutoGenerate: true,
			InTags:       []string{"p0"},
			APIAddress:   "127.0.0.1",
			APIPort:      10085,
			LogPath:      "./v2.log",
		},
		Sync: structures.Sync{
			Interval:  60,
			FailDelay: 5,
			Timeout:   5,
		},
	}
)

func TestGetNodeInfo(t *testing.T) {
	var (
		v2bCtl *V2bController
	)
	v2bCtl = new(V2bController)
	_ = v2bCtl.Init(testCfg, 0)
	err := v2bCtl.GetNodeInfo(false)
	if err != nil {
		t.Errorf("Post Failed %s", err)
	}
	fmt.Printf("%+v", v2bCtl.NodeInfo)
}

func TestGetUsers(t *testing.T) {
	var (
		v2bCtl *V2bController
	)
	v2bCtl = new(V2bController)
	_ = v2bCtl.Init(testCfg, 0)
	err := v2bCtl.GetNodeInfo(false)
	users, err := v2bCtl.GetUser()
	if err != nil {
		t.Errorf("Post Failed %s", err)
	}
	fmt.Printf("%+v", users)
}

func TestPostTraffic(t *testing.T) {

	var (
		v2bCtl      *V2bController
		trafficData = &[]structures.UserTraffic{
			{
				Up:   5333,
				Down: 545454,
				Id:   1,
			},
		}
	)
	v2bCtl = new(V2bController)
	_ = v2bCtl.Init(testCfg, 0)
	err := v2bCtl.PostTraffic(trafficData)

	if err != nil {
		t.Errorf("Post Failed %s", err)
	}
}
