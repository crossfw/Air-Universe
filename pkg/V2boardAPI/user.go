package V2boardAPI

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

func getUser(node *V2bController) (userList *[]structures.UserInfo, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("get users from v2board failed %s", r))
		}
	}()

	userList = new([]structures.UserInfo)
	user := structures.UserInfo{}

	client := &http.Client{Timeout: 40 * time.Second}
	defer client.CloseIdleConnections()
	apiURL := ""
	switch node.NodeInfo.Protocol {
	case "vmess":
		apiURL = "api/v1/server/Deepbwork/user"
	case "trojan":
		apiURL = "api/v1/server/TrojanTidalab/user"
	case "ss":
		apiURL = "api/v1/server/ShadowsocksTidalab/user"
	}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s?node_id=%v&token=%s&local_port=1", node.URL, apiURL, node.NodeInfo.Id, node.Key), nil)
	if err != nil {
		return
	}
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

	numOfUsers := len(rtn.Get("data").MustArray())

	for u := 0; u < numOfUsers; u++ {
		user.Id = uint32(rtn.Get("data").GetIndex(u).Get("id").MustInt())
		user.Level = 0
		user.InTag = node.NodeInfo.Tag
		user.Tag = fmt.Sprintf("%s-%s", strconv.FormatUint(uint64(user.Id), 10), user.InTag)
		user.Protocol = node.NodeInfo.Protocol
		switch node.NodeInfo.Protocol {
		case "ss":
			user.Password = rtn.Get("data").GetIndex(u).Get("secret").MustString()
			user.CipherType = rtn.Get("data").GetIndex(u).Get("cipher").MustString()
			//set ss port
			if u == 0 && node.NodeInfo.Protocol == "ss" {
				node.NodeInfo.ListenPort = uint32(rtn.Get("data").GetIndex(u).Get("port").MustInt())
			}
		case "trojan":
			user.Uuid = rtn.Get("data").GetIndex(u).Get("trojan_user").Get("password").MustString()
		case "vmess":
			user.Uuid = rtn.Get("data").GetIndex(u).Get("v2ray_user").Get("uuid").MustString()
			user.AlterId = uint32(rtn.Get("data").GetIndex(u).Get("v2ray_user").Get("alter_id").MustInt())
		}

		*userList = append(*userList, user)
	}
	return
}
