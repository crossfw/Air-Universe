package SSPanelAPI

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

func postTraffic(node *SspController, trafficData *[]structures.UserTraffic) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("unplanned outages when post traffic data")
		}
	}()
	type trafficType struct {
		Data []structures.UserTraffic `json:"data"`
	}
	var body trafficType
	// build post json
	body.Data = *trafficData
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return
	}
	client := &http.Client{Timeout: 40 * time.Second}
	defer client.CloseIdleConnections()
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/mod_mu/users/traffic?key=%s&node_id=%v", node.URL, node.Key, node.NodeInfo.Id), bytes.NewBuffer(bodyJson))
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
		return errors.New(fmt.Sprintf("Server error - %s", rtn.Get("data").MustString()))
	}

	return
}
