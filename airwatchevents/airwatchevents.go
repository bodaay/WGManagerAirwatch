package airwatchevents

import (
	"WGManagerAirwatch/airwatchapi"
	"WGManagerAirwatch/airwatchevents/resource"
	"WGManagerAirwatch/utils"
	"WGManagerAirwatch/wgairwatch"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

//{"EventId":662,"EventType":"Delete Device Requested","DeviceId":81,"DeviceFriendlyName":"DELETE IN PROGRESS...","EnrollmentEmailAddress":"iphonexs2@ssb.local","EnrollmentUserName":"iphonexs2","EventTime":"2021-01-05T09:55:10.7702835Z",
//"EnrollmentStatus":"Unknown","CompromisedStatus":null,"CompromisedTimeStamp":"0001-01-01T00:00:00","ComplianceStatus":null,"PhoneNumber":null,"Udid":"00008020-001A5D1E21E2002E","SerialNumber":"F2LXKH8DKPH2","MACAddress":null,
//"DeviceIMEI":null,"EnrollmentUserId":0,"AssetNumber":"00008020-001A5D1E21E2002E","Platform":null,"OperatingSystem":null,"Ownership":null,"SIMMCC":null,"CurrentMCC":null,"OrganizationGroupName":null}

func checkIPAccess(clientip string, allowedIPScidr []string) bool {
	return true
	ip := net.ParseIP(clientip)
	for _, aips := range allowedIPScidr {
		_, ipnet, err := net.ParseCIDR(aips)
		if err != nil {
			log.Println(err)
			return false
		}
		if ipnet.Contains(ip) {
			return true
		}
	}

	return false
}

//StartAdminClient start the REST API Echo Server for inserting watermark
func StartEventsClient(wgConfig *wgairwatch.WGConfigAirwatch) error {
	e := echo.New()
	// const subserviceIdentifier = "StartWebClient"
	configureClientWebServer(e)
	configureAllRoutesClient(e, wgConfig)
	address := (wgConfig.AirwatchEventsListenAddress + ":" + strconv.Itoa(int(wgConfig.AirwatchEventsListenPort)))
	//err := e.StartTLS(address, (config.RootCertFile), (config.RootCertKey))
	if wgConfig.AirwatchEventsUseTLS {
		e.StartTLS(address, (wgConfig.AirwatchEventsTLSCert), (wgConfig.AirwatchEventsTLSKey))

	} else {
		e.Start(address)

	}
	return nil
}
func configureClientWebServer(e *echo.Echo) {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.BodyLimit("100M"))
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))
	mime.AddExtensionType(".js", "application/javascript") //This will solve some windows shit issue, when it will serve javascript file as text/plain, read more about it at:https://github.com/labstack/echo/issues/1038

}

func configureAllRoutesClient(e *echo.Echo, wgConfig *wgairwatch.WGConfigAirwatch) {
	postEventsReceived(e, wgConfig)
	// deleteEventsReceived(e, wgConfig)
	// putEventsReceived(e, wgConfig)
	getEventsReceived(e, wgConfig)

}

func postEventsReceived(e *echo.Echo, wgConfig *wgairwatch.WGConfigAirwatch) {
	e.POST("/events", func(c echo.Context) error {
		IsAllowed := checkIPAccess(c.RealIP(), wgConfig.AirwatchEventsAllowedIPS)
		if !IsAllowed {
			return c.String(http.StatusUnauthorized, fmt.Sprintf("You are not allowed to access, ip: %s", c.RealIP()))
		}
		u := new(resource.AirWatchEventReceived)

		if err := c.Bind(u); err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		url := fmt.Sprintf("%s/API/mdm/devices/%d/customattributes", wgConfig.AirwatchRestAPIAddress, u.DeviceID)
		defaultApplicationGroup := "com.addrek.wireguard"
		log.Println(url)
		//check if its allocation event trigger
		for _, eids := range wgConfig.AllocateClientEventsIDs {
			if eids != u.EventID {
				continue
			}
			log.Printf("Request to add ClientID: %s   For Organization Name: %s\n", u.AssetNumber, u.OrganizationGroupName)
			vpnconfig, err := airwatchapi.AllocateClient(u.AssetNumber, u.OrganizationGroupName, wgConfig)
			if err != nil {
				return c.String(http.StatusBadRequest, err.Error())
			}

			// log.Println(vpnconfig)

			entries := strings.Split(vpnconfig, "\n")

			m := make(map[string]string)
			for _, e := range entries {
				parts := strings.Split(e, "=")
				if len(parts) > 1 {
					m[strings.Trim(parts[0], " \n\r")] = strings.Trim(parts[1], " \n\r")
				}
			}

			cAttributes := resource.UpdateCustomAttributesRequest{}

			// wireguard.isLocked
			attributeIsLocked := resource.CustomAttribute{
				Name:             "wireguard.isLocked",
				Value:            "true",
				ApplicationGroup: defaultApplicationGroup,
			}
			cAttributes.CustomAttributes = append(cAttributes.CustomAttributes, attributeIsLocked)

			// wireguard.profile
			attributeProfile := resource.CustomAttribute{
				Name:             "wireguard.profile",
				Value:            fmt.Sprintf("%s,%s,%s", utils.TrimString(m["Address"]), utils.TrimString(m["PublicKey"]), utils.TrimString(m["PrivateKey"])),
				ApplicationGroup: defaultApplicationGroup,
			}
			cAttributes.CustomAttributes = append(cAttributes.CustomAttributes, attributeProfile)

			// wireguard.allowedIPS
			tempstring := ""
			for _, k := range strings.Split(m["AllowedIPs"], ",") {
				tempstring += k + ","
			}
			tempstring = tempstring[:len(tempstring)-1]
			attributeAllowedIPs := resource.CustomAttribute{
				Name:             "wireguard.allowedIPS",
				Value:            tempstring,
				ApplicationGroup: defaultApplicationGroup,
			}
			cAttributes.CustomAttributes = append(cAttributes.CustomAttributes, attributeAllowedIPs)

			// wireguard.dnsServers
			tempstring = ""
			for _, k := range strings.Split(m["DNS"], ",") {
				tempstring += utils.TrimString(k) + ","
			}
			tempstring = tempstring[:len(tempstring)-1]
			attributeDNSServers := resource.CustomAttribute{
				Name:             "wireguard.dnsServers",
				Value:            tempstring,
				ApplicationGroup: defaultApplicationGroup,
			}
			cAttributes.CustomAttributes = append(cAttributes.CustomAttributes, attributeDNSServers)

			// wireguard.endpoint
			attributeEndpoint := resource.CustomAttribute{
				Name:             "wireguard.endpoint",
				Value:            utils.TrimString(m["Endpoint"]),
				ApplicationGroup: defaultApplicationGroup,
			}
			cAttributes.CustomAttributes = append(cAttributes.CustomAttributes, attributeEndpoint)

			// wireguard.persistentKeepAlive
			attributePersistentKeepAlive := resource.CustomAttribute{
				Name:             "wireguard.persistentKeepAlive",
				Value:            utils.TrimString(m["PersistentKeepalive"]),
				ApplicationGroup: defaultApplicationGroup,
			}
			cAttributes.CustomAttributes = append(cAttributes.CustomAttributes, attributePersistentKeepAlive)

			// wireguard.onDemandWifiEnabled
			attributeOnDemandWifiEnabled := resource.CustomAttribute{
				Name:             "wireguard.onDemandWifiEnabled",
				Value:            "true",
				ApplicationGroup: defaultApplicationGroup,
			}
			cAttributes.CustomAttributes = append(cAttributes.CustomAttributes, attributeOnDemandWifiEnabled)

			// wireguard.onDemandWifiSsids
			// tempstring = ""
			// for _, k := range strings.Split(m["DNS"], ",") {
			// 	tempstring += utils.TrimString(k) + ","
			// }
			// tempstring = tempstring[:len(tempstring)-1]
			attributeOnDemandWifiSsids := resource.CustomAttribute{
				Name:             "wireguard.onDemandWifiSsids",
				Value:            "WifiSSid1,WifiSSid2",
				ApplicationGroup: defaultApplicationGroup,
			}
			cAttributes.CustomAttributes = append(cAttributes.CustomAttributes, attributeOnDemandWifiSsids)

			// wireguard.onDemandWifiSsidType
			attributeOnDemandWifiSsidType := resource.CustomAttribute{
				Name:             "wireguard.onDemandWifiSsidType",
				Value:            "any",
				ApplicationGroup: defaultApplicationGroup,
			}
			cAttributes.CustomAttributes = append(cAttributes.CustomAttributes, attributeOnDemandWifiSsidType)

			// wireguard.onDemandCelluar
			attributeOnDemandCelluar := resource.CustomAttribute{
				Name:             "wireguard.onDemandCelluar",
				Value:            "true",
				ApplicationGroup: defaultApplicationGroup,
			}
			cAttributes.CustomAttributes = append(cAttributes.CustomAttributes, attributeOnDemandCelluar)

			jsonData, err := json.MarshalIndent(cAttributes, "", "  ")
			if err != nil {
				return c.String(http.StatusBadRequest, err.Error())
			}

			log.Println(string(jsonData))
			client := resty.New()
			res, err := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}).R().
				SetHeader("Content-Type", "application/json").
				SetHeader("aw-tenant-code", wgConfig.AirwatchRestAPIToken).
				SetBasicAuth(wgConfig.AirwatchRestAPIUsername, wgConfig.AirwatchRestAPIPassword).
				SetBody(jsonData).
				Put(url)

			if err != nil {
				log.Println(err)
				return c.String(http.StatusBadRequest, err.Error())
			}
			log.Println(res.Status())
			return c.String(http.StatusOK, res.Status()) //no need to break out, since we are already returning from this shit
		}
		for _, eids := range wgConfig.RevokeClientEventsIDs {
			if eids != u.EventID {
				continue
			}
			err := airwatchapi.RevokeClient(u.AssetNumber, u.OrganizationGroupName, wgConfig)
			if err != nil {
				return c.String(http.StatusBadRequest, err.Error())
			}
			cAttributes := resource.UpdateCustomAttributesRequest{}
			// wireguard.profile
			attributeProfile := resource.CustomAttribute{
				Name:             "wireguard.profile",
				Value:            "none",
				ApplicationGroup: defaultApplicationGroup,
			}
			cAttributes.CustomAttributes = append(cAttributes.CustomAttributes, attributeProfile)
			jsonData, err := json.MarshalIndent(cAttributes, "", "  ")
			if err != nil {
				return c.String(http.StatusBadRequest, err.Error())
			}

			log.Println(string(jsonData))
			// client := resty.New()
			// res, err := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}).R().
			// 	SetHeader("Content-Type", "application/json").
			// 	SetHeader("aw-tenant-code", wgConfig.AirwatchRestAPIToken).
			// 	SetBasicAuth(wgConfig.AirwatchRestAPIUsername, wgConfig.AirwatchRestAPIPassword).
			// 	SetBody(jsonData).
			// 	Put(url)

			return c.String(http.StatusOK, "Client Revoked") //no need to break out, since we are already returning from this shit
		}
		// fmt.Printf("-----------------\n")
		// fmt.Printf("%+v", u)
		// fmt.Printf("-----------------\n")
		return c.JSONPretty(http.StatusOK, "no matching event id", "  ")
	})
}

// func deleteEventsReceived(e *echo.Echo, wgConfig *wgairwatch.WGConfigAirwatch) {
// 	e.DELETE("/events", func(c echo.Context) error {
// 		IsAllowed := checkIPAccess(c.RealIP(), wgConfig.AirwatchEventsAllowedIPS)
// 		if !IsAllowed {
// 			return c.String(http.StatusUnauthorized, fmt.Sprintf("You are not allowed to access, ip: %s", c.RealIP()))
// 		}
// 		body, err := ioutil.ReadAll(c.Request().Body)
// 		if err != nil {
// 			return err
// 		}
// 		fmt.Printf("-----------------\n")
// 		fmt.Println(string(body))
// 		fmt.Printf("-----------------\n")
// 		return c.JSONPretty(http.StatusOK, "ok", "  ")
// 	})
// }

// func putEventsReceived(e *echo.Echo, wgConfig *wgairwatch.WGConfigAirwatch) {
// 	e.PUT("/events", func(c echo.Context) error {
// 		IsAllowed := checkIPAccess(c.RealIP(), wgConfig.AirwatchEventsAllowedIPS)
// 		if !IsAllowed {
// 			return c.String(http.StatusUnauthorized, fmt.Sprintf("You are not allowed to access, ip: %s", c.RealIP()))
// 		}

// 		body, err := ioutil.ReadAll(c.Request().Body)
// 		if err != nil {
// 			return err
// 		}
// 		fmt.Printf("-----------------\n")
// 		fmt.Println(string(body))
// 		fmt.Printf("-----------------\n")
// 		return c.JSONPretty(http.StatusOK, "ok", "  ")
// 	})
// }

func getEventsReceived(e *echo.Echo, wgConfig *wgairwatch.WGConfigAirwatch) {
	e.GET("/events", func(c echo.Context) error {
		IsAllowed := checkIPAccess(c.RealIP(), wgConfig.AirwatchEventsAllowedIPS)
		if !IsAllowed {
			return c.String(http.StatusUnauthorized, fmt.Sprintf("You are not allowed to access, ip: %s", c.RealIP()))
		}

		body, err := ioutil.ReadAll(c.Request().Body)
		if err != nil {
			return err
		}
		fmt.Printf("-----------------\n")
		fmt.Println(string(body))
		fmt.Printf("-----------------\n")
		return c.JSONPretty(http.StatusOK, "ok", "  ")
	})
}
