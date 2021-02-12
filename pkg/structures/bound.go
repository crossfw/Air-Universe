package structures

type Settings struct {
}
type StreamSettings struct {
}
type Sniffing struct {
	Enabled      bool     `json:"enabled"`
	DestOverride []string `json:"destOverride"`
}
type Allocate struct {
	Strategy    string `json:"strategy"`
	Refresh     int    `json:"refresh"`
	Concurrency int    `json:"concurrency"`
}
type Inbound struct {
	Listen         string         `json:"listen"`
	Port           int            `json:"port"`
	Protocol       string         `json:"protocol"`
	Settings       Settings       `json:"settings"`
	StreamSettings StreamSettings `json:"streamSettings"`
	Tag            string         `json:"tag"`
	Sniffing       Sniffing       `json:"sniffing"`
	Allocate       Allocate       `json:"allocate"`
}
