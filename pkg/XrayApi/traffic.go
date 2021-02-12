package XrayApi

import (
	"context"
	"fmt"
	statsService "github.com/xtls/xray-core/app/stats/command"
	"log"
)

func queryUserTraffic(client statsService.StatsServiceClient, email string, reset bool) {
	uplinkStr := fmt.Sprintf("user>>>%s>>>traffic>>>uplink", email)
	downlinkStr := fmt.Sprintf("user>>>%s>>>traffic>>>downlink", email)
	uplinkResp, err := client.GetStats(context.Background(), &statsService.GetStatsRequest{
		Name:   uplinkStr,
		Reset_: reset,
	})

	if err != nil {
		log.Printf("failed to call grpc command: %v", err)
	}

	downlinkResp, err := client.GetStats(context.Background(), &statsService.GetStatsRequest{
		Name:   downlinkStr,
		Reset_: reset,
	})

	if err != nil {
		log.Printf("failed to call grpc command: %v", err)
	}

	log.Printf("name: %v", email)
	log.Printf("uplink: %v", uplinkResp.Stat.Value)
	log.Printf("downlink: %v", downlinkResp.Stat.Value)
	log.Printf("total: %v", uplinkResp.Stat.Value+downlinkResp.Stat.Value)
}
