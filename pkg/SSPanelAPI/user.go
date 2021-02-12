package SSPanelAPI

import (
	"errors"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/crossfw/Air-Universe/pkg/structures"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

func (node *NodeInfo) GetUser(baseCfg *structures.BaseConfig) (userList *[]structures.UserInfo, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("get users from sspanel failed")
		}
	}()
	userList = new([]structures.UserInfo)
	user := structures.UserInfo{}
	client := &http.Client{Timeout: 10 * time.Second}
	defer client.CloseIdleConnections()
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/mod_mu/users?key=%s&node_id=%v", baseCfg.Panel.URL, baseCfg.Panel.Key, node.Id), nil)
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
	if err != nil {
		return nil, err
	}

	for u := 0; u < numOfUsers; u++ {
		user.Id = uint32(rtn.Get("data").GetIndex(u).Get("id").MustInt())
		user.Uuid = rtn.Get("data").GetIndex(u).Get("uuid").MustString()
		user.AlertId = node.AlertID
		user.Level = 0
		user.InTag = baseCfg.Proxy.InTags[node.Id]
		user.Tag = fmt.Sprintf("%s-%s", strconv.FormatUint(uint64(user.Id), 10), user.InTag)
		user.Protocol = node.Protocol
		user.MaxClients = uint32(rtn.Get("data").GetIndex(u).Get("node_connector").MustInt())

		userSL := uint32(rtn.Get("data").GetIndex(u).Get("node_speedlimit").MustInt())
		// The minimal value decide SpeedLimit
		if userSL < node.SpeedLimit {
			user.SpeedLimit = userSL
		} else {
			user.SpeedLimit = node.SpeedLimit
		}

		*userList = append(*userList, user)
	}

	return userList, nil
}
