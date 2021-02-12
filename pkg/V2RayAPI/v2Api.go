package V2RayAPI

import (
	"context"
	"fmt"
	"github.com/crossfw/Air-Universe/pkg/structures"
	"google.golang.org/grpc"
	"strconv"

	"v2ray.com/core/app/proxyman/command"
	statsService "v2ray.com/core/app/stats/command"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/proxy/vmess"
)

// level will for control speed limit.
// https://github.com/v2fly/v2ray-core/pull/403

func v2AddUser(c command.HandlerServiceClient, user *structures.UserInfo) error {
	// 区分不同组的用户 Email = id-tag
	userEmail := fmt.Sprintf("%s-%s", strconv.FormatUint(uint64(user.Id), 10), user.InTag)
	_, err := c.AlterInbound(context.Background(), &command.AlterInboundRequest{
		Tag: user.InTag,
		Operation: serial.ToTypedMessage(&command.AddUserOperation{
			User: &protocol.User{
				Level: user.Level,
				Email: userEmail,
				Account: serial.ToTypedMessage(&vmess.Account{
					Id:               user.Uuid,
					AlterId:          user.AlertId,
					SecuritySettings: &protocol.SecurityConfig{Type: protocol.SecurityType_AUTO},
				}),
			},
		}),
	})
	if err != nil {
		//log.Printf("failed to call grpc command: %v", err)
		return err
	} else {
		return nil
	}
}

func v2RemoveUser(c command.HandlerServiceClient, user *structures.UserInfo) error {
	userEmail := fmt.Sprintf("%s-%s", strconv.FormatUint(uint64(user.Id), 10), user.InTag)
	_, err := c.AlterInbound(context.Background(), &command.AlterInboundRequest{
		Tag: user.InTag,
		Operation: serial.ToTypedMessage(&command.RemoveUserOperation{
			Email: userEmail,
		}),
	})
	if err != nil {
		//log.Printf("failed to call grpc command: %v", err)
		return err
	} else {
		return nil
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
		traffic = 0
	}

	return
}

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
		ut.Up, err = v2QueryUserTraffic(*StatsClient, fmt.Sprintf("%s-%s", strconv.FormatUint(uint64(u.Id), 10), u.InTag), "up")
		ut.Down, err = v2QueryUserTraffic(*StatsClient, fmt.Sprintf("%s-%s", strconv.FormatUint(uint64(u.Id), 10), u.InTag), "down")
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

//gRPC 操作一律用指针完成
func V2InitApi(cfg *structures.BaseConfig) (*command.HandlerServiceClient, *statsService.StatsServiceClient, *grpc.ClientConn, error) {
	cmdConn, err := grpc.Dial(fmt.Sprintf("%s:%d", cfg.Proxy.APIAddress, cfg.Proxy.APIPort), grpc.WithInsecure())
	if err != nil {
		return nil, nil, nil, err
	}
	hsClient := command.NewHandlerServiceClient(cmdConn)
	ssClient := statsService.NewStatsServiceClient(cmdConn)

	return &hsClient, &ssClient, cmdConn, nil
}
