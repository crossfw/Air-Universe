package SSPanelAPI

import (
	"errors"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/crossfw/Air-Universe/pkg/structures"
	regexp "github.com/dlclark/regexp2"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

/*
[url, port, alertId, isTLS, transportMode]   (.*?)(?=;)
path	(?<=path=).*(?=\|)|(?<=path=).*
host	(?<=host=).*(?=\|)|(?<=host=).*
*/

func String2Uint32(s string) (uint32, error) {
	t, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint32(t), err
}

func GetNodeInfo(cfg *structures.BaseConfig, node *structures.NodeInfo, idIndex uint32) (changed bool, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("get users from sspanel failed")
		}
	}()

	var nodeInfo *structures.NodeInfo
	nodeInfo = new(structures.NodeInfo)
	//nodeInfo = new(NodeInfo)
	client := &http.Client{Timeout: 10 * time.Second}
	defer client.CloseIdleConnections()
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/mod_mu/nodes/%v/info?key=%s", cfg.Panel.URL, cfg.Panel.NodeIDs[idIndex], cfg.Panel.Key), nil)
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

	nodeInfo.RawInfo = rtn.Get("data").Get("server").MustString()
	nodeInfo.Sort = uint32(rtn.Get("data").Get("sort").MustInt())
	nodeInfo.Id = cfg.Panel.NodeIDs[idIndex]
	nodeInfo.IdIndex = idIndex
	nodeInfo.SpeedLimit = uint32(rtn.Get("data").Get("node_speedlimit").MustInt())
	if cfg.Proxy.Cert.KeyPath != "" && cfg.Proxy.Cert.CertPath != "" {
		nodeInfo.Cert = cfg.Proxy.Cert
	}
	nodeInfo.Tag = cfg.Proxy.InTags[idIndex]
	switch nodeInfo.Sort {
	case 11:
		nodeInfo.Protocol = "vmess"
		err = parseVmessRawInfo(nodeInfo)
	case 12:
		nodeInfo.Protocol = "vmess"
		err = parseVmessRawInfo(nodeInfo)
		// Force Relay
		nodeInfo.EnableProxyProtocol = true
	case 14:
		nodeInfo.Protocol = "trojan"
		err = parseTrojanRawInfo(nodeInfo)
	}

	if *nodeInfo == *node {
		return false, nil
	} else {
		*node = *nodeInfo
		return true, nil
	}
}

/*
[url, port, alertId, isTLS, transportMode]   (^|(?<=;))([^;]*)(?=;)
path	(?<=path=).*?(?=\|)|(?<=path=).*
host	(?<=host=).*?(?=\|)|(?<=host=).*
*/
func parseVmessRawInfo(node *structures.NodeInfo) (err error) {
	reBasicInfos, _ := regexp.Compile("(^|(?<=;))([^;]*)(?=;)", 1)
	rePath, _ := regexp.Compile("(?<=path=).*?(?=\\|)|(?<=path=).*", 1)
	reHost, _ := regexp.Compile("(?<=host=).*?(?=\\|)|(?<=host=).*", 1)
	reInsidePort, _ := regexp.Compile("(?<=inside_port=).*?(?=\\|)|(?<=inside_port=).*", 1)
	reRelay, _ := regexp.Compile("(?<=relay=).*?(?=\\|)|(?<=relay=)", 1)

	basicInfos, _ := reBasicInfos.FindStringMatch(node.RawInfo)
	var basicInfoArray []string
	for basicInfos != nil {
		basicInfoArray = append(basicInfoArray, basicInfos.String())
		basicInfos, _ = reBasicInfos.FindNextMatch(basicInfos)
	}
	mPath, _ := rePath.FindStringMatch(node.RawInfo)
	mHost, _ := reHost.FindStringMatch(node.RawInfo)
	mRelay, _ := reRelay.FindStringMatch(node.RawInfo)
	mInsidePort, _ := reInsidePort.FindStringMatch(node.RawInfo)
	//insidePort := mInsidePort
	if len(basicInfoArray) == 5 {
		node.Url = basicInfoArray[0]
		if mInsidePort == nil {
			node.ListenPort, _ = String2Uint32(basicInfoArray[1])
		} else {
			node.ListenPort, _ = String2Uint32(mInsidePort.String())
		}
		node.AlertID, _ = String2Uint32(basicInfoArray[2])

		node.TransportMode = basicInfoArray[3]

		if basicInfoArray[4] == "tls" {
			node.EnableTLS = true
		} else {
			node.EnableTLS = false
		}

	} else {
		err = errors.New("panel config missing params")
	}

	if mPath != nil {
		// First cheater is "\", remove it.
		node.Path = mPath.String()
	}
	if mRelay != nil {
		node.EnableProxyProtocol, _ = strconv.ParseBool(mRelay.String())
	} else {
		node.EnableProxyProtocol = false
	}
	if mHost != nil {
		node.Host = mHost.String()
	}

	return
}

func parseTrojanRawInfo(node *structures.NodeInfo) (err error) {
	reUrl, _ := regexp.Compile("(^|(?<=;))([^;]*)(?=;)", 1)
	rePort, _ := regexp.Compile("(?<=port=).*?(?=\\|)|(?<=port=).*", 1)
	reHost, _ := regexp.Compile("(?<=host=).*?(?=\\|)|(?<=host=).*", 1)
	reRelay, _ := regexp.Compile("(?<=relay=).*?(?=\\|)|(?<=relay=)", 1)
	reListenPort, _ := regexp.Compile("(?<=#).*", 1)

	mUrl, _ := reUrl.FindStringMatch(node.RawInfo)
	mPort, _ := rePort.FindStringMatch(node.RawInfo)
	mHost, _ := reHost.FindStringMatch(node.RawInfo)
	mRelay, _ := reRelay.FindStringMatch(node.RawInfo)

	if mUrl != nil {
		// First cheater is "\", remove it.
		node.Url = mUrl.String()
	}

	// Listen port
	if mPort != nil {
		portRaw := mPort.String()
		mListenPort, _ := reListenPort.FindStringMatch(portRaw)
		if mListenPort != nil {
			node.ListenPort, _ = String2Uint32(mListenPort.String())
		} else {
			node.ListenPort, _ = String2Uint32(portRaw)
		}
	}

	if mRelay != nil {
		node.EnableProxyProtocol, _ = strconv.ParseBool(mRelay.String())
	} else {
		node.EnableProxyProtocol = false
	}
	if mHost != nil {
		node.Host = mHost.String()
	}

	node.TransportMode = "tcp"
	node.EnableTLS = true
	return
}
