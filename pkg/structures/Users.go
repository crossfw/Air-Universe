package structures

import (
	"errors"
	"fmt"
	"reflect"
)

type UserInfo struct {
	Id      uint32
	Uuid    string
	AlterId uint32
	// Level will use for speed limit
	Level uint32
	InTag string
	// Tag = Id + “-” + InTag
	Tag string
	// Protocol Vmess, trojan..
	Protocol   string
	CipherType string
	Password   string
	SpeedLimit uint32
	MaxClients uint32
	// 单端口承载用户标识，true代表该用户为单端口承载用户
	SSConfig bool
}

type UserTraffic struct {
	Id   uint32 `json:"user_id"`
	Up   int64  `json:"u"`
	Down int64  `json:"d"`
}

type UserIP struct {
	Id      uint32
	InTag   string
	AliveIP []string
}

func FindUserDiffer(before, now *[]UserInfo) (remove, add *[]UserInfo, err error) {
	defer func() {
		if r := recover(); r != nil {
			remove = new([]UserInfo)
			add = new([]UserInfo)
			err = errors.New(fmt.Sprintf("model FindUserDiffer cause error - %s", r))
		}
	}()

	remove = new([]UserInfo)
	add = new([]UserInfo)
	// 对于空的对象要处理下，因为会死循环
	if len(*before) == 0 {
		return remove, now, err
	} else if len(*now) == 0 {
		return before, add, err
	}

	n := 0
	b := 0
	//nLastAppear := false
	//bLastAppear := false
	for true {
		if (*before)[b] == (*now)[n] {
			n++
			b++
		} else if (*before)[b].Id < (*now)[n].Id {
			// (*before)[b] has been removed
			*remove = append(*remove, (*before)[b])
			b++
		} else if (*before)[b].Id > (*now)[n].Id {
			// (*now)[n] has been inserted
			*add = append(*add, (*now)[n])
			n++
		} else if (*before)[b].Id == (*now)[n].Id && reflect.DeepEqual((*before)[b], (*now)[n]) == false {
			//user (*before)[b] changed uuid
			*remove = append(*remove, (*before)[b])
			*add = append(*add, (*now)[n])
			n++
			b++
			// Last one will tagged
			//continue
		}
		// any userList finished, break and add remainder users to remove or add
		if n == len(*now) || b == len(*before) {
			break
		}
	}

	// left users will add or remove
	if b != len(*before) {
		for u := b; u < len(*before); u++ {
			*remove = append(*remove, (*before)[u])
		}
	} else if n != len(*now) {
		for u := n; u < len(*now); u++ {
			*add = append(*add, (*now)[u])
		}
	}

	return
}
