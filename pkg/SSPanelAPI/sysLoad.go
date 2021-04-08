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

type PostLoad struct {
	Uptime uint64 `json:"uptime"`
	Load   string `json:"load"`
}

func postSysLoad(node *SspController, LoadData *structures.SysLoad) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("unplanned outages when post system loads data to panel")
		}
	}()
	var body PostLoad
	// build post json
	body.Uptime = LoadData.Uptime
	body.Load = fmt.Sprintf("%f %f %f", LoadData.Load1, LoadData.Load5, LoadData.Load15)

	bodyJson, err := json.Marshal(body)
	if err != nil {
		return err
	}
	client := &http.Client{Timeout: 40 * time.Second}
	defer client.CloseIdleConnections()
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/mod_mu/nodes/%v/info?key=%s", node.URL, node.NodeInfo.Id, node.Key), bytes.NewBuffer(bodyJson))
	if err != nil {
		return err
	}
	// Use json type
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	rtn, err := simplejson.NewJson(bodyText)
	if err != nil {
		return err
	}
	if rtn.Get("ret").MustInt() != 1 {
		return errors.New(fmt.Sprintf("Server error - %s", rtn.Get("data").MustString()))
	}

	return
}
