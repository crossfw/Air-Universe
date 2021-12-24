package DjangoSSPanelAPI

import (
	"errors"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/crossfw/Air-Universe/pkg/structures"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func getNodeInfo(node *DjSspController, closeTLS bool) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("get nodeInfo from v2board failed %s", r))
		}
	}()

	client := &http.Client{Timeout: 40 * time.Second}
	defer client.CloseIdleConnections()
	apiURL := ""
	switch node.NodeInfo.Protocol {
	case "vmess":
		apiURL = "api/vmess_server_config"
	case "trojan":
		apiURL = "api/trojan_server_config"
	case "ss":
		node.NodeInfo.TransportMode = "tcp"
		node.NodeInfo.EnableTLS = false
		return err
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
	if rtn.Get("message").MustString() != "" {
		return errors.New(fmt.Sprintf("Server error - %s", rtn.Get("message").MustString()))
	}

	switch node.NodeInfo.Protocol {
	case "vmess":
		err = parseVmessRawInfo(rtn, node.NodeInfo, closeTLS)
	case "trojan":
		err = parseTrojanRawInfo(rtn, node.NodeInfo, closeTLS)
	}
	if err != nil {
		return
	}
	return nil
}

func parseTrojanRawInfo(rtnJson *simplejson.Json, node *structures.NodeInfo, closeTLS bool) (err error) {
	inboundInfo := rtnJson.Get("inbounds").GetIndex(0)
	node.ListenPort = uint32(inboundInfo.Get("port").MustInt())
	node.TransportMode = inboundInfo.Get("streamSettings").Get("network").MustString()
	if closeTLS == false {
		node.EnableTLS = true
	} else {
		node.EnableTLS = false
	}
	return err
}

func parseVmessRawInfo(rtnJson *simplejson.Json, node *structures.NodeInfo, closeTLS bool) (err error) {
	inboundInfo := rtnJson.Get("inbounds").GetIndex(0)
	node.ListenPort = uint32(inboundInfo.Get("port").MustInt())
	node.TransportMode = inboundInfo.Get("streamSettings").Get("network").MustString()

	switch node.TransportMode {
	case "ws":
		wsPath := inboundInfo.Get("streamSettings").Get("wsSettings").Get("path").MustString()
		realPathIndex := strings.Index(wsPath, "?")

		if realPathIndex != 0 {
			node.Path = string([]byte(wsPath)[0:realPathIndex])
		} else {
			node.Path = wsPath
		}
		//node.Host = inboundInfo.Get("streamSettings").Get("wsSettings").Get("headers").Get("Host").MustString()
		if closeTLS == false {
			node.EnableTLS = true
		} else {
			node.EnableTLS = false
		}
	default:
		node.TransportMode = "tcp"
	}

	return err
}
