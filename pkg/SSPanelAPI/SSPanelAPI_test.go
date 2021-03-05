package SSPanelAPI

import (
	"fmt"
	"github.com/crossfw/Air-Universe/pkg/SysLoad"
	"github.com/crossfw/Air-Universe/pkg/structures"
	log "github.com/sirupsen/logrus"
	"testing"
)

var (
	testCfg = &structures.BaseConfig{
		Panel: structures.Panel{
			Type:    "sspanel",
			URL:     "",
			Key:     "",
			NodeIDs: []uint32{24},
		},
		Proxy: structures.Proxy{
			Type:         "xray",
			AlertID:      1,
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
	nodeInfo = &structures.NodeInfo{
		Id: 44,
	}
)

func TestPostTraffic(t *testing.T) {

	trafficData := &[]structures.UserTraffic{
		{
			Up:   5,
			Down: 545454,
			Id:   1,
		},
	}

	ret, err := PostTraffic(testCfg, nodeInfo, trafficData)
	if ret != 1 && err != nil {
		log.Println(err)
		t.Errorf("Post Failed")
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
	)

	ret, err := PostUsersIP(testCfg, &userIPs)
	t.Log(err)
	if ret != 1 && err != nil {
		t.Errorf("Post Failed")
	}
}

func TestPostSysLoad(t *testing.T) {
	loaData, err := SysLoad.GetSysLoad()
	if err != nil {
		t.Errorf("Post Failed")
	}
	ret, err := PostSysLoad(testCfg, nodeInfo, loaData)
	if err != nil {
		t.Errorf("Post Failed %s", err)
	}
	fmt.Println(ret)
}
