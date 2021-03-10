package V2RayAPI

import (
	"context"
	"github.com/crossfw/Air-Universe/pkg/structures"
	"github.com/v2fly/v2ray-core/v4/app/proxyman/command"
	"github.com/v2fly/v2ray-core/v4/common/protocol"
	"github.com/v2fly/v2ray-core/v4/common/serial"
	"github.com/v2fly/v2ray-core/v4/proxy/trojan"
	"github.com/v2fly/v2ray-core/v4/proxy/vmess"
)

// level will for control speed limit.
// https://github.com/v2fly/v2ray-core/pull/403

func addVmessUser(c command.HandlerServiceClient, user *structures.UserInfo) error {
	_, err := c.AlterInbound(context.Background(), &command.AlterInboundRequest{
		Tag: user.InTag,
		Operation: serial.ToTypedMessage(&command.AddUserOperation{
			User: &protocol.User{
				Level: user.Level,
				Email: user.Tag,
				Account: serial.ToTypedMessage(&vmess.Account{
					Id:               user.Uuid,
					AlterId:          user.AlertId,
					SecuritySettings: &protocol.SecurityConfig{Type: protocol.SecurityType_AUTO},
				}),
			},
		}),
	})
	if err != nil {
		//log.Printf("failed to call grpc command: %v", err)
		return err
	} else {
		return nil
	}
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

func removeUser(c command.HandlerServiceClient, user *structures.UserInfo) error {
	_, err := c.AlterInbound(context.Background(), &command.AlterInboundRequest{
		Tag: user.InTag,
		Operation: serial.ToTypedMessage(&command.RemoveUserOperation{
			Email: user.Tag,
		}),
	})
	if err != nil {
		//log.Printf("failed to call grpc command: %v", err)
		return err
	} else {
		return nil
	}
}
