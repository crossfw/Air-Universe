package SSPanelAPI

import "github.com/crossfw/Air-Universe/pkg/structures"

type NodeInfo struct {
	Id                  uint32
	IdIndex             uint32
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
	Cert                structures.Cert
}
