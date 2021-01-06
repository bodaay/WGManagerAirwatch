package main

import (
	"WGManagerAirwatch/airwatchevents"
	wgairwatch "WGManagerAirwatch/wgairwatch"
	"log"
	"os"
	"time"
)

func main() {
	defaultConfigFilePath := "wgmanairwatchconfig.json"
	if len(os.Args) > 1 {
		defaultConfigFilePath = os.Args[1]
	}
	// runningAsRoot, err := utils.CheckIfAdminOrRoot()
	// if err != nil {
	// 	panic(err)
	// }
	// if !runningAsRoot {
	// 	log.Fatalln("You must run this app as Admin or Root!")
	// }

	//Load the config file
	var wgc wgairwatch.WGConfigAirwatch
	err := wgc.ParseConfigFile(defaultConfigFilePath)
	if err != nil {
		newconfig, err := wgc.CreateDefaultconfig(defaultConfigFilePath)
		if err != nil {
			panic(err)
		}
		wgc = *newconfig
	}
	for {
		res, err := wgc.VerifyWGManager()
		if err != nil {
			log.Println(err)
			time.Sleep(3 * time.Second)
			continue
		}
		log.Println(res.String())
		break
	}

	airwatchevents.StartEventsClient(&wgc)
	// webapi.StartAdminClient(&wgc)
}
