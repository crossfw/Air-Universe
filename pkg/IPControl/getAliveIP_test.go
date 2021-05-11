package IPControl

import (
	"fmt"
	"github.com/crossfw/Air-Universe/pkg/structures"
	"testing"
)

//# v2ray.com/core/common/buf
//vet.exe: ..\..\vendor\v2ray.com\core\common\buf\readv_constraint_windows.go:10:6: reason declared but not used
var (
	baseCfg = &structures.BaseConfig{
		Panel: structures.Panel{
			Type: "sspanel",
		},
		Proxy: structures.Proxy{
			Type:       "v2ray",
			AlterID:    1,
			InTags:     []string{"p0"},
			APIAddress: "127.0.0.1",
			APIPort:    10085,
			LogPath:    "\\locTest\\v2.log",
		},
		Sync: structures.Sync{
			Interval:  60,
			FailDelay: 5,
			Timeout:   5,
		},
	}
)

func TestGetIP(t *testing.T) {
	a, err := ReadLog(baseCfg)
	if err != nil {
		t.Errorf("Failed %s", err)
	}
	fmt.Println(*a)
}

func TestCaptureDetail(t *testing.T) {
	line := "2021/03/15 19:16:04 127.0.0.1:46948 accepted tcp:accepted-sngapm.qcloud.com:443 email: 3668-node2"
	ret, err := captureDetail(line)
	if err != nil {
		t.Errorf("Failed %s", err)
	}
	fmt.Println(ret)
}
