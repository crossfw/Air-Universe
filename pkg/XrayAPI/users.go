package XrayAPI

import (
	"context"
	"github.com/crossfw/Air-Universe/pkg/structures"
	"github.com/xtls/xray-core/app/proxyman/command"
	"github.com/xtls/xray-core/common/protocol"
	"github.com/xtls/xray-core/common/serial"
	"github.com/xtls/xray-core/proxy/shadowsocks"
	"github.com/xtls/xray-core/proxy/trojan"
	"github.com/xtls/xray-core/proxy/vless"
	"github.com/xtls/xray-core/proxy/vmess"
)

func addVmessUser(client command.HandlerServiceClient, user *structures.UserInfo) error {
	_, err := client.AlterInbound(context.Background(), &command.AlterInboundRequest{
		Tag: user.InTag,
		Operation: serial.ToTypedMessage(&command.AddUserOperation{
			User: &protocol.User{
				Level: user.Level,
				Email: user.Tag,
				Account: serial.ToTypedMessage(&vmess.Account{
					Id:      user.Uuid,
					AlterId: user.AlterId,
				}),
			},
		}),
	})
	return err
}

func addVlessUser(client command.HandlerServiceClient, user *structures.UserInfo) error {
	_, err := client.AlterInbound(context.Background(), &command.AlterInboundRequest{
		Tag: user.InTag,
		Operation: serial.ToTypedMessage(&command.AddUserOperation{
			User: &protocol.User{
				Level: user.Level,
				Email: user.Tag,
				Account: serial.ToTypedMessage(&vless.Account{
					Id:   user.Uuid,
					Flow: "xtls-rprx-direct",
				}),
			},
		}),
	})
	return err
}

func addSSUser(client command.HandlerServiceClient, user *structures.UserInfo) error {
	var ssCipherType shadowsocks.CipherType
	switch user.CipherType {
	case "aes-128-gcm":
		ssCipherType = shadowsocks.CipherType_AES_128_GCM
	case "aes-256-gcm":
		ssCipherType = shadowsocks.CipherType_AES_256_GCM
	case "chacha20-ietf-poly1305":
		ssCipherType = shadowsocks.CipherType_CHACHA20_POLY1305
	}

	_, err := client.AlterInbound(context.Background(), &command.AlterInboundRequest{
		Tag: user.InTag,
		Operation: serial.ToTypedMessage(&command.AddUserOperation{
			User: &protocol.User{
				Level: user.Level,
				Email: user.Tag,
				Account: serial.ToTypedMessage(&shadowsocks.Account{
					Password:   user.Password,
					CipherType: ssCipherType,
				}),
			},
		}),
	})
	return err
}

func addTrojanUser(client command.HandlerServiceClient, user *structures.UserInfo) error {
	_, err := client.AlterInbound(context.Background(), &command.AlterInboundRequest{
		Tag: user.InTag,
		Operation: serial.ToTypedMessage(&command.AddUserOperation{
			User: &protocol.User{
				Level: user.Level,
				Email: user.Tag,
				Account: serial.ToTypedMessage(&trojan.Account{
					Password: user.Uuid,
				}),
			},
		}),
	})
	return err
}

func removeUser(client command.HandlerServiceClient, user *structures.UserInfo) error {
	_, err := client.AlterInbound(context.Background(), &command.AlterInboundRequest{
		Tag: user.InTag,
		Operation: serial.ToTypedMessage(&command.RemoveUserOperation{
			Email: user.Tag,
		}),
	})
	return err
}
