package main

import (
	"fmt"
	"testing"
)

func TestConfigParse(t *testing.T) {
	var (
		cfgPath = "C:\\Users\\bendf\\Documents\\my-research\\air\\Air-Series\\Air-Universe\\config\\v2rayssp_json\\v2rayssp_sample.json"
	)

	cfg, err := parseBaseConfig(&cfgPath)
	if err != nil {
		t.Errorf("err\n")
	}
	fmt.Println(cfg)

}
