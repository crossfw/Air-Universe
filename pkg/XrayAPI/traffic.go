package XrayAPI

import (
	"context"
	"fmt"
	statsService "github.com/xtls/xray-core/app/stats/command"
)

func queryUserTraffic(c statsService.StatsServiceClient, userId, direction string) (traffic int64, err error) {
	// var userTraffic *string
	traffic = 0
	ptn := fmt.Sprintf("user>>>%s>>>traffic>>>%slink", userId, direction)
	resp, err := c.QueryStats(context.Background(), &statsService.QueryStatsRequest{
		Pattern: ptn,
		Reset_:  true, // reset traffic data everytime
	})
	if err != nil {
		return
	}
	// Get traffic data
	stat := resp.GetStat()

	if len(stat) != 0 {
		traffic = stat[0].Value
	} else {
		traffic = 0
	}
	return
}
