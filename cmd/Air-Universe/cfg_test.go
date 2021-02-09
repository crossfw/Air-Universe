package main

import (
	"fmt"
	"testing"
)

func TestConfigParse(t *testing.T) {
	var (
		cfgPath = "config\\v2rayssp_json\\v2rayssp_sample.json"
	)

	cfg, err := ParseBaseConfig(&cfgPath)
	if err != nil {
		t.Errorf("err\n")
	}
	fmt.Println(cfg)

}
