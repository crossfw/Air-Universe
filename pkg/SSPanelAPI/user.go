package SSPanelAPI

import (
	"errors"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/crossfw/Air-Universe/pkg/structures"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

func getUser(node *SspController) (userList *[]structures.UserInfo, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("get users from sspanel %s", r))
		}
	}()
	userList = new([]structures.UserInfo)
	user := structures.UserInfo{}
	client := &http.Client{Timeout: 40 * time.Second}
	defer client.CloseIdleConnections()
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/mod_mu/users?key=%s&node_id=%v", node.URL, node.Key, node.NodeInfo.Id), nil)
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
	if rtn.Get("ret").MustInt() != 1 {
		return nil, errors.New(fmt.Sprintf("Server error - %s", rtn.Get("data").MustString()))
	}

	numOfUsers := len(rtn.Get("data").MustArray())
	numOfSkipUsers := 0
	for u := 0; u < numOfUsers; u++ {
		user.Id = uint32(rtn.Get("data").GetIndex(u).Get("id").MustInt())
		user.Uuid = rtn.Get("data").GetIndex(u).Get("uuid").MustString()
		if user.Uuid == "" && node.NodeInfo.Protocol != "ss" {
			numOfSkipUsers++
			continue
		}
		user.Password = rtn.Get("data").GetIndex(u).Get("passwd").MustString()
		user.AlterId = node.NodeInfo.AlterID
		user.Level = 0
		user.InTag = node.NodeInfo.Tag
		user.Tag = fmt.Sprintf("%s-%s", strconv.FormatUint(uint64(user.Id), 10), user.InTag)
		user.Protocol = node.NodeInfo.Protocol
		user.MaxClients = uint32(rtn.Get("data").GetIndex(u).Get("node.NodeInfo_connector").MustInt())

		userSL := uint32(rtn.Get("data").GetIndex(u).Get("node_speedlimit").MustInt())
		// The minimal value decide SpeedLimit
		if userSL > 0 && userSL < node.NodeInfo.SpeedLimit {
			user.SpeedLimit = userSL
		} else if node.NodeInfo.SpeedLimit > 0 {
			user.SpeedLimit = node.NodeInfo.SpeedLimit
		} else {
			user.SpeedLimit = userSL
		}

		//单端口承载用户判定, 请在配置文件中打开为后端下发偏移后端口选项
		isMultiUser := rtn.Get("data").GetIndex(u).Get("is_multi_user").MustInt()
		if isMultiUser > 0 {
			user.SSConfig = true
			if node.NodeInfo.Protocol == "ss" {
				node.NodeInfo.ListenPort = uint32(rtn.Get("data").GetIndex(u).Get("port").MustInt())
				node.NodeInfo.CipherType = rtn.Get("data").GetIndex(u).Get("method").MustString()
			}
		} else {
			user.SSConfig = false
		}

		*userList = append(*userList, user)
	}

	//写入加密方式，为避免承载用户不是第一个，所以单独拉出来循环
	if node.NodeInfo.Protocol == "ss" {
		for u := 0; u < len(*userList); u++ {
			(*userList)[u].CipherType = node.NodeInfo.CipherType
		}
	}
	// 如果有无效用户提示
	if numOfSkipUsers != 0 {
		log.Warnf("There are %v users who haven't valid UUID. Please check your panel.", numOfSkipUsers)
	}
	return userList, nil
}
