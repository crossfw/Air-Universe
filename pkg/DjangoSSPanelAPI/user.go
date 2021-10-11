package DjangoSSPanelAPI

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

func getUser(node *DjSspController) (userList *[]structures.UserInfo, err error) {
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
		apiURL = "api/user_vmess_config"
	case "trojan":
		apiURL = "api/user_trojan_config"
	case "ss":
		apiURL = "api/user_ss_config"
	}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s/%v/?token=%s", node.URL, apiURL, node.NodeInfo.Id, node.Key), nil)
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

	numOfUsers := len(rtn.Get("configs").MustArray())

	for u := 0; u < numOfUsers; u++ {
		valid := rtn.Get("configs").GetIndex(u).Get("enable").MustBool()
		if valid == false {
			continue
		}
		user.Id = uint32(rtn.Get("configs").GetIndex(u).Get("user_id").MustInt())
		user.Level = 0
		user.InTag = node.NodeInfo.Tag
		user.Tag = fmt.Sprintf("%s-%s", strconv.FormatUint(uint64(user.Id), 10), user.InTag)
		user.Protocol = node.NodeInfo.Protocol
		switch node.NodeInfo.Protocol {
		case "ss":
			user.Password = rtn.Get("configs").GetIndex(u).Get("password").MustString()
			user.CipherType = rtn.Get("configs").GetIndex(u).Get("method").MustString()
			//set ss port
			if u == 0 && node.NodeInfo.Protocol == "ss" {
				node.NodeInfo.ListenPort = uint32(rtn.Get("configs").GetIndex(u).Get("port").MustInt())
			}
		case "trojan":
			user.Uuid = rtn.Get("configs").GetIndex(u).Get("password").MustString()
		case "vmess":
			user.Uuid = rtn.Get("configs").GetIndex(u).Get("uuid").MustString()
			user.AlterId = uint32(rtn.Get("configs").GetIndex(u).Get("alter_id").MustInt())
		}

		*userList = append(*userList, user)
	}
	return
}
