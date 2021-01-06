package airwatchapi

import (
	"WGManagerAirwatch/utils"
	"WGManagerAirwatch/wgairwatch"
	"fmt"
	"log"
	"strings"

	"github.com/go-resty/resty/v2"
)

/*
{
    "Instancename":"wg01",
    "Clientuuid":"client1"
}
*/

func AllocateClient(uniqueIdentifier string, mapName string, wgConfig *wgairwatch.WGConfigAirwatch) (vpnconfig string, err error) {
	client := resty.New()
	url := ""
	if wgConfig.WGManagerUseTLS {
		url = fmt.Sprintf("https://%s:%d/api/client", wgConfig.WGManagerAddressIP, wgConfig.WGManagerPort)
	} else {
		url = fmt.Sprintf("http://%s:%d/api/client", wgConfig.WGManagerAddressIP, wgConfig.WGManagerPort)
	}
	instanceName := strings.ToLower(wgConfig.WGConfigDefaultInstanceName) //by default will set it to default instance name
	for _, imrans := range wgConfig.WGConfigAirwatchMaps {
		if utils.TrimString(strings.ToLower(imrans.OrganizationName)) == utils.TrimString(strings.ToLower(mapName)) {
			instanceName = strings.ToLower(imrans.MapWGInstanceName)
			break
		}
	}
	log.Printf("Matched instance name for that organizaion name: %s\n", instanceName)
	postbody := fmt.Sprintf(`{"Instancename":"%s", "Clientuuid":"%s"}`, instanceName, uniqueIdentifier)
	res, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody([]byte(postbody)).
		Post(url)
	if err != nil {
		return "", err
	}
	// log.Println(res.String())
	return res.String(), nil
}

func RevokeClient(uniqueIdentifier string, mapName string, wgConfig *wgairwatch.WGConfigAirwatch) error {
	client := resty.New()
	url := ""
	if wgConfig.WGManagerUseTLS {
		url = fmt.Sprintf("https://%s:%d/api/client/all", wgConfig.WGManagerAddressIP, wgConfig.WGManagerPort)
	} else {
		url = fmt.Sprintf("http://%s:%d/api/client/all", wgConfig.WGManagerAddressIP, wgConfig.WGManagerPort)
	}
	instanceName := utils.TrimString(strings.ToLower(wgConfig.WGConfigDefaultInstanceName)) //by default will set it to default instance name
	for _, imrans := range wgConfig.WGConfigAirwatchMaps {
		if strings.ToLower(imrans.MapWGInstanceName) == strings.ToLower(mapName) {
			instanceName = strings.ToLower(imrans.MapWGInstanceName)
			break
		}
	}
	postbody := fmt.Sprintf(`{"Instancename":"%s", "Clientuuid":"%s"}`, instanceName, uniqueIdentifier)
	_, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody([]byte(postbody)).
		Delete(url)
	if err != nil {
		return err
	}
	// log.Println(res.String())
	return nil
}
