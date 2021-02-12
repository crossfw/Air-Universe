package structures

type ssp struct {
	Ret  int  `json:"ret"`
	Data Data `json:"data"`
}
type Data struct {
	NodeGroup      int    `json:"node_group"`
	NodeClass      int    `json:"node_class"`
	NodeSpeedlimit int    `json:"node_speedlimit"`
	TrafficRate    int    `json:"traffic_rate"`
	MuOnly         int    `json:"mu_only"`
	Sort           int    `json:"sort"`
	Server         string `json:"server"`
	Type           string `json:"type"`
}
