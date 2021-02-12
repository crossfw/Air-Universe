package XrayApi

import (
	"fmt"
	"github.com/crossfw/Air-Universe/pkg/structures"
	"github.com/xtls/xray-core/app/proxyman/command"
	statsService "github.com/xtls/xray-core/app/stats/command"
	"google.golang.org/grpc"
)

func InitApi(cfg *structures.BaseConfig) (*command.HandlerServiceClient, *statsService.StatsServiceClient, *grpc.ClientConn, error) {
	cmdConn, err := grpc.Dial(fmt.Sprintf("%s:%d", cfg.Proxy.APIAddress, cfg.Proxy.APIPort), grpc.WithInsecure())
	if err != nil {
		return nil, nil, nil, err
	}
	hsClient := command.NewHandlerServiceClient(cmdConn)
	ssClient := statsService.NewStatsServiceClient(cmdConn)

	return &hsClient, &ssClient, cmdConn, nil
}
