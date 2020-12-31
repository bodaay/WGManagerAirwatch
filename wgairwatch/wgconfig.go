package wgairwatch

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"
)

const defaultAPIListenAdderss = "0.0.0.0"
const defaultAPIListenPort = 9696
const defaultAPIUseTLS = false
const defaultAPICertFile = "/etc/ssl/wgmanairwatch/wgmanairwatch.cert"
const defaultAPIKeyFile = "/etc/ssl/wgmanairwatch/wgmanairwatch.key"

var defaultAllowedIPsCIDR = []string{"0.0.0.0/32"}

const defaultWGMangerAddressIP = "127.0.0.1"
const defaultWGMangerPort = 6969

var defaultAllocateClientEventsIDS = []uint64{642}
var defaultRevokeClientEventsIDS = []uint64{642}
var defaultWGConfigMaps = []*WGConfigAirwatchInstanceMap{
	{
		OrganizationName:  "work1",
		MapWGInstanceName: "wg01",
	},
}

//WGConfig Global Configuration For WGManager
type WGConfigAirwatch struct {
	sync.Mutex
	APIListenAddress        string                         `json:"api_listen_address"`
	APIListenPort           uint16                         `json:"api_listen_port"`
	APIUseTLS               bool                           `json:"api_use_tls"`
	APITLSCert              string                         `json:"apitls_cert"`
	APITLSKey               string                         `json:"apitls_key"`
	APIAllowedIPS           []string                       `json:"api_allowed_ips"`
	WGManagerAddressIP      string                         `json:"wgmanager_address_ip"`
	WGManagerPort           uint16                         `json:"wgmanager_address_port"`
	AllocateClientEventsIDs []uint64                       `json:"allocate_client_event_ids"`
	RevokeClientEventsIDs   []uint64                       `json:"revoke_client_event_ids"`
	WGConfigAirwatchMaps    []*WGConfigAirwatchInstanceMap `json:""wgconfig_maps`
}
type WGConfigAirwatchInstanceMap struct {
	OrganizationName  string `json:"organization_name"`
	MapWGInstanceName string `json:"map_wg_instance_name"`
}

//WGInstanceConfig Per Instance Configuration

//CreateDefaultconfig Create Default Config file based on our constants
func (w *WGConfigAirwatch) CreateDefaultconfig(configpath string) (*WGConfigAirwatch, error) {
	var wgdefault WGConfigAirwatch
	wgdefault.APIListenAddress = defaultAPIListenAdderss
	wgdefault.APIListenPort = defaultAPIListenPort

	wgdefault.APIUseTLS = defaultAPIUseTLS
	wgdefault.APITLSCert = defaultAPICertFile
	wgdefault.APITLSKey = defaultAPIKeyFile
	wgdefault.APIAllowedIPS = defaultAllowedIPsCIDR
	wgdefault.WGManagerAddressIP = defaultWGMangerAddressIP
	wgdefault.WGManagerPort = defaultWGMangerPort
	wgdefault.AllocateClientEventsIDs = defaultAllocateClientEventsIDS
	wgdefault.RevokeClientEventsIDs = defaultRevokeClientEventsIDS
	wgdefault.WGConfigAirwatchMaps = defaultWGConfigMaps
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
