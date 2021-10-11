package DjangoSSPanelAPI

import (
	"fmt"
	"github.com/crossfw/Air-Universe/pkg/SysLoad"
	"github.com/crossfw/Air-Universe/pkg/structures"
	"testing"
)

var (
	testCfg = &structures.BaseConfig{
		Panel: structures.Panel{
			Type:      "django-sspanel",
			URL:       "",
			Key:       "",
			NodeIDs:   []uint32{1},
			NodesType: []string{"vmess"},
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

func TestPostTraffic(t *testing.T) {

	var (
		sspCtl      *DjSspController
		trafficData = &[]structures.UserTraffic{
			{
				Up:   5,
				Down: 545454,
				Id:   1,
			},
		}
	)
	sspCtl = new(DjSspController)
	_ = sspCtl.Init(testCfg, 0)
	err := sspCtl.PostTraffic(trafficData)

	if err != nil {
		t.Errorf("Post Failed %s", err)
	}
}

func TestAliveIPost(t *testing.T) {
	var (
		userIPs = []structures.UserIP{
			{
				Id:      1,
				InTag:   "p0",
				AliveIP: []string{"1.1.1.1", "2.2.2.2"},
			},
			{
				Id:      2,
				InTag:   "p0",
				AliveIP: []string{"1.1.1.1", "2.2.2.2"},
			},
		}
		sspCtl *DjSspController
	)

	sspCtl = new(DjSspController)
	_ = sspCtl.Init(testCfg, 0)

	err := sspCtl.PostAliveIP(testCfg, &userIPs)
	t.Log(err)
	if err != nil {
		t.Errorf("Post Failed")
	}
}

func TestPostSysLoad(t *testing.T) {
	var (
		sspCtl *DjSspController
	)
	loaData, err := SysLoad.GetSysLoad()
	sspCtl = new(DjSspController)
	_ = sspCtl.Init(testCfg, 0)
	err = sspCtl.PostSysLoad(loaData)
	if err != nil {
		t.Errorf("Post Failed %s", err)
	}
	fmt.Println(sspCtl.NodeInfo)
}

func TestGetNodeInfo(t *testing.T) {
	var (
		sspCtl *DjSspController
	)
	sspCtl = new(DjSspController)
	_ = sspCtl.Init(testCfg, 0)
	err := sspCtl.GetNodeInfo(false)
	if err != nil {
		t.Errorf("Post Failed %s", err)
	}
	fmt.Println(sspCtl.NodeInfo)
}

func TestGetUsers(t *testing.T) {
	var (
		sspCtl *DjSspController
	)
	sspCtl = new(DjSspController)
	_ = sspCtl.Init(testCfg, 0)
	err := sspCtl.GetNodeInfo(false)
	users, err := sspCtl.GetUser()
	if err != nil {
		t.Errorf("Post Failed %s", err)
	}
	fmt.Println(users)
}
