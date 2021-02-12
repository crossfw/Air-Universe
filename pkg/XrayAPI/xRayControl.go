package XrayAPI

import (
	"github.com/crossfw/Air-Universe/pkg/structures"
	"github.com/xtls/xray-core/app/proxyman/command"
	statsService "github.com/xtls/xray-core/app/stats/command"
)

func XrayAddVmessUsers(hsClient *command.HandlerServiceClient, users *[]structures.UserInfo) (err error) {
	for _, u := range *users {
		err := addV2rayVmessUser(*hsClient, &u)
		if err != nil {
			return err
		}
	}
	return
}

func XrayAddTrojanUsers(hsClient *command.HandlerServiceClient, users *[]structures.UserInfo) (err error) {
	for _, u := range *users {
		err := addTrojanUser(*hsClient, &u)
		if err != nil {
			return err
		}
	}
	return
}

func XrayRemoveUsers(hsClient *command.HandlerServiceClient, users *[]structures.UserInfo) (err error) {
	for _, u := range *users {
		err := removeUser(*hsClient, &u)
		if err != nil {
			return err
		}
	}
	return
}

func XrayQueryUsersTraffic(StatsClient *statsService.StatsServiceClient, users *[]structures.UserInfo) (usersTraffic *[]structures.UserTraffic, err error) {
	usersTraffic = new([]structures.UserTraffic)
	var ut structures.UserTraffic

	for _, u := range *users {
		ut.Id = u.Id
		ut.Up, err = queryUserTraffic(*StatsClient, u.Tag, "up")
		ut.Down, err = queryUserTraffic(*StatsClient, u.Tag, "down")
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
