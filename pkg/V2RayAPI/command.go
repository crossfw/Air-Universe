package V2RayAPI

import (
	"github.com/crossfw/Air-Universe/pkg/structures"
	"google.golang.org/grpc"
	v2Cmd "v2ray.com/core/app/proxyman/command"
	v2StatsService "v2ray.com/core/app/stats/command"
)

type V2rayCommand interface {
	Init() error
	AddUsers(user *[]structures.UserInfo) error
	RemoveUsers(user *[]structures.UserInfo) error
	QueryUsersTraffic(user *[]structures.UserInfo) (*[]structures.UserTraffic, error)
	AddVmessInbound()
	AddTrojanInbound()
	AddSSInbound()
	RemoveInbound()
}

type V2rayController struct {
	HsClient *v2Cmd.HandlerServiceClient
	SsClient *v2StatsService.StatsServiceClient
	CmdConn  *grpc.ClientConn
}
