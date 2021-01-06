package wgairwatch

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"
)

const defaultAirwatchEventsListenAdderss = "0.0.0.0"
const defaultAirwatchEventsListenPort = 9696

// const defaultAirwatchEventsUseTLS = false
const defaultAirwatchEventsCertFile = "/etc/ssl/wgmanairwatch/wgmanairwatch.cert"
const defaultAirwatchEventsKeyFile = "/etc/ssl/wgmanairwatch/wgmanairwatch.key"

var defaultAirwatchEventsIPsCIDR = []string{"0.0.0.0/0", "192.168.101.30/32", "192.168.101.35/32", "192.168.101.40/32"}

const defaultWGMangerAddressIP = "127.0.0.1"
const defaultWGMangerPort = 6969

var defaultAllocateClientEventsIDS = []uint64{148, 642}
var defaultRevokeClientEventsIDS = []uint64{39, 662}
var defaultWGConfigMaps = []*WGConfigAirwatchInstanceMap{
	{
		OrganizationName:  "work1",
		MapWGInstanceName: "wg01",
	},
}

const defaultInstanceName = "wg01"
const defaultAirwatchRestAPIAddress = "mymdm.airwatch.org"

//WGConfig Global Configuration For WGManager
type WGConfigAirwatch struct {
	sync.Mutex
	AirwatchEventsListenAddress string                         `json:"airwatch_events_listen_address"`
	AirwatchEventsListenPort    uint16                         `json:"airwatch_events_listen_port"`
	AirwatchEventsUseTLS        bool                           `json:"airwatch_events_use_tls"`
	AirwatchEventsTLSCert       string                         `json:"airwatch_events_tls_cert"`
	AirwatchEventsTLSKey        string                         `json:"airwatch_events_tls_key"`
	AirwatchEventsAllowedIPS    []string                       `json:"airwatch_events_allowed_ips"`
	WGManagerAddressIP          string                         `json:"wgmanager_address_ip"`
	WGManagerPort               uint16                         `json:"wgmanager_address_port"`
	WGManagerUseTLS             bool                           `json:"wgmanager_use_tls"`
	AllocateClientEventsIDs     []uint64                       `json:"allocate_client_event_ids"`
	RevokeClientEventsIDs       []uint64                       `json:"revoke_client_event_ids"`
	WGConfigAirwatchMaps        []*WGConfigAirwatchInstanceMap `json:"wgconfig_airwatch_maps"`
	WGConfigDefaultInstanceName string                         `json:"wgconfig_default_instance_name"`
	AirwatchRestAPIAddress      string                         `json:"airwatch_rest_api_address"`
	AirwatchRestAPIToken        string                         `json:"airwatch_rest_api_token"`
	AirwatchRestAPIUsername     string                         `json:"airwatch_rest_api_username"`
	AirwatchRestAPIPassword     string                         `json:"airwatch_rest_api_password"`
}
type WGConfigAirwatchInstanceMap struct {
	OrganizationName  string `json:"organization_name"`
	MapWGInstanceName string `json:"map_wg_instance_name"`
}

//WGInstanceConfig Per Instance Configuration

//CreateDefaultconfig Create Default Config file based on our constants
func (w *WGConfigAirwatch) CreateDefaultconfig(configpath string) (*WGConfigAirwatch, error) {
	var wgdefault WGConfigAirwatch
	wgdefault.AirwatchEventsListenAddress = defaultAirwatchEventsListenAdderss
	wgdefault.AirwatchEventsListenPort = defaultAirwatchEventsListenPort

	// wgdefault.AirwatchEventsUseTLS = defaultAirwatchEventsUseTLS
	wgdefault.AirwatchEventsTLSCert = defaultAirwatchEventsCertFile
	wgdefault.AirwatchEventsTLSKey = defaultAirwatchEventsKeyFile
	wgdefault.AirwatchEventsAllowedIPS = defaultAirwatchEventsIPsCIDR
	wgdefault.WGManagerAddressIP = defaultWGMangerAddressIP
	wgdefault.WGManagerPort = defaultWGMangerPort
	wgdefault.AllocateClientEventsIDs = defaultAllocateClientEventsIDS
	wgdefault.RevokeClientEventsIDs = defaultRevokeClientEventsIDS
	wgdefault.WGConfigAirwatchMaps = defaultWGConfigMaps
	wgdefault.WGConfigDefaultInstanceName = defaultInstanceName
	wgdefault.AirwatchRestAPIAddress = defaultAirwatchRestAPIAddress

	err := wgdefault.SaveConfigFile(configpath)
	if err != nil {
		return nil, err
	}
	return &wgdefault, nil
}

//ParseConfigFile Parse Config File by specified path
func (w *WGConfigAirwatch) ParseConfigFile(configpath string) error {
	data, err := ioutil.ReadFile(configpath)
	if err != nil {
		return err
	}
	err = w.ParseConfig(string(data))
	if err != nil {
		return err
	}
	return nil
}

//ParseConfig Parse Config string
func (w *WGConfigAirwatch) ParseConfig(configstring string) error {
	err := json.Unmarshal([]byte(configstring), w)
	if err != nil {
		return err
	}

	return nil
}

//SaveConfigFile Save the file into the specified path
func (w *WGConfigAirwatch) SaveConfigFile(configpath string) error {
	jsondata, err := json.MarshalIndent(w, "", "  ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(configpath, jsondata, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}
