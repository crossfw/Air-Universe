package XrayAPI

import (
	"context"
	"github.com/xtls/xray-core/app/proxyman"
	"github.com/xtls/xray-core/app/proxyman/command"
	"github.com/xtls/xray-core/common/net"
	"github.com/xtls/xray-core/common/serial"
	"github.com/xtls/xray-core/core"
	"github.com/xtls/xray-core/transport/internet"
	"github.com/xtls/xray-core/transport/internet/websocket"
)

func addInbound(client command.HandlerServiceClient) error {
	_, err := client.AddInbound(context.Background(), &command.AddInboundRequest{
		Inbound: &core.InboundHandlerConfig{
			Tag: "p6",
			ReceiverSettings: serial.ToTypedMessage(&proxyman.ReceiverConfig{
				PortRange: net.SinglePortRange(23333),
				Listen: net.NewIPOrDomain(net.AnyIP),
				SniffingSettings: &proxyman.SniffingConfig{
					Enabled: true,
					DestinationOverride: []string{"http", "tls"},
				},
				StreamSettings: &internet.StreamConfig{
					TransportSettings: []*internet.TransportConfig{
						{
							Protocol: internet.TransportProtocol_WebSocket,
							Settings: serial.ToTypedMessage(&websocket.Config{

							}),
						}
,
					}
				}
			}
		},
	})
},
}
