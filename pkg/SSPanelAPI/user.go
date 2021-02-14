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

func GetUser(cfg *structures.BaseConfig, node *structures.NodeInfo) (userList *[]structures.UserInfo, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("get users from sspanel failed")
		}
	}()
	userList = new([]structures.UserInfo)
	user := structures.UserInfo{}
	client := &http.Client{Timeout: 10 * time.Second}
	defer client.CloseIdleConnections()
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/mod_mu/users?key=%s&node_id=%v", cfg.Panel.URL, cfg.Panel.Key, node.Id), nil)
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
		user.Password = rtn.Get("data").GetIndex(u).Get("passwd").MustString()
		user.AlertId = node.AlertID
		user.Level = 0
		user.InTag = cfg.Proxy.InTags[node.IdIndex]
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

		//单端口承载用户判定, 请在配置文件中打开为后端下发偏移后端口选项
		isMultiUser := rtn.Get("data").GetIndex(u).Get("is_multi_user").MustInt()
		if isMultiUser > 0 {
			user.SSConfig = true
			if node.Protocol == "ss" {
				node.ListenPort = uint32(rtn.Get("data").GetIndex(u).Get("port").MustInt())
				node.CipherType = rtn.Get("data").GetIndex(u).Get("method").MustString()
			}
		} else {
			user.SSConfig = false
		}

		*userList = append(*userList, user)
	}

	//写入加密方式，为避免承载用户不是第一个，所以单独拉出来循环
	if node.Protocol == "ss" {
		for _, u := range *userList {
			u.CipherType = node.CipherType
		}
	}

	return userList, nil
}
