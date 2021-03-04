package SSPanelAPI

import (
	"github.com/crossfw/Air-Universe/pkg/structures"
	log "github.com/sirupsen/logrus"
	"testing"
)

var (
	testCfg = &structures.BaseConfig{
		Panel: structures.Panel{
			Type: "sspanel",
			URL:  "",
			Key:  "",
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
	trafficData = &[]structures.UserTraffic{
		{
			Up:   5,
			Down: 545454,
			Id:   1,
		},
	}
)

func TestPostTraffic(t *testing.T) {
	ret, err := PostTraffic(testCfg, nodeInfo, trafficData)
	if ret != 1 && err != nil {
		log.Println(err)
		t.Errorf("Post Failed")
	}
}
