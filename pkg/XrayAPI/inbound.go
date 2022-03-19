package XrayAPI

import (
	"context"
	"errors"
	"github.com/crossfw/Air-Universe/pkg/structures"
	"github.com/xtls/xray-core/app/proxyman"
	"github.com/xtls/xray-core/app/proxyman/command"
	"github.com/xtls/xray-core/common/net"
	"github.com/xtls/xray-core/common/protocol/tls/cert"
	"github.com/xtls/xray-core/common/serial"
	"github.com/xtls/xray-core/core"
	"github.com/xtls/xray-core/infra/conf"
	ssInbound "github.com/xtls/xray-core/proxy/shadowsocks"
	trojanInbound "github.com/xtls/xray-core/proxy/trojan"
	vmessInbound "github.com/xtls/xray-core/proxy/vmess/inbound"
	"github.com/xtls/xray-core/transport/internet"
	"github.com/xtls/xray-core/transport/internet/http"
	"github.com/xtls/xray-core/transport/internet/kcp"
	"github.com/xtls/xray-core/transport/internet/tcp"
	"github.com/xtls/xray-core/transport/internet/tls"
	"github.com/xtls/xray-core/transport/internet/websocket"
)

func addInbound(client command.HandlerServiceClient, node *structures.NodeInfo) (err error) {
	var (
		protocolName      string
		transportSettings []*internet.TransportConfig
		securityType      string
		securitySettings  []*serial.TypedMessage
		proxySetting      *serial.TypedMessage
	)

	switch node.Protocol {
	case "vmess":
		proxySetting = serial.ToTypedMessage(&vmessInbound.Config{})
	case "trojan":
		proxySetting = serial.ToTypedMessage(&trojanInbound.ServerConfig{})
	case "ss":
		proxySetting = serial.ToTypedMessage(&ssInbound.ServerConfig{
			Network: []net.Network{2, 3},
		})
	case "vless":
		err = errors.New("unsupported to auto create VLESS inbounds")
		return err
	}

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
	case "kcp":
		protocolName = "mkcp"
		transportSettings = []*internet.TransportConfig{
			{
				ProtocolName: protocolName,
				Settings:     serial.ToTypedMessage(&kcp.Config{}),
			},
		}
	case "http":
		protocolName = "http"
		transportSettings = []*internet.TransportConfig{
			{
				ProtocolName: protocolName,
				Settings: serial.ToTypedMessage(&http.Config{
					Host: []string{node.Host},
					Path: node.Path,
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
		builtCert, err := certConfig.Build()
		if err != nil {
			return err
		}
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

	_, err = client.AddInbound(context.Background(), &command.AddInboundRequest{
		Inbound: &core.InboundHandlerConfig{
			Tag: node.Tag,
			ReceiverSettings: serial.ToTypedMessage(&proxyman.ReceiverConfig{
				//PortList: &net.PortList{
				//	Range: []*net.PortRange{net.SinglePortRange(net.Port(node.ListenPort))},
				//},
				PortList: &net.PortList{Range: []*net.PortRange{net.SinglePortRange(net.Port(node.ListenPort))}},
				Listen:   net.NewIPOrDomain(net.AnyIP),
				SniffingSettings: &proxyman.SniffingConfig{
					Enabled:             node.EnableSniffing,
					DestinationOverride: []string{"http", "tls"},
				},
				StreamSettings: &internet.StreamConfig{
					ProtocolName:      protocolName,
					TransportSettings: transportSettings,
					SecurityType:      securityType,
					SecuritySettings:  securitySettings,
				},
			}),
			ProxySettings: proxySetting,
		},
	})

	return err
}

func removeInbound(client command.HandlerServiceClient, node *structures.NodeInfo) error {
	_, err := client.RemoveInbound(context.Background(), &command.RemoveInboundRequest{
		Tag: node.Tag,
	})
	return err
}
