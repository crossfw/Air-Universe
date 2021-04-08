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

type userIPData struct {
	IP     string `json:"ip"`
	UserID uint32 `json:"user_id"`
}

type postIPType struct {
	Data []userIPData `json:"data"`
}

func postIP(baseCfg *structures.BaseConfig, idIndex uint32, userIP *postIPType) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("unplanned outages when post alive IP data")
		}
	}()

	bodyJson, err := json.Marshal(userIP)
	if err != nil {
		//errors.New("post body structure is invalid")
		return err
	}
	client := &http.Client{Timeout: 40 * time.Second}
	defer client.CloseIdleConnections()
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/mod_mu/users/aliveip?key=%s&node_id=%v", baseCfg.Panel.URL, baseCfg.Panel.Key, baseCfg.Panel.NodeIDs[idIndex]), bytes.NewBuffer(bodyJson))
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

func postUsersIP(baseCfg *structures.BaseConfig, userIP *[]structures.UserIP) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("unplanned outages when process alive IP data")
		}
	}()

	// 按 inbound tagId 依次推送
	for tagId := 0; tagId < len(baseCfg.Proxy.InTags); tagId++ {
		aliveIPData := postIPType{}
		ipRecord := userIPData{}

		for _, user := range *userIP {
			if user.InTag == baseCfg.Proxy.InTags[tagId] {
				ipRecord.UserID = user.Id
				for _, ip := range user.AliveIP {
					ipRecord.IP = ip
					aliveIPData.Data = append(aliveIPData.Data, ipRecord)
				}
			} else {
				continue
			}
		}

		// 只推送有数据的id
		if len(aliveIPData.Data) != 0 {
			err := postIP(baseCfg, uint32(tagId), &aliveIPData)
			if err != nil {
				return err
			}
		}
	}
	return
}
