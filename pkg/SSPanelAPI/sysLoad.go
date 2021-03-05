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

type PostLoad struct {
	Uptime uint64 `json:"uptime"`
	Load   string `json:"load"`
}

func PostSysLoad(cfg *structures.BaseConfig, node *structures.NodeInfo, LoadData *structures.SysLoad) (ret int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("post system loads data to sspanel failed")
		}
	}()
	var body PostLoad
	// build post json
	body.Uptime = LoadData.Uptime
	body.Load = fmt.Sprintf("%f %f %f", LoadData.Load1, LoadData.Load5, LoadData.Load15)

	bodyJson, err := json.Marshal(body)
	if err != nil {
		log.Println("Post body error")
		return 0, err
	}
	client := &http.Client{Timeout: time.Duration(cfg.Sync.Timeout) * time.Second}
	defer client.CloseIdleConnections()
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/mod_mu/nodes/%v/info?key=%s", cfg.Panel.URL, node.Id, cfg.Panel.Key), bytes.NewBuffer(bodyJson))
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
