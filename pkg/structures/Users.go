package structures

type UserInfo struct {
	Id      int
	Uuid    string
	AlertId uint32
	Level   uint32
	InTag   string
}

type UserTraffic struct {
	Id   int   `json:"user_id"`
	Up   int64 `json:"u"`
	Down int64 `json:"d"`
}
