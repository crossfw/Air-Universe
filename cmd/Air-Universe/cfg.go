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
		Panel: structures.Panel{
			Type: "sspanel",
		},
		Proxy: structures.Proxy{
			Type:         "xray",
			AlertID:      1,
			AutoGenerate: true,
			InTags:       []string{"p0"},
			APIAddress:   "127.0.0.1",
			APIPort:      10085,
			LogPath:      "./v2.log",
		},
		Sync: structures.Sync{
			Interval:       60,
			FailDelay:      5,
			Timeout:        5,
			PostIPInterval: 300,
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
			log.Warnln("InTags length isn't equal to nodeID length, adding inTags")
			for n := len(baseCfg.Proxy.InTags); n < len(baseCfg.Panel.NodeIDs); n++ {
				baseCfg.Proxy.InTags = append(baseCfg.Proxy.InTags, fmt.Sprintf("p%v", n))
			}
		}
	}
	log.Println(*baseCfg)
	return baseCfg, nil
}
