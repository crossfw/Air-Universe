package main

import (
	"encoding/json"
	"fmt"
	"github.com/crossfw/Air-Universe/pkg/structures"
	log "github.com/sirupsen/logrus"
	"os"
)

// Default config
var (
	baseCfg = &structures.BaseConfig{
		Log: structures.Log{
			LogLevel: "info",
			Access:   "",
		},
		Panel: structures.Panel{
			Type: "sspanel",
		},
		Proxy: structures.Proxy{
			Type:           "xray",
			AlterID:        1,
			AutoGenerate:   true,
			InTags:         []string{},
			APIAddress:     "127.0.0.1",
			APIPort:        10085,
			LogPath:        "/var/log/au/xr.log",
			ForceCloseTLS:  false,
			EnableSniffing: true,
			Cert: structures.Cert{
				CertPath: "/usr/local/share/au/server.crt",
				KeyPath:  "/usr/local/share/au/server.key",
			},
		},
		Sync: structures.Sync{
			Interval:       60,
			FailDelay:      5,
			Timeout:        5,
			PostIPInterval: 90,
		},
	}
)

func ParseBaseConfig(configPath *string) (*structures.BaseConfig, error) {
	file, err := os.Open(*configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(baseCfg); err != nil {
		return nil, err
	}
	if baseCfg.Proxy.AutoGenerate == true {
		if len(baseCfg.Proxy.InTags) < len(baseCfg.Panel.NodeIDs) {
			log.Infof("InTags length isn't equal to nodeID length, adding inTags")
			for n := len(baseCfg.Proxy.InTags); n < len(baseCfg.Panel.NodeIDs); n++ {
				baseCfg.Proxy.InTags = append(baseCfg.Proxy.InTags, fmt.Sprintf("p%v", n))
			}
		}
	}
	return baseCfg, nil
}
