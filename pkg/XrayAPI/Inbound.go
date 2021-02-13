package XrayAPI

import (
	"context"
	"github.com/crossfw/Air-Universe/pkg/SSPanelAPI"
	"github.com/xtls/xray-core/app/proxyman"
	"github.com/xtls/xray-core/app/proxyman/command"
	"github.com/xtls/xray-core/common/net"
	"github.com/xtls/xray-core/common/protocol"
	"github.com/xtls/xray-core/common/protocol/tls/cert"
	"github.com/xtls/xray-core/common/serial"
	"github.com/xtls/xray-core/core"
	"github.com/xtls/xray-core/infra/conf"
	"github.com/xtls/xray-core/proxy/vmess"
	"github.com/xtls/xray-core/proxy/vmess/inbound"
	"github.com/xtls/xray-core/transport/internet"
	_ "github.com/xtls/xray-core/transport/internet"
	"github.com/xtls/xray-core/transport/internet/tcp"
	_ "github.com/xtls/xray-core/transport/internet/tcp"
	"github.com/xtls/xray-core/transport/internet/tls"
	"github.com/xtls/xray-core/transport/internet/websocket"
)

func addInbound(client command.HandlerServiceClient, node *SSPanelAPI.NodeInfo) error {
	var (
		protocolName      string
		transportSettings []*internet.TransportConfig
		securityType      string
		securitySettings  []*serial.TypedMessage
	)
	switch node.TransportMode {
	case "ws":
		protocolName = "websocket"
		if node.Path == "" {
			node.Path = "/"
		}
		header := []*websocket.Header{
			{
				Key:   "Host",
				Value: node.Host,
			},
		}

		transportSettings = []*internet.TransportConfig{
			{
				ProtocolName: protocolName,
				Settings: serial.ToTypedMessage(&websocket.Config{
					Path:                node.Path,
					Header:              header,
					AcceptProxyProtocol: node.EnableProxyProtocol,
				},
				),
			},
		}

	case "tcp":
		protocolName = "tcp"
		transportSettings = []*internet.TransportConfig{
			{
				ProtocolName: protocolName,
				Settings: serial.ToTypedMessage(&tcp.Config{
					AcceptProxyProtocol: node.EnableProxyProtocol,
				}),
			},
		}
	}

	if node.EnableTLS == true && node.Cert.CertPath != "" && node.Cert.KeyPath != "" {
		// Use custom cert file
		certConfig := &conf.TLSCertConfig{
			CertFile: node.Cert.CertPath,
			KeyFile:  node.Cert.KeyPath,
		}
		builtCert, _ := certConfig.Build()
		securityType = serial.GetMessageType(&tls.Config{})
		securitySettings = []*serial.TypedMessage{
			serial.ToTypedMessage(&tls.Config{
				Certificate: []*tls.Certificate{builtCert},
			}),
		}
	} else if node.EnableTLS == true {
		// Auto build cert
		securityType = serial.GetMessageType(&tls.Config{})
		securitySettings = []*serial.TypedMessage{
			serial.ToTypedMessage(&tls.Config{
				Certificate: []*tls.Certificate{tls.ParseCertificate(cert.MustGenerate(nil))},
			}),
		}
	} else {
		// Disable TLS
		securityType = ""
		securitySettings = nil
	}

	_, err := client.AddInbound(context.Background(), &command.AddInboundRequest{
		Inbound: &core.InboundHandlerConfig{
			//Tag: fmt.Sprintf("p%v", node.IdIndex),
			Tag: "p0",
			ReceiverSettings: serial.ToTypedMessage(&proxyman.ReceiverConfig{
				PortRange: net.SinglePortRange(net.Port(node.ListenPort)),
				Listen:    net.NewIPOrDomain(net.AnyIP),
				SniffingSettings: &proxyman.SniffingConfig{
					Enabled:             true,
					DestinationOverride: []string{"http", "tls"},
				},
				StreamSettings: &internet.StreamConfig{
					ProtocolName:      protocolName,
					TransportSettings: transportSettings,
					SecurityType:      securityType,
					SecuritySettings:  securitySettings,
				},
			}),
			ProxySettings: serial.ToTypedMessage(&inbound.Config{
				Detour: &inbound.DetourConfig{
					To: "direct",
				},
			}),
		},
	})

	return err
}

func addInboundManual(client command.HandlerServiceClient) error {
	_, err := client.AddInbound(context.Background(), &command.AddInboundRequest{
		Inbound: &core.InboundHandlerConfig{
			Tag: "p0",
			ReceiverSettings: serial.ToTypedMessage(&proxyman.ReceiverConfig{
				PortRange: net.SinglePortRange(23333),
				Listen:    net.NewIPOrDomain(net.AnyIP),
				SniffingSettings: &proxyman.SniffingConfig{
					Enabled:             true,
					DestinationOverride: []string{"http", "tls"},
				},
				StreamSettings: &internet.StreamConfig{
					ProtocolName: "websocket",
					TransportSettings: []*internet.TransportConfig{
						{
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
					SecurityType: serial.GetMessageType(&tls.Config{}),
					SecuritySettings: []*serial.TypedMessage{
						serial.ToTypedMessage(&tls.Config{
							Certificate: []*tls.Certificate{tls.ParseCertificate(cert.MustGenerate(nil))},
						}),
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
