package v2rayApi

import (
	"context"
	"errors"
	"fmt"
	"github.com/crossfw/Air-Universe/pkg/structures"
	"log"
	"strconv"

	"v2ray.com/core/app/proxyman/command"
	statsService "v2ray.com/core/app/stats/command"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/proxy/vmess"
)

type userInfo structures.UserInfo
type userTraffic structures.UserTraffic

// level will for control speed limit.
// https://github.com/v2fly/v2ray-core/pull/403

func v2AddUser(c command.HandlerServiceClient, user *userInfo) {
	_, err := c.AlterInbound(context.Background(), &command.AlterInboundRequest{
		Tag: user.InTag,
		Operation: serial.ToTypedMessage(&command.AddUserOperation{
			User: &protocol.User{
				Level: user.Level,
				Email: strconv.Itoa(user.Id),
				Account: serial.ToTypedMessage(&vmess.Account{
					Id:               user.Uuid,
					AlterId:          user.AlertId,
					SecuritySettings: &protocol.SecurityConfig{Type: protocol.SecurityType_AUTO},
				}),
			},
		}),
	})
	if err != nil {
		log.Printf("failed to call grpc command: %v", err)
	} else {
		//log.Printf("ok: %v", resp)
	}
}

func v2RemoveUser(c command.HandlerServiceClient, user *userInfo) {
	_, err := c.AlterInbound(context.Background(), &command.AlterInboundRequest{
		Tag: user.InTag,
		Operation: serial.ToTypedMessage(&command.RemoveUserOperation{
			Email: strconv.Itoa(user.Id),
		}),
	})
	if err != nil {
		log.Printf("failed to call grpc command: %v", err)
	} else {
		//log.Printf("ok: %v", resp)
	}
}

func v2QueryUserTraffic(c statsService.StatsServiceClient, userId, direction string) (traffic int64, err error) {
	// var userTraffic *string
	traffic = 0
	ptn := fmt.Sprintf("user>>>%s>>>traffic>>>%slink", userId, direction)
	resp, err := c.QueryStats(context.Background(), &statsService.QueryStatsRequest{
		Pattern: ptn,
		Reset_:  true, // reset traffic data everytime
	})
	if err != nil {
		return
	}
	// Get traffic data
	stat := resp.GetStat()
	// look at v2ray.com/core/app/stats/command stat structure
	if len(stat) != 0 {
		traffic = stat[0].Value
	} else {
		err = errors.New("query traffic error - wrong data")
	}
	return
}

func AddUsers(hsClient *command.HandlerServiceClient, users *[]userInfo) error {
	log.Println("Adding users:", users)
	//cmdConn, err := grpc.Dial(fmt.Sprintf("%s:%d", v2ApiAddr, v2ApiPort), grpc.WithInsecure())
	//if err != nil {
	//	return err
	//}
	//hsClient := command.NewHandlerServiceClient(cmdConn)
	for _, u := range *users {
		v2AddUser(*hsClient, &u)
	}

	return nil
}

func RemoveUsers(hsClient *command.HandlerServiceClient, users *[]userInfo) error {
	log.Println("Adding users:", users)
	//cmdConn, err := grpc.Dial(fmt.Sprintf("%s:%d", v2ApiAddr, v2ApiPort), grpc.WithInsecure())
	//if err != nil {
	//	return err
	//}
	//hsClient := command.NewHandlerServiceClient(cmdConn)
	for _, u := range *users {
		v2RemoveUser(*hsClient, &u)
	}

	return nil
}

func QueryUsersTraffic(StatsClient *statsService.StatsServiceClient, users *[]userInfo) (usersTraffic *[]userTraffic, err error) {
	usersTraffic = new([]userTraffic)
	var ut userTraffic
	//cmdConn, err := grpc.Dial(fmt.Sprintf("%s:%d", v2ApiAddr, v2ApiPort), grpc.WithInsecure())
	//if err != nil {
	//	return nil, err
	//}
	//hsClientTraffic := statsService.NewStatsServiceClient(cmdConn)
	for _, u := range *users {
		ut.Id = u.Id
		ut.Up, _ = v2QueryUserTraffic(*StatsClient, strconv.Itoa(u.Id), "up")
		ut.Down, _ = v2QueryUserTraffic(*StatsClient, strconv.Itoa(u.Id), "down")
		// when a user used this node, post traffic data
		if ut.Up+ut.Down > 0 {
			*usersTraffic = append(*usersTraffic, ut)
		}
	}

	return usersTraffic, nil
}
