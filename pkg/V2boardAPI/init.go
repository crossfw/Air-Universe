package V2boardAPI

import (
	"errors"
	"github.com/crossfw/Air-Universe/pkg/structures"
	"strings"
)

type V2bController struct {
	URL      string
	Key      string
	NodeInfo *structures.NodeInfo
}

func (v2bCtl *V2bController) Init(cfg *structures.BaseConfig, idIndex uint32) (err error) {
	v2bCtl.NodeInfo = new(structures.NodeInfo)
	v2bCtl.URL = cfg.Panel.URL
	v2bCtl.Key = cfg.Panel.Key
	v2bCtl.NodeInfo.Id = cfg.Panel.NodeIDs[idIndex]
	v2bCtl.NodeInfo.IdIndex = idIndex
	// 预先写入，如果没有获取到节点配置则使用配置文件的alterID
	v2bCtl.NodeInfo.AlterID = cfg.Proxy.AlterID
	v2bCtl.NodeInfo.Tag = cfg.Proxy.InTags[idIndex]
	v2bCtl.NodeInfo.Cert = cfg.Proxy.Cert
	v2bCtl.NodeInfo.EnableSniffing = cfg.Proxy.EnableSniffing
	// Not force
	if len(cfg.Panel.NodesProxyProtocol) > int(idIndex) {
		v2bCtl.NodeInfo.EnableProxyProtocol = cfg.Panel.NodesProxyProtocol[idIndex]
	} else {
		v2bCtl.NodeInfo.EnableProxyProtocol = false
	}

	switch strings.ToLower(cfg.Panel.NodesType[idIndex]) {
	case "v2ray":
		v2bCtl.NodeInfo.Protocol = "vmess"
	case "vmess":
		v2bCtl.NodeInfo.Protocol = "vmess"
	case "trojan":
		v2bCtl.NodeInfo.Protocol = "trojan"
	case "ss":
		v2bCtl.NodeInfo.Protocol = "ss"
	default:
		err = errors.New("unsupported protocol")
	}

	return err
}
