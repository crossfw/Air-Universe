package sspApi

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

func GetUser(baseCfg *structures.BaseConfig, idIndex uint32) (userList *[]structures.UserInfo, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("get users from sspanel failed")
		}
	}()
	userList = new([]structures.UserInfo)
	user := structures.UserInfo{}
	client := &http.Client{Timeout: 10 * time.Second}
	defer client.CloseIdleConnections()
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/mod_mu/users?key=%s&node_id=%v", baseCfg.Panel.URL, baseCfg.Panel.Key, baseCfg.Panel.NodeIDs[idIndex]), nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	rtn, err := simplejson.NewJson(bodyText)
	if err != nil {
		return nil, err
	}

	numOfUsers := len(rtn.Get("data").MustArray())

	for u := 0; u < numOfUsers; u++ {
		user.Id = rtn.Get("data").GetIndex(u).Get("id").MustInt()
		user.Uuid = rtn.Get("data").GetIndex(u).Get("uuid").MustString()
		user.AlertId = baseCfg.Proxy.AlertID
		user.Level = 0
		user.InTag = baseCfg.Proxy.InTags[idIndex]
		*userList = append(*userList, user)
	}

	return userList, nil
}

func PostTraffic(baseCfg *structures.BaseConfig, idIndex uint32, trafficData *[]structures.UserTraffic) (ret int, err error) {
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
	client := &http.Client{Timeout: time.Duration(baseCfg.Sync.Timeout) * time.Second}
	defer client.CloseIdleConnections()
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/mod_mu/users/traffic?key=%s&node_id=%v", baseCfg.Panel.URL, baseCfg.Panel.Key, baseCfg.Panel.NodeIDs[idIndex]), bytes.NewBuffer(bodyJson))
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
