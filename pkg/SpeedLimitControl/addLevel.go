package SpeedLimitControl

import "github.com/crossfw/Air-Universe/pkg/structures"

/*
	向上匹配限速策略，即如果策略限制只有 1Mbps, 10Mbps 用户限速5Mbps 最终会匹配到 10Mbps
*/

func AddLevel(users *[]structures.UserInfo, sl []float32) (err error) {
	var speedIndex uint32
	// 不限速策略，默认使用level0
	for userIndex := 0; userIndex < len(*users); userIndex++ {
		userSpeedLimit := float32((*users)[userIndex].SpeedLimit)

		if userSpeedLimit == 0 || userSpeedLimit > sl[len(sl)-1] {
			(*users)[userIndex].Level = 0
			continue
		}

		for speedIndex = 1; int(speedIndex) < len(sl); speedIndex++ {
			if userSpeedLimit > sl[speedIndex] {
				continue
			} else if userSpeedLimit <= sl[speedIndex] {
				(*users)[userIndex].Level = speedIndex
				break
			}
		}
	}
	return err
}
