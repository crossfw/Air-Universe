package V2RayAPI

import (
	"github.com/crossfw/Air-Universe/pkg/structures"
	"google.golang.org/grpc"
	v2Cmd "v2ray.com/core/app/proxyman/command"
	v2StatsService "v2ray.com/core/app/stats/command"
)

type V2rayController struct {
	HsClient *v2Cmd.HandlerServiceClient
	SsClient *v2StatsService.StatsServiceClient
	CmdConn  *grpc.ClientConn
}

func (V2rayCtl *V2rayController) AddInbound(node *structures.PanelCmd) (err error) {
	panic("PlaceHolder")
}

func (V2rayCtl *V2rayController) RemoveInbound(node *structures.PanelCmd) (err error) {
	panic("PlaceHolder")
}
