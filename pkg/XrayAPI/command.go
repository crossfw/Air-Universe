package XrayAPI

import (
	xrayCmd "github.com/xtls/xray-core/app/proxyman/command"
	xrayStatsService "github.com/xtls/xray-core/app/stats/command"
	"google.golang.org/grpc"
)

type XrayController struct {
	HsClient *xrayCmd.HandlerServiceClient
	SsClient *xrayStatsService.StatsServiceClient
	CmdConn  *grpc.ClientConn
}
