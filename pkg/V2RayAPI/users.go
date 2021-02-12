package V2RayAPI

import (
	"context"
	"github.com/crossfw/Air-Universe/pkg/structures"
	"v2ray.com/core/app/proxyman/command"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/proxy/vmess"
)

// level will for control speed limit.
// https://github.com/v2fly/v2ray-core/pull/403

func v2AddUser(c command.HandlerServiceClient, user *structures.UserInfo) error {
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

func v2RemoveUser(c command.HandlerServiceClient, user *structures.UserInfo) error {
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
