package DjangoSSPanelAPI

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/crossfw/Air-Universe/pkg/structures"
	"net/http"
	"time"
)

type UserTraffic struct {
	UserId          uint32 `json:"user_id"`
	DownloadTraffic int64  `json:"dt"`
	UploadTraffic   int64  `json:"ut"`
}

type syncReq struct {
	UserTraffics *[]UserTraffic `json:"user_traffics"`
}

type SSUserTraffic struct {
	UserId          uint32 `json:"user_id"`
	DownloadTraffic int64  `json:"download_traffic"`
	UploadTraffic   int64  `json:"upload_traffic"`
}

type SSsyncReq struct {
	UserTraffics *[]SSUserTraffic `json:"data"`
}

func postTraffic(node *DjSspController, trafficData *[]structures.UserTraffic) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("unplanned outages when post traffic data")
		}
	}()
	var bodyJson []byte
	if node.NodeInfo.Protocol == "ss" {
		PostTrafficJson := new([]SSUserTraffic)
		for _, userTraffic := range *trafficData {
			var ut SSUserTraffic
			ut.UserId = userTraffic.Id
			ut.UploadTraffic = userTraffic.Up
			ut.DownloadTraffic = userTraffic.Down

			*PostTrafficJson = append(*PostTrafficJson, ut)
		}

		bodyJson, err = json.Marshal(&SSsyncReq{UserTraffics: PostTrafficJson})
		if err != nil {
			return
		}
	} else {
		PostTrafficJson := new([]UserTraffic)
		for _, userTraffic := range *trafficData {
			var ut UserTraffic
			ut.UserId = userTraffic.Id
			ut.UploadTraffic = userTraffic.Up
			ut.DownloadTraffic = userTraffic.Down

			*PostTrafficJson = append(*PostTrafficJson, ut)
		}

		bodyJson, err = json.Marshal(&syncReq{UserTraffics: PostTrafficJson})
		if err != nil {
			return
		}
	}
	client := &http.Client{Timeout: 40 * time.Second}
	defer client.CloseIdleConnections()
	apiURL := ""
	switch node.NodeInfo.Protocol {
	case "vmess":
		apiURL = "api/user_vmess_config"
	case "trojan":
		apiURL = "api/user_trojan_config"
	case "ss":
		apiURL = "api/user_ss_config"
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/%v/?token=%s", node.URL, apiURL, node.NodeInfo.Id, node.Key), bytes.NewBuffer(bodyJson))
	if err != nil {
		return
	}
	// Use json type
	req.Header.Set("Content-Type", "application/json")
	_, err = client.Do(req)
	if err != nil {
		return
	}
	return
}
