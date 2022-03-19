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

func (djsspCtl *DjSspController) Init(cfg *structures.BaseConfig, idIndex uint32) (err error) {
	djsspCtl.NodeInfo = new(structures.NodeInfo)
	djsspCtl.URL = cfg.Panel.URL
	djsspCtl.Key = cfg.Panel.Key
	djsspCtl.NodeInfo.Id = cfg.Panel.NodeIDs[idIndex]
	djsspCtl.NodeInfo.IdIndex = idIndex
	// 预先写入，如果没有获取到节点配置则使用配置文件的alterID
	djsspCtl.NodeInfo.AlterID = cfg.Proxy.AlterID
	djsspCtl.NodeInfo.Tag = cfg.Proxy.InTags[idIndex]
	djsspCtl.NodeInfo.Cert = cfg.Proxy.Cert
	djsspCtl.NodeInfo.EnableSniffing = cfg.Proxy.EnableSniffing

	switch strings.ToLower(cfg.Panel.NodesType[idIndex]) {
	case "v2ray":
		djsspCtl.NodeInfo.Protocol = "vmess"
	case "vmess":
		djsspCtl.NodeInfo.Protocol = "vmess"
	case "trojan":
		djsspCtl.NodeInfo.Protocol = "trojan"
	case "ss":
		djsspCtl.NodeInfo.Protocol = "ss"
	default:
		err = errors.New("unsupported protocol")
	}

	return err
}
