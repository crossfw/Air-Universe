package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/crossfw/Air-Universe/pkg/IPControl"
	"github.com/crossfw/Air-Universe/pkg/SpeedLimitControl"
	"github.com/crossfw/Air-Universe/pkg/SysLoad"
	"github.com/crossfw/Air-Universe/pkg/structures"
	log "github.com/sirupsen/logrus"
	"os"
	"reflect"
	"sync"
	"time"
)

const (
	VERSION = "0.8.4"
)

type WaitGroupWrapper struct {
	sync.WaitGroup
}

func init() {
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
			log.Errorf("Failed to read config file - %s", err)
			os.Exit(1)
		}

		switch baseCfg.Log.LogLevel {
		case "debug":
			log.SetLevel(log.DebugLevel)
		case "info":
			log.SetLevel(log.InfoLevel)
		case "warning":
			log.SetLevel(log.WarnLevel)
		case "error":
			log.SetLevel(log.ErrorLevel)
		case "panic":
			log.SetLevel(log.PanicLevel)

		}

		if baseCfg.Log.Access != "" {
			file, err := os.OpenFile(baseCfg.Log.Access, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
			if err == nil {
				log.Infof("Log file will save at %s", baseCfg.Log.Access)
				log.SetOutput(file)
			} else {
				log.Warn("Failed to log to file, using default stderr")
			}
		}

		err = checkCfg()
		if err != nil {
			log.Errorf("Failed to check config file - %s", err)
			os.Exit(1)
		} else {
			log.Debugf("Successfully check config file")
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
			err = errors.New(fmt.Sprintf("%v (nodeId) main thread error - %s", baseCfg.Panel.NodeIDs[idIndex], r))
			log.Errorf("NodeID: %v IDIndex %v - Thread error - %s", baseCfg.Panel.NodeIDs[idIndex], idIndex, r)
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
	nodeID := baseCfg.Panel.NodeIDs[idIndex]
	nodeBefore = new(structures.NodeInfo)
	usersBefore = new([]structures.UserInfo)
	usersNow = new([]structures.UserInfo)
	usersTraffic = new([]structures.UserTraffic)

	// Get gRpc client and init v2ray api connection
	proxyClient, err = initProxyCore()
	if err != nil {
		log.Warnf("NodeID: %v IDIndex %v - Failed to init proxy-core client - %s", nodeID, idIndex, err)
		os.Exit(1)
	}
	panelClient, err = initPanel(idIndex)
	if err != nil {
		log.Warnf("NodeID: %v IDIndex %v - Failed to init panel client %s", nodeID, idIndex, err)
		os.Exit(1)
	} else {
		log.Debugf("NodeID: %v IDIndex %v - Successfully init", nodeID, idIndex)
	}

	log.Infof("NodeID: %v IDIndex %v - Successfully start sync thread", nodeID, idIndex)
	for {
		// Repeat until success
		for {
			err = panelClient.GetNodeInfo(baseCfg.Proxy.ForceCloseTLS)
			if err != nil {
				log.Warnf("NodeID: %v IDIndex %v - Failed to obtain node info - %s", nodeID, idIndex, err)
			} else {
				log.Debugf("NodeID: %v IDIndex %v - Successfully obtain node info", nodeID, idIndex)
				break
			}
			time.Sleep(time.Duration(baseCfg.Sync.FailDelay) * time.Second)
		}

		for {
			usersNow, err = panelClient.GetUser()
			if err != nil {
				log.Warnf("NodeID: %v IDIndex %v - Failed to obtain users info - %s", nodeID, idIndex, err)
			} else {
				log.Debugf("NodeID: %v IDIndex %v - Successfully obtain users info", nodeID, idIndex)
				break
			}
			time.Sleep(time.Duration(baseCfg.Sync.FailDelay) * time.Second)
		}

		if baseCfg.Proxy.SpeedLimitLevel != nil {
			err = SpeedLimitControl.AddLevel(usersNow, baseCfg.Proxy.SpeedLimitLevel)
			if err != nil {
				log.Warnf("NodeID: %v IDIndex %v - Failed to add level to users - %s", nodeID, idIndex, err)
			}
		}

		if reflect.DeepEqual(*panelClient.GetNowInfo(), *nodeBefore) == false && baseCfg.Proxy.AutoGenerate == true {
			for {
				// Is the first time to add inbound.
				if nodeBefore.Tag == "" {
					err = proxyClient.AddInbound(panelClient.GetNowInfo())
					if err != nil {
						log.Warnf("NodeID: %v IDIndex %v - Failed to add inbound - %s", nodeID, idIndex, err)
						// 第一次未成功先删除后添加
						nodeBefore.Tag = panelClient.GetNowInfo().Tag
					} else {
						break
					}
				} else {
					log.Infof("NodeID: %v IDIndex %v - Node info changed ", nodeID, idIndex)
					// 用户置0，在删除后添加
					usersBefore = new([]structures.UserInfo)
					// 对于已存在的节点，先remove再add
					err = proxyClient.RemoveInbound(nodeBefore)
					if err != nil {
						log.Warnf("NodeID: %v IDIndex %v - Failed to remove inbound - %s", nodeID, idIndex, err)
						//continue
					} else {
						log.Debugf("NodeID: %v IDIndex %v - Successfully remove inbound", nodeID, idIndex)
					}
					err = proxyClient.AddInbound(panelClient.GetNowInfo())
					if err != nil {
						log.Warnf("NodeID: %v IDIndex %v - Failed to add inbound - %s", nodeID, idIndex, err)
					} else {
						break
					}
				}
				time.Sleep(time.Duration(baseCfg.Sync.FailDelay) * time.Second)
			}

			*nodeBefore = *panelClient.GetNowInfo()
			log.Infof("NodeID: %v IDIndex %v - Successfully add inbound", nodeID, idIndex)
		}
		useRemove, userAdd, err := structures.FindUserDiffer(usersBefore, usersNow)
		if err != nil {
			log.Warnf("NodeID: %v IDIndex %v - Failed to process users info - %s", nodeID, idIndex, err)
		}

		// Remove first, if user changed uuid, remove old then add new.
		if len(*useRemove) > 0 {
			err = proxyClient.RemoveUsers(useRemove)
			if err != nil {
				log.Warnf("NodeID: %v IDIndex %v - Failed to remove users - %s", nodeID, idIndex, err)
			} else {
				log.Infof("NodeID: %v IDIndex %v - Remove users num: %v", nodeID, idIndex, len(*useRemove))
			}
		}

		if len(*userAdd) > 0 {
			err = proxyClient.AddUsers(userAdd)
			if err != nil {
				log.Warnf("NodeID: %v IDIndex %v - Failed to add users - %s", nodeID, idIndex, err)
			} else {
				log.Infof("NodeID: %v IDIndex %v - Add users num: %v", nodeID, idIndex, len(*userAdd))
			}
		}

		// Sync_interval
		time.Sleep(time.Duration(baseCfg.Sync.Interval) * time.Second)

		usersTraffic, err = proxyClient.QueryUsersTraffic(usersNow)
		if err != nil {
			log.Warnf("NodeID: %v IDIndex %v - Failed to query users traffic - %s", nodeID, idIndex, err)
		} else {
			log.Debugf("NodeID: %v IDIndex %v - Quary user traffic success - %+v", nodeID, idIndex, *usersTraffic)
		}

		// Every query will reset traffic statics, post traffic data will loop until success
		for {
			if len(*usersTraffic) > 0 {
				err = panelClient.PostTraffic(usersTraffic)
			}
			if err != nil {
				log.Warnf("NodeID: %v IDIndex %v - Failed to post users traffic - %s", nodeID, idIndex, err)
			} else {
				log.Debugf("NodeID: %v IDIndex %v - Post Traffic data success - %+v", nodeID, idIndex, *usersTraffic)
				break
			}
			time.Sleep(time.Duration(baseCfg.Sync.FailDelay) * time.Second)
		}
		if err != nil {
			log.Error(err)
		}

		loaData, err := SysLoad.GetSysLoad()
		if baseCfg.Panel.Type == "sspanel" {
			err = panelClient.PostSysLoad(loaData)
			if err != nil {
				log.Warnf("NodeID: %v IDIndex %v - Failed to post system load - %s", nodeID, idIndex, err)
			} else {
				log.Debugf("NodeID: %v IDIndex %v - Successfully post system load - %+v", nodeID, idIndex, *loaData)
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

	panelClient, err := initPanel(0)

	for {
		time.Sleep(time.Duration(baseCfg.Sync.PostIPInterval) * time.Second)
		usersIp, err := IPControl.ReadLog(baseCfg)
		if err != nil {
			log.Warnf("Failed to get alive IP - %s", err)
		} else {
			log.Debugf("Alive IP data %+v", *usersIp)
		}

		err = panelClient.PostAliveIP(baseCfg, usersIp)
		if err != nil {
			log.Warnf("Failed to post alive IP - %s", err)
		}

		err = IPControl.ClearLog(baseCfg)
		if err != nil {
			log.Warnf("Failed to clear proxy-core log - %s", err)
		}
	}
}

func main() {
	var wg *WaitGroupWrapper
	wg = new(WaitGroupWrapper)
	// delay 2 s to wait proxy-core start
	time.Sleep(time.Duration(2) * time.Second)
	for idIndex := 0; idIndex < len(baseCfg.Panel.NodeIDs); idIndex++ {
		wg.Add(1)
		go nodeSync(uint32(idIndex), wg)
		// 延迟执行，防止在多节点时面板和代理内核崩溃
		time.Sleep(time.Duration(5) * time.Second)
	}
	wg.Add(1)
	if baseCfg.Panel.Type == "sspanel" {
		go postUsersIP(wg)
	}

	// wait
	wg.Wait()
}
