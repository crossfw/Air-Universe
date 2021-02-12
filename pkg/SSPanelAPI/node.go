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

type NodeInfo struct {
	Id                  uint32
	SpeedLimit          uint32 `json:"node_speedlimit"`
	Sort                uint32 `json:"sort"`
	RawInfo             string `json:"server"`
	Url                 string
	Protocol            string
	ListenPort          uint32
	AlertID             uint32
	EnableTLS           bool
	EnableProxyProtocol bool
	TransportMode       string
	Path                string
	Host                string
}

func String2Uint32(s string) (uint32, error) {
	uint64, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint32(uint64), err
}

func GetNodeInfo(baseCfg *structures.BaseConfig, idIndex uint32) (nodeInfo *NodeInfo, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("get users from sspanel failed")
		}
	}()
	nodeInfo = new(NodeInfo)
	client := &http.Client{Timeout: 10 * time.Second}
	defer client.CloseIdleConnections()
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/mod_mu/nodes/%v/info?key=%s", baseCfg.Panel.URL, baseCfg.Panel.NodeIDs[idIndex], baseCfg.Panel.Key), nil)
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
	nodeInfo.Id = idIndex
	nodeInfo.SpeedLimit = uint32(rtn.Get("data").Get("node_speedlimit").MustInt())

	return
}

/*
[url, port, alertId, isTLS, transportMode]   (^|(?<=;))([^;]*)(?=;)
path	(?<=path=).*?(?=\|)|(?<=path=).*
host	(?<=host=).*?(?=\|)|(?<=host=).*
*/
func (node *NodeInfo) parseRawInfo() (err error) {
	reBasicInfos, _ := regexp.Compile("(^|(?<=;))([^;]*)(?=;)", 1)
	rePath, _ := regexp.Compile("(?<=path=).*?(?=\\|)|(?<=path=).*", 1)
	reHost, _ := regexp.Compile("(?<=host=).*?(?=\\|)|(?<=host=).*", 1)
	reInsidePort, _ := regexp.Compile("(?<=inside_port=).*?(?=\\|)|(?<=inside_port=).*", 1)
	reRelay, _ := regexp.Compile("\\|relay", 1)

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
	node.Url = basicInfoArray[0]
	if mInsidePort == nil {
		node.ListenPort, _ = String2Uint32(basicInfoArray[1])
	} else {
		node.ListenPort, _ = String2Uint32(mInsidePort.String())
	}
	node.AlertID, _ = String2Uint32(basicInfoArray[2])

	if basicInfoArray[3] == "tls" {
		node.EnableTLS = true
	} else {
		node.EnableTLS = false
	}

	node.TransportMode = basicInfoArray[4]
	if mPath != nil {
		// First cheater is "\", remove it.
		node.Path = mPath.String()[1:]
	}
	if mRelay != nil {
		node.EnableProxyProtocol = true
	}
	if mHost != nil {
		node.Host = mHost.String()
	}

	//fmt.Println(basicInfoArray)

	return
}
