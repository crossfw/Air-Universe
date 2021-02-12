package structures

type Node struct {
	NodeSpeedLimit int    `json:"node_speedlimit"`
	Sort           int    `json:"sort"`
	Server         string `json:"server"`
	NodeID         uint32
}
