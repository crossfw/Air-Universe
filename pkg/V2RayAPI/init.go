package V2RayAPI

import (
	"errors"
	"fmt"
	"github.com/crossfw/Air-Universe/pkg/structures"
	v2BoundsCmd "github.com/v2fly/v2ray-core/v4/app/proxyman/command"
	v2StatsService "github.com/v2fly/v2ray-core/v4/app/stats/command"
	"google.golang.org/grpc"
)

func (v2rayCtl *V2rayController) Init(cfg *structures.BaseConfig) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("init apt of V2ray error - %s", r))
		}
	}()
	v2rayCtl.CmdConn, err = grpc.Dial(fmt.Sprintf("%s:%d", cfg.Proxy.APIAddress, cfg.Proxy.APIPort), grpc.WithInsecure())
	if err != nil {
		return
	}
	hsClient := v2BoundsCmd.NewHandlerServiceClient(v2rayCtl.CmdConn)
	ssClient := v2StatsService.NewStatsServiceClient(v2rayCtl.CmdConn)

	v2rayCtl.HsClient = &hsClient
	v2rayCtl.SsClient = &ssClient

	return
}
