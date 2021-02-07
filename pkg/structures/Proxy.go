package structures

import (
	"v2ray.com/core/app/proxyman/command"
	statsService "v2ray.com/core/app/stats/command"
)

type V2rayController struct {
	HsClient *command.HandlerServiceClient
	SsClient *statsService.StatsServiceClient
}
