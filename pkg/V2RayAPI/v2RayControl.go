package V2RayAPI

import (
	"github.com/crossfw/Air-Universe/pkg/structures"
	"v2ray.com/core/app/proxyman/command"
	statsService "v2ray.com/core/app/stats/command"
)

func V2AddUsers(hsClient *command.HandlerServiceClient, users *[]structures.UserInfo) error {

	for _, u := range *users {
		err := v2AddUser(*hsClient, &u)
		if err != nil {
			return err
		}
	}

	return nil
}

func V2RemoveUsers(hsClient *command.HandlerServiceClient, users *[]structures.UserInfo) error {
	for _, u := range *users {
		err := v2RemoveUser(*hsClient, &u)
		if err != nil {
			return err
		}
	}

	return nil
}

func V2QueryUsersTraffic(StatsClient *statsService.StatsServiceClient, users *[]structures.UserInfo) (usersTraffic *[]structures.UserTraffic, err error) {
	usersTraffic = new([]structures.UserTraffic)
	var ut structures.UserTraffic

	for _, u := range *users {
		ut.Id = u.Id
		ut.Up, err = v2QueryUserTraffic(*StatsClient, u.Tag, "up")
		ut.Down, err = v2QueryUserTraffic(*StatsClient, u.Tag, "down")
		// when a user used this node, post traffic data
		if ut.Up+ut.Down > 0 {
			*usersTraffic = append(*usersTraffic, ut)
		}
		if err != nil {
			return
		}
	}
	return
}
