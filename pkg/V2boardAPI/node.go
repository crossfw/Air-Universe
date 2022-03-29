package V2boardAPI

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/crossfw/Air-Universe/pkg/structures"
	"io/ioutil"
	"net/http"
	"time"
)

func getNodeInfo(node *V2bController, closeTLS bool) (err error) {
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
		apiURL = "api/v1/server/Deepbwork/config"
	case "trojan":
		apiURL = "api/v1/server/TrojanTidalab/config"
	case "ss":
		node.NodeInfo.TransportMode = "tcp"
		node.NodeInfo.EnableTLS = false
		return err
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
	node.ListenPort = uint32(rtnJson.Get("local_port").MustInt())
	node.Host = rtnJson.Get("ssl").Get("sni").MustString()
	node.TransportMode = "tcp"
	if closeTLS == false {
		node.EnableTLS = true
	} else {
		node.EnableTLS = false
	}
	return err
}

func parseVmessRawInfo(rtnJson *simplejson.Json, node *structures.NodeInfo, closeTLS bool) (err error) {
	// Thanks XrayR
	inboundInfo := simplejson.New()
	if tmpInboundInfo, ok := rtnJson.CheckGet("inbound"); ok {
		inboundInfo = tmpInboundInfo
		// Compatible with v2board 1.5.5-dev
	} else if tmpInboundInfo, ok := rtnJson.CheckGet("inbounds"); ok {
		tmpInboundInfo := tmpInboundInfo.MustArray()
		marshal_byte, _ := json.Marshal(tmpInboundInfo[0].(map[string]interface{}))
		inboundInfo, _ = simplejson.NewJson(marshal_byte)
	} else {
		return fmt.Errorf("Unable to find inbound(s) in the nodeInfo.")
	}
	node.ListenPort = uint32(inboundInfo.Get("port").MustInt())
	node.TransportMode = inboundInfo.Get("streamSettings").Get("network").MustString()

	switch node.TransportMode {
	case "ws":
		node.Path = inboundInfo.Get("streamSettings").Get("wsSettings").Get("path").MustString()
		node.Host = inboundInfo.Get("streamSettings").Get("wsSettings").Get("headers").Get("Host").MustString()
	}

	if inboundInfo.Get("streamSettings").Get("security").MustString() == "tls" && closeTLS == false {
		node.EnableTLS = true
	} else {
		node.EnableTLS = false
	}

	return err
}
