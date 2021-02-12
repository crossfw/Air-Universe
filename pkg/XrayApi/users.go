package XrayApi

import (
	"fmt"
	"github.com/xtls/xray-core/app/proxyman/command"
	statsService "github.com/xtls/xray-core/app/stats/command"
	"github.com/xtls/xray-core/common/protocol"
	"github.com/xtls/xray-core/common/serial"
	"github.com/xtls/xray-core/proxy/trojan"
	"github.com/xtls/xray-core/proxy/vless"
	"github.com/xtls/xray-core/proxy/vmess"
	"log"
)

func addV2rayVmessUser(client command.HandlerServiceClient, inboundTag string, level uint32, email string, id string, alterID uint32) error {
	_, err := client.AlterInbound(context.Background(), &command.AlterInboundRequest{
		Tag: inboundTag,
		Operation: serial.ToTypedMessage(&command.AddUserOperation{
			User: &protocol.User{
				Level: level,
				Email: email,
				Account: serial.ToTypedMessage(&vmess.Account{
					Id:      id,
					AlterId: alterID,
				}),
			},
		}),
	})
	return err
}

func addV2rayVlessUser(client command.HandlerServiceClient, inboundTag string, level uint32, email string, id string, flow string) error {
	_, err := client.AlterInbound(context.Background(), &command.AlterInboundRequest{
		Tag: inboundTag,
		Operation: serial.ToTypedMessage(&command.AddUserOperation{
			User: &protocol.User{
				Level: level,
				Email: email,
				Account: serial.ToTypedMessage(&vless.Account{
					Id:   id,
					Flow: flow,
				}),
			},
		}),
	})
	return err
}

func addV2rayTrojanUser(client command.HandlerServiceClient, inboundTag string, level uint32, email string, password string) error {
	_, err := client.AlterInbound(context.Background(), &command.AlterInboundRequest{
		Tag: inboundTag,
		Operation: serial.ToTypedMessage(&command.AddUserOperation{
			User: &protocol.User{
				Level: level,
				Email: email,
				Account: serial.ToTypedMessage(&trojan.Account{
					Password: password,
				}),
			},
		}),
	})
	return err
}

func removeV2rayUser(client command.HandlerServiceClient, inboundTag string, email string) error {
	_, err := client.AlterInbound(context.Background(), &command.AlterInboundRequest{
		Tag: inboundTag,
		Operation: serial.ToTypedMessage(&command.RemoveUserOperation{
			Email: email,
		}),
	})
	return err
}

func getV2rayUserTraffic(client statsService.StatsServiceClient, email string, reset bool) {
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
