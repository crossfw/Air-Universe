package V2RayAPI

import (
	"errors"
	"github.com/crossfw/Air-Universe/pkg/structures"
	v2BoundsCmd "github.com/v2fly/v2ray-core/v4/app/proxyman/command"
	v2StatsService "github.com/v2fly/v2ray-core/v4/app/stats/command"
	"google.golang.org/grpc"
)

type V2rayController struct {
	HsClient *v2BoundsCmd.HandlerServiceClient
	SsClient *v2StatsService.StatsServiceClient
	CmdConn  *grpc.ClientConn
}

func (v2rayCtl *V2rayController) AddInbound(node *structures.NodeInfo) (err error) {
	return addInbound(*v2rayCtl.HsClient, node)
}

func (v2rayCtl *V2rayController) RemoveInbound(node *structures.NodeInfo) (err error) {
	return removeInbound(*v2rayCtl.HsClient, node)
}

func (v2rayCtl *V2rayController) AddUsers(users *[]structures.UserInfo) (err error) {

	for _, u := range *users {
		switch u.Protocol {
		case "vmess":
			err = addVmessUser(*v2rayCtl.HsClient, &u)
		case "trojan":
			err = addTrojanUser(*v2rayCtl.HsClient, &u)
		case "ss":
			err = errors.New("V2Ray-core not support Shadowsocks protocol")
		}

		if err != nil {
			return err
		}
	}
	return

	return nil
}

func (v2rayCtl *V2rayController) RemoveUsers(users *[]structures.UserInfo) error {
	for _, u := range *users {
		err := removeUser(*v2rayCtl.HsClient, &u)
		if err != nil {
			return err
		}
	}

	return nil
}

func (v2rayCtl *V2rayController) QueryUsersTraffic(users *[]structures.UserInfo) (usersTraffic *[]structures.UserTraffic, err error) {
	usersTraffic = new([]structures.UserTraffic)
	var ut structures.UserTraffic

	for _, u := range *users {
		ut.Id = u.Id
		ut.Up, err = v2QueryUserTraffic(*v2rayCtl.SsClient, u.Tag, "up")
		ut.Down, err = v2QueryUserTraffic(*v2rayCtl.SsClient, u.Tag, "down")
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
