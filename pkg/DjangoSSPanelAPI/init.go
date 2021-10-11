package DjangoSSPanelAPI

import (
	"errors"
	"github.com/crossfw/Air-Universe/pkg/structures"
	"strings"
)

type DjSspController struct {
	URL      string
	Key      string
	NodeInfo *structures.NodeInfo
}

func (sspCtl *DjSspController) Init(cfg *structures.BaseConfig, idIndex uint32) (err error) {
	sspCtl.NodeInfo = new(structures.NodeInfo)
	sspCtl.URL = cfg.Panel.URL
	sspCtl.Key = cfg.Panel.Key
	sspCtl.NodeInfo.Id = cfg.Panel.NodeIDs[idIndex]
	sspCtl.NodeInfo.IdIndex = idIndex
	// 预先写入，如果没有获取到节点配置则使用配置文件的alterID
	sspCtl.NodeInfo.AlterID = cfg.Proxy.AlterID
	sspCtl.NodeInfo.Tag = cfg.Proxy.InTags[idIndex]
	sspCtl.NodeInfo.Cert = cfg.Proxy.Cert

	switch strings.ToLower(cfg.Panel.NodesType[idIndex]) {
	case "v2ray":
		sspCtl.NodeInfo.Protocol = "vmess"
	case "vmess":
		sspCtl.NodeInfo.Protocol = "vmess"
	case "trojan":
		sspCtl.NodeInfo.Protocol = "trojan"
	case "ss":
		sspCtl.NodeInfo.Protocol = "ss"
	default:
		err = errors.New("unsupported protocol")
	}

	return err
}
