package XrayAPI

import (
	"github.com/crossfw/Air-Universe/pkg/structures"
	xrayCmd "github.com/xtls/xray-core/app/proxyman/command"
	xrayStatsService "github.com/xtls/xray-core/app/stats/command"
	"google.golang.org/grpc"
)

type XrayCommand interface {
	Init() error
	AddUsers(user *[]structures.UserInfo) error
	RemoveUsers(user *[]structures.UserInfo) error
	QueryUserTraffic(user *[]structures.UserInfo) (*[]structures.UserTraffic, error)
	AddVmessInbound()
	AddTrojanInbound()
	AddSSInbound()
	RemoveInbound()
}

type XrayController struct {
	HsClient *xrayCmd.HandlerServiceClient
	SsClient *xrayStatsService.StatsServiceClient
	CmdConn  *grpc.ClientConn
}
