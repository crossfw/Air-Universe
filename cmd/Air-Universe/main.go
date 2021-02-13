package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/crossfw/Air-Universe/pkg/IPControl"
	sspApi "github.com/crossfw/Air-Universe/pkg/SSPanelAPI"
	"github.com/crossfw/Air-Universe/pkg/structures"
	log "github.com/sirupsen/logrus"
	"os"
	"sync"
	"time"
	_ "time"
)

const (
	VERSION = "0.3.0"
)

type WaitGroupWrapper struct {
	sync.WaitGroup
}

func init() {
	log.SetLevel(log.DebugLevel)

	var (
		printVersion bool
		configPath   string
	)
	//log.SetReportCaller(true)

	flag.BoolVar(&printVersion, "v", false, "print version")
	flag.StringVar(&configPath, "c", "locTest\\test.json", "configure file")
	flag.Parse()

	if printVersion {
		_, _ = fmt.Fprintf(os.Stdout, "Air-Universe %s\n", VERSION)
		os.Exit(0)
	}

	if configPath != "" {
		_, err := ParseBaseConfig(&configPath)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
		err = checkCfg()
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
		return
	}

	if flag.NFlag() == 0 {
		flag.PrintDefaults()
		os.Exit(0)
	}

}

func nodeSync(idIndex uint32, w *WaitGroupWrapper) (err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
			err = errors.New(fmt.Sprintf("%v (nodeId) main thread error - %s", baseCfg.Panel.NodeIDs[idIndex], r))
			log.Error(err)
			w.Done()
		}
	}()
	var (
		usersBefore, usersNow *[]structures.UserInfo
		usersTraffic          *[]structures.UserTraffic
		apiClient             structures.ProxyCommand
		nodeNow               *structures.NodeInfo
	)
	nodeNow = new(structures.NodeInfo)
	usersBefore = new([]structures.UserInfo)
	usersNow = new([]structures.UserInfo)
	usersTraffic = new([]structures.UserTraffic)

	// Get gRpc client and init v2ray api connection
	apiClient, err = initProxyCore()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	for {
		changed, err := getNodeInfo(nodeNow, idIndex)
		if err != nil {
			log.Error(err)
		}
		// Try add first, if no error cause, it's the first time to add, else remove then add until no error
		if changed == true && baseCfg.Proxy.AutoGenerate == true {
			err = apiClient.AddInbound(nodeNow)
			for err != nil {
				err = apiClient.RemoveInbound(nodeNow)
				err = apiClient.AddInbound(nodeNow)
				if err == nil {
					break
				}
				log.Warnf("Add inbound Failed", err)
				time.Sleep(time.Duration(baseCfg.Sync.FailDelay) * time.Second)
			}
			log.Printf("Added inbound %s", nodeNow.Tag)
		}

		usersNow, err = getUsers(nodeNow)
		if err != nil {
			log.Error(err)
		}
		useRemove, userAdd, err := structures.FindUserDiffer(usersBefore, usersNow)
		if err != nil {
			log.Error(err)
		}

		// Remove first, if user change uuid, remove old then add new.
		if useRemove != nil {
			log.Debugf(fmt.Sprint("Remove users ", *useRemove))
			err = apiClient.RemoveUsers(useRemove)
			if err != nil {
				log.Error(err)
			}
		}

		if userAdd != nil {
			log.Debugf(fmt.Sprint("Add users ", *userAdd))
			err = apiClient.AddUsers(userAdd)
			if err != nil {
				log.Error(err)
			}
		}

		// Sync_interval
		time.Sleep(time.Duration(baseCfg.Sync.Interval) * time.Second)

		usersTraffic, err = apiClient.QueryUsersTraffic(usersNow)
		if err != nil {
			log.Error(err)
		}
		log.Debugf(fmt.Sprint("Traffic data ", *usersTraffic))
		for err != nil {
			_, err = postUsersTraffic(nodeNow, usersTraffic)
			if err != nil {
				log.Error(err)
			}
		}

		usersBefore = usersNow
	}
}

func postUsersIP(w *WaitGroupWrapper) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("post users' IP thread error - %s", r))
			log.Error(err)
			w.Done()
		}
	}()
	for {
		time.Sleep(time.Duration(baseCfg.Sync.Interval) * time.Second)
		usersIp, err := IPControl.ReadLog(baseCfg)
		if err != nil {
			log.Error(err)
		}
		log.Debugf(fmt.Sprint("IP data ", *usersIp))
		ret, err := sspApi.PostUsersIP(baseCfg, usersIp)
		if ret != 1 || err != nil {
			log.Error(err)
		}
	}
}

func main() {
	var wg *WaitGroupWrapper
	wg = new(WaitGroupWrapper)

	for idIndex := 0; idIndex < len(baseCfg.Panel.NodeIDs); idIndex++ {
		wg.Add(1)
		go nodeSync(uint32(idIndex), wg)
	}
	wg.Add(1)
	go postUsersIP(wg)

	// wait
	wg.Wait()
}
