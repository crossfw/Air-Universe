package SSPanelAPI

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/crossfw/Air-Universe/pkg/structures"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func (node *NodeInfo) PostTraffic(cfg *structures.BaseConfig, trafficData *[]structures.UserTraffic) (ret int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("post traffic data to sspanel failed")
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
		log.Println("Post body error")
		return 0, err
	}
	log.Println("Traffic data", body)
	client := &http.Client{Timeout: time.Duration(cfg.Sync.Timeout) * time.Second}
	defer client.CloseIdleConnections()
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/mod_mu/users/traffic?key=%s&node_id=%v", cfg.Panel.URL, cfg.Panel.Key, node.Id), bytes.NewBuffer(bodyJson))
	if err != nil {
		return 0, err
	}
	// Use json type
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	rtn, err := simplejson.NewJson(bodyText)
	if err != nil {
		return 0, err
	}

	return rtn.Get("ret").MustInt(), nil
}
