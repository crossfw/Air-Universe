package main

import (
	"flag"
	"fmt"
	v2rayApi "github.com/crossfw/Air-Universe/pkg/V2rayApi"
	"github.com/crossfw/Air-Universe/pkg/structures"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
	_ "time"
)

const (
	VERSION = "0.0.5"
)

func init() {
	log.SetLevel(log.DebugLevel)

	var (
		printVersion bool
		configPath   string
	)

	flag.BoolVar(&printVersion, "V", false, "print version")
	flag.StringVar(&configPath, "C", "config/Air-Universe_json/test.json", "configure file")
	flag.Parse()

	if printVersion {
		_, _ = fmt.Fprintf(os.Stdout, "V2ray-ssp %s\n", VERSION)
		os.Exit(0)
	}

	if configPath != "" {
		_, err := parseBaseConfig(&configPath)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
		err = checkCfg()
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	}

	//if flag.NFlag() == 0 {
	//	flag.PrintDefaults()
	//	os.Exit(0)
	//}

}

func sync(idIndex uint32) error {
	var (
		err                   error
		v2Client              v2rayController
		usersBefore, usersNow *[]structures.UserInfo
		usersTraffic          *[]structures.UserTraffic
	)
	usersBefore = new([]structures.UserInfo)
	usersNow = new([]structures.UserInfo)
	usersTraffic = new([]structures.UserTraffic)

	for {
		v2Client.HsClient, v2Client.SsClient, err = v2rayApi.V2InitApi(baseCfg)
		if err != nil {
			log.Error(err)
		}
		usersNow, err = GetUserSelector(idIndex)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
		useRemove, userAdd, err := structures.FindUserDiffer(usersBefore, usersNow)
		if err != nil {
			log.Error(err)
		}

		if userAdd != nil {
			log.Debugf(fmt.Sprint("Add users ", *userAdd))
			err = v2Client.v2rayAddUsers(userAdd)
			if err != nil {
				log.Error(err)
			}
		}
		if useRemove != nil {
			log.Debugf(fmt.Sprint("Remove users ", *useRemove))
			err = v2Client.v2rayRemoveUsers(useRemove)
			if err != nil {
				log.Error(err)
			}
		}

		// Sync_interval
		time.Sleep(time.Duration(baseCfg.Sync.Interval) * time.Second)

		usersTraffic, err = v2Client.v2rayQueryTraffic(usersNow)
		if err != nil {
			log.Error(err)
		}
		_, err = PostUserSelector(idIndex, usersTraffic)
		if err != nil {
			log.Error(err)
		}
		usersBefore = usersNow
	}
}

func main() {
	log.Debug("start")

	for idIndex := 0; idIndex < len(baseCfg.Panel.NodeIDs); idIndex++ {
		go sync(uint32(idIndex))
	}

	time.Sleep(time.Duration(100000) * time.Second)
}

//	log.Println("Started.")
//log.Println("Getting users from", cfg.Url)
//
//// Get users until success
//now, err := v2rayssp.GetUser(cfg.Url, cfg.Key, cfg.NodeId)
//for err != nil {
//	log.Println("Failed to get users")
//	now, err = v2rayssp.GetUser(cfg.Url, cfg.Key, cfg.NodeId)
//	time.Sleep(time.Duration(cfg.FailDelay) * time.Second)
//}
//log.Println("Add users to V2ray-core.")
//for _, tag := range cfg.InTags {
//	_ = v2rayssp.AddUser(now, tag, cfg.APIAddress, cfg.APIPort, cfg.AlertId)
//}
//log.Println("Finish adding users.")
//
//// To find difference in 2 Users array
//before := now
//for {
//	log.Printf("waitting %v s \n", cfg.SyncInterval)
//	// Sleep and post traffic data and get new users
//	time.Sleep(time.Duration(cfg.SyncInterval) * time.Second)
//	nowTraffic, _ := v2rayssp.QueryTraffic(before, cfg.APIAddress, cfg.APIPort)
//	// POST until success, because we'll reset traffic data in each query.
//	log.Println("Post traffic to server.")
//	ret, err := v2rayssp.PostTraffic(cfg.Url, cfg.Key, cfg.NodeId, nowTraffic)
//	for err != nil && ret != 1 {
//		log.Println("Failed to post traffic")
//		ret, err = v2rayssp.PostTraffic(cfg.Url, cfg.Key, cfg.NodeId, nowTraffic)
//		time.Sleep(time.Duration(cfg.FailDelay) * time.Second)
//	}
//	log.Println("Finish posting traffic.")
//
//	// Get new users
//	log.Println("Getting users.")
//	now, err = v2rayssp.GetUser(cfg.Url, cfg.Key, cfg.NodeId)
//	for err != nil {
//		log.Println("Failed to get users")
//		now, err = v2rayssp.GetUser(cfg.Url, cfg.Key, cfg.NodeId)
//		time.Sleep(time.Duration(cfg.FailDelay) * time.Second)
//	}
//	log.Println("Finish getting users.")
//	log.Println("Add & remove users via users difference")
//	remove, add, _ := v2rayssp.FindUserDiffer(before, now)
//
//	for _, tag := range cfg.InTags {
//		v2rayssp.RemoveUser(remove, tag, cfg.APIAddress, cfg.APIPort, cfg.AlertId)
//		v2rayssp.AddUser(add, tag, cfg.APIAddress, cfg.APIPort, cfg.AlertId)
//	}
//	before = now
//	log.Println("Finish adding & removing")
//}
