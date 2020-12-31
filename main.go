package main

import (
	"WGManagerAirwatch/webapi"
	wgairwatch "WGManagerAirwatch/wgairwatch"
	"os"
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
	//Search the path for instances configuration files

	webapi.StartClient(&wgc)
	// webapi.StartAdminClient(&wgc)
}
