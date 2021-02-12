package V2RayAPI

import (
	"errors"
	"fmt"
	"github.com/crossfw/Air-Universe/pkg/structures"
	"google.golang.org/grpc"
	"v2ray.com/core/app/proxyman/command"
	statsService "v2ray.com/core/app/stats/command"
)

func (V2rayCtl *V2rayController) Init(cfg *structures.BaseConfig) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("init apt of V2ray error - %s", r))
		}
	}()
	V2rayCtl.CmdConn, err = grpc.Dial(fmt.Sprintf("%s:%d", cfg.Proxy.APIAddress, cfg.Proxy.APIPort), grpc.WithInsecure())
	if err != nil {
		return
	}
	hsClient := command.NewHandlerServiceClient(V2rayCtl.CmdConn)
	ssClient := statsService.NewStatsServiceClient(V2rayCtl.CmdConn)

	V2rayCtl.HsClient = &hsClient
	V2rayCtl.SsClient = &ssClient

	return
}
