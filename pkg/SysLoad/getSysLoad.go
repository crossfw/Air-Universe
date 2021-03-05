package SysLoad

import (
	"github.com/crossfw/Air-Universe/pkg/structures"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
)

func GetSysLoad() (sysLoad *structures.SysLoad, err error) {
	sysLoad = new(structures.SysLoad)
	sLoad, err := load.Avg()
	if err != nil {
		return nil, err
	}
	sysLoad.Load1 = sLoad.Load1
	sysLoad.Load5 = sLoad.Load5
	sysLoad.Load15 = sLoad.Load15

	sysLoad.Uptime, err = host.Uptime()
	if err != nil {
		return nil, err
	}

	return
}
