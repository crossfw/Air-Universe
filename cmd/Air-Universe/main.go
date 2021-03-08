package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/crossfw/Air-Universe/pkg/IPControl"
	"github.com/crossfw/Air-Universe/pkg/SysLoad"
	"github.com/crossfw/Air-Universe/pkg/structures"
	log "github.com/sirupsen/logrus"
	"os"
	"reflect"
	"sync"
	"time"
)

const (
	VERSION = "0.3.4"
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
	flag.StringVar(&configPath, "c", "", "configure file")
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
		proxyClient           structures.ProxyCommand
		panelClient           structures.PanelCommand
		nodeBefore            *structures.NodeInfo
	)

	nodeBefore = new(structures.NodeInfo)
	usersBefore = new([]structures.UserInfo)
	usersNow = new([]structures.UserInfo)
	usersTraffic = new([]structures.UserTraffic)

	// Get gRpc client and init v2ray api connection
	proxyClient, err = initProxyCore()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	panelClient, err = initNode(idIndex)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	for {
		err = panelClient.GetNodeInfo()
		if err != nil {
			log.Error(err)
		}

		usersNow, err = panelClient.GetUser()
		if err != nil {
			log.Error(err)
		}
		// Try add first, if no error cause, it's the first time to add, else remove then add until no error
		if reflect.DeepEqual(*panelClient.GetNowInfo(), *nodeBefore) == false && baseCfg.Proxy.AutoGenerate == true {
			err = proxyClient.RemoveInbound(nodeBefore)
			err = proxyClient.AddInbound(panelClient.GetNowInfo())
			if err != nil {
				log.Warnf("Add inbound Error", err)
			}
			for err != nil {
				err = proxyClient.RemoveInbound(nodeBefore)
				if err != nil {
					log.Warnf("Remove inbound Error", err)
				}
				time.Sleep(time.Duration(30) * time.Second)
				err = proxyClient.AddInbound(panelClient.GetNowInfo())
				if err == nil {
					break
				}
				log.Warnf("Add inbound Failed", err)
				time.Sleep(time.Duration(baseCfg.Sync.FailDelay) * time.Second)
			}
			log.Printf("Added inbound %s", panelClient.GetNowInfo())
			usersBefore = new([]structures.UserInfo)
		}
		*nodeBefore = *panelClient.GetNowInfo()
		useRemove, userAdd, err := structures.FindUserDiffer(usersBefore, usersNow)
		if err != nil {
			log.Error(err)
		}

		// Remove first, if user change uuid, remove old then add new.
		if useRemove != nil {
			log.Debugf(fmt.Sprint("Remove users ", *useRemove))
			err = proxyClient.RemoveUsers(useRemove)
			if err != nil {
				log.Error(err)
			}
		}

		if userAdd != nil {
			log.Debugf(fmt.Sprint("Add users ", *userAdd))
			err = proxyClient.AddUsers(userAdd)
			if err != nil {
				log.Error(err)
			}
		}

		// Sync_interval
		time.Sleep(time.Duration(baseCfg.Sync.Interval) * time.Second)

		usersTraffic, err = proxyClient.QueryUsersTraffic(usersNow)
		if err != nil {
			log.Error(err)
		}
		log.Debugf(fmt.Sprint("Traffic data ", *usersTraffic))
		for {
			err = panelClient.PostTraffic(usersTraffic)
			if err != nil {
				log.Warnf("Failed to post traffic - %s", err)
				time.Sleep(time.Duration(baseCfg.Sync.FailDelay) * time.Second)
			} else {
				break
			}
		}
		if err != nil {
			log.Error(err)
		}

		loaData, err := SysLoad.GetSysLoad()
		err = panelClient.PostSysLoad(loaData)
		if err != nil {
			log.Warnf("Failed to post sys load - %s", err)
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

	panelClient, err := initNode(0)

	for {
		time.Sleep(time.Duration(baseCfg.Sync.PostIPInterval) * time.Second)
		usersIp, err := IPControl.ReadLog(baseCfg)
		if err != nil {
			log.Error(err)
		}
		log.Debugf(fmt.Sprint("IP data ", *usersIp))
		err = panelClient.PostAliveIP(baseCfg, usersIp)
		if err != nil {
			log.Error(err)
		}
		err = IPControl.ClearLog(baseCfg)
		if err != nil {
			log.Error("Clear log error - ", err)
		}
	}
}

func main() {
	var wg *WaitGroupWrapper
	wg = new(WaitGroupWrapper)

	for idIndex := 0; idIndex < len(baseCfg.Panel.NodeIDs); idIndex++ {
		wg.Add(1)
		go nodeSync(uint32(idIndex), wg)
		// 延迟执行，防止在多节点时面板和代理内核崩溃
		time.Sleep(time.Duration(1) * time.Second)
	}
	wg.Add(1)
	go postUsersIP(wg)

	// wait
	wg.Wait()
}
