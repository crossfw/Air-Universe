package V2RayAPI

import (
	"errors"
	"fmt"
	"github.com/crossfw/Air-Universe/pkg/structures"
	"google.golang.org/grpc"
	"v2ray.com/core/app/proxyman/command"
	statsService "v2ray.com/core/app/stats/command"
)

func InitApi(cfg *structures.BaseConfig, v2gRpc *structures.V2rayController) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("init apt of V2ray error - %s", r))
		}
	}()
	v2gRpc = new(structures.V2rayController)
	v2gRpc.CmdConn, err = grpc.Dial(fmt.Sprintf("%s:%d", cfg.Proxy.APIAddress, cfg.Proxy.APIPort), grpc.WithInsecure())
	if err != nil {
		return
	}
	hsClient := command.NewHandlerServiceClient(v2gRpc.CmdConn)
	ssClient := statsService.NewStatsServiceClient(v2gRpc.CmdConn)

	v2gRpc.HsClient = &hsClient
	v2gRpc.SsClient = &ssClient

	return
}
