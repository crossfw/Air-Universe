package XrayApi

import (
	"github.com/crossfw/Air-Universe/pkg/structures"
	"github.com/xtls/xray-core/app/proxyman/command"
)

func XrayAddVmessUsers(hsClient *command.HandlerServiceClient, users *[]structures.UserInfo) (err error) {
	for _, u := range *users {
		err := addV2rayVmessUser(*hsClient, &u)
		if err != nil {
			return err
		}
	}
	return
}
