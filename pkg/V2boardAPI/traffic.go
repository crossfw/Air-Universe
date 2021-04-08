package V2boardAPI

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/crossfw/Air-Universe/pkg/structures"
	"io/ioutil"
	"net/http"
	"time"
)

func postTraffic(node *V2bController, trafficData *[]structures.UserTraffic) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("unplanned outages when post traffic data")
		}
	}()

	bodyJson, err := json.Marshal(*trafficData)
	if err != nil {
		return
	}
	client := &http.Client{Timeout: 40 * time.Second}
	defer client.CloseIdleConnections()
	apiURL := ""
	switch node.NodeInfo.Protocol {
	case "vmess":
		apiURL = "api/v1/server/Deepbwork/submit"
	case "trojan":
		apiURL = "api/v1/server/TrojanTidalab/submit"
	case "ss":
		apiURL = "api/v1/server/ShadowsocksTidalab/submit"
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s?node_id=%v&token=%s", node.URL, apiURL, node.NodeInfo.Id, node.Key), bytes.NewBuffer(bodyJson))
	if err != nil {
		return
	}
	// Use json type
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	rtn, err := simplejson.NewJson(bodyText)
	if err != nil {
		return
	}
	if rtn.Get("ret").MustInt() != 1 {
		return errors.New("server error or node not found")
	}

	return
}
