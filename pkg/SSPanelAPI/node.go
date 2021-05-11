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
[url, port, alterId, isTLS, transportMode]   (.*?)(?=;)
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

func getNodeInfo(node *SspController, closeTLS bool) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("get nodeInfo from sspanel failed %s", r))
		}
	}()

	client := &http.Client{Timeout: 40 * time.Second}
	defer client.CloseIdleConnections()
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/mod_mu/nodes/%v/info?key=%s", node.URL, node.NodeInfo.Id, node.Key), nil)
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
	// ret not equal to 1 means sspanel caused an error or node not fond.
	if rtn.Get("ret").MustInt() != 1 {
		return errors.New(fmt.Sprintf("Server error - %s", rtn.Get("data").MustString()))
	}

	node.NodeInfo.RawInfo = rtn.Get("data").Get("server").MustString()
	node.NodeInfo.Sort = uint32(rtn.Get("data").Get("sort").MustInt())

	node.NodeInfo.SpeedLimit = uint32(rtn.Get("data").Get("node_speedlimit").MustInt())

	switch node.NodeInfo.Sort {
	case 0:
		node.NodeInfo.Protocol = "ss"
		err = parseSSRawInfo(node.NodeInfo)
	case 10:
		node.NodeInfo.Protocol = "ss"
		err = parseSSRawInfo(node.NodeInfo)
		node.NodeInfo.EnableProxyProtocol = true
	case 11:
		node.NodeInfo.Protocol = "vmess"
		err = parseVmessRawInfo(node.NodeInfo, closeTLS)
	case 12:
		node.NodeInfo.Protocol = "vmess"
		err = parseVmessRawInfo(node.NodeInfo, closeTLS)
		// Force Relay
		node.NodeInfo.EnableProxyProtocol = true
	case 14:
		node.NodeInfo.Protocol = "trojan"
		err = parseTrojanRawInfo(node.NodeInfo, closeTLS)
	}

	return nil
}

/*
[url, port, alterId, isTLS, transportMode]   (^|(?<=;))([^;]*)(?=;)
path	(?<=path=).*?(?=\|)|(?<=path=).*
host	(?<=host=).*?(?=\|)|(?<=host=).*
*/
func parseVmessRawInfo(node *structures.NodeInfo, closeTLS bool) (err error) {
	reBasicInfos, _ := regexp.Compile("(^|(?<=;))([^;]*)(?=;)", 1)
	rePath, _ := regexp.Compile("(?<=path=).*?(?=\\||\\?)|(?<=path=).*", 1)
	reHost, _ := regexp.Compile("(?<=host=).*?(?=\\|)|(?<=host=).*", 1)
	reInsidePort, _ := regexp.Compile("(?<=inside_port=).*?(?=\\|)|(?<=inside_port=).*", 1)
	reRelay, _ := regexp.Compile("(?<=relay=).*?(?=\\|)|(?<=relay=)", 1)
	reVless, _ := regexp.Compile("(?<=enable_vless=).*?(?=\\|)|(?<=enable_vless=)", 1)

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
	mVless, _ := reVless.FindStringMatch(node.RawInfo)

	//insidePort := mInsidePort
	if len(basicInfoArray) == 5 {
		node.Url = basicInfoArray[0]
		if mInsidePort == nil {
			node.ListenPort, _ = String2Uint32(basicInfoArray[1])
		} else {
			node.ListenPort, _ = String2Uint32(mInsidePort.String())
		}
		node.AlterID, _ = String2Uint32(basicInfoArray[2])

		node.EnableTLS = false
		for _, transM := range []int{3, 4} {
			switch basicInfoArray[transM] {
			case "tcp":
				node.TransportMode = "tcp"
			case "ws":
				node.TransportMode = "ws"
			case "kcp":
				node.TransportMode = "kcp"
			case "http":
				node.TransportMode = "http"
			case "tls":
				if closeTLS == false {
					node.EnableTLS = true
				}
			}
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
	if mVless != nil {
		node.Protocol = "vless"
	}

	return
}

func parseTrojanRawInfo(node *structures.NodeInfo, closeTLS bool) (err error) {
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
	if closeTLS == false {
		node.EnableTLS = true
	} else {
		node.EnableTLS = false
	}

	return
}

func parseSSRawInfo(node *structures.NodeInfo) (err error) {
	node.TransportMode = "tcp"
	node.EnableTLS = false
	return err
}
