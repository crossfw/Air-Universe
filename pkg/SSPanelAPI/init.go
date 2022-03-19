package SSPanelAPI

import "github.com/crossfw/Air-Universe/pkg/structures"

type SspController struct {
	URL      string
	Key      string
	NodeInfo *structures.NodeInfo
}

func (sspCtl *SspController) Init(cfg *structures.BaseConfig, idIndex uint32) (err error) {
	sspCtl.NodeInfo = new(structures.NodeInfo)
	sspCtl.URL = cfg.Panel.URL
	sspCtl.Key = cfg.Panel.Key
	sspCtl.NodeInfo.Id = cfg.Panel.NodeIDs[idIndex]
	sspCtl.NodeInfo.IdIndex = idIndex
	// 预先写入，如果没有获取到节点配置则使用配置文件的alterID
	sspCtl.NodeInfo.AlterID = cfg.Proxy.AlterID
	sspCtl.NodeInfo.Tag = cfg.Proxy.InTags[idIndex]
	sspCtl.NodeInfo.Cert = cfg.Proxy.Cert
	sspCtl.NodeInfo.EnableSniffing = cfg.Proxy.EnableSniffing

	return err
}
