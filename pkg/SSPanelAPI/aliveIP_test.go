package SSPanelAPI

import (
	"github.com/crossfw/Air-Universe/pkg/structures"
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

func TestAliveIPost(t *testing.T) {
	var (
		userIPs []structures.UserIP
		userIP  = structures.UserIP{
			Id:      1,
			InTag:   "p0",
			AliveIP: []string{"1.1.1.1", "2.2.2.2"},
		}
	)
	userIPs = append(userIPs, userIP)
	userIP.Id = 4
	userIP.AliveIP = []string{"3.3.3.3", "4.4.4.4"}
	userIPs = append(userIPs, userIP)

	ret, err := PostUsersIP(baseCfg, &userIPs)
	t.Log(err)
	if ret != 1 && err != nil {
		t.Errorf("Post Failed")
	}
}
