package XrayAPI

import (
	"context"
	"github.com/xtls/xray-core/app/proxyman"
	"github.com/xtls/xray-core/app/proxyman/command"
	"github.com/xtls/xray-core/common/net"
	"github.com/xtls/xray-core/common/protocol"
	"github.com/xtls/xray-core/common/serial"
	"github.com/xtls/xray-core/core"
	"github.com/xtls/xray-core/proxy/vmess"
	"github.com/xtls/xray-core/proxy/vmess/inbound"
	"github.com/xtls/xray-core/transport/internet"
	_ "github.com/xtls/xray-core/transport/internet"
	_ "github.com/xtls/xray-core/transport/internet/tcp"
	"github.com/xtls/xray-core/transport/internet/websocket"
)

func addInbound(client command.HandlerServiceClient) error {
	_, err := client.AddInbound(context.Background(), &command.AddInboundRequest{
		Inbound: &core.InboundHandlerConfig{
			Tag: "p6",
			ReceiverSettings: serial.ToTypedMessage(&proxyman.ReceiverConfig{
				PortRange: net.SinglePortRange(23333),
				Listen:    net.NewIPOrDomain(net.AnyIP),
				//SniffingSettings: &proxyman.SniffingConfig{
				//	Enabled: true,
				//	DestinationOverride: []string{"http", "tls"},
				//},
				StreamSettings: &internet.StreamConfig{
					Protocol:     internet.TransportProtocol_WebSocket,
					ProtocolName: "websocket",
					TransportSettings: []*internet.TransportConfig{
						{
							Protocol:     internet.TransportProtocol_WebSocket,
							ProtocolName: "websocket",
							Settings: serial.ToTypedMessage(&websocket.Config{
								Path: "/videos",
								Header: []*websocket.Header{
									{
										Key:   "Host",
										Value: "xray.com",
									},
								},
								AcceptProxyProtocol: false,
							},
							//Settings: serial.ToTypedMessage(&tcp.Config{
							//	AcceptProxyProtocol: false,
							//},
							),
						},
					},
				},
			}),
			ProxySettings: serial.ToTypedMessage(&inbound.Config{
				User: []*protocol.User{
					{
						Email: "123",
						Level: 0,
						Account: serial.ToTypedMessage(&vmess.Account{
							Id:               "25fbb183-ad22-4a06-bd83-1e397b74a4ce",
							AlterId:          32,
							SecuritySettings: &protocol.SecurityConfig{Type: protocol.SecurityType_AUTO},
						}),
					},
				},
				Detour: &inbound.DetourConfig{
					To: "direct",
				},
			}),
		},
	})
	return err
}
