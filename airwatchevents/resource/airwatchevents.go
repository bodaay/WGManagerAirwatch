package resource

//AirWatchEventReceived is the JSON request recieved from airwatch
type AirWatchEventReceived struct {
	EventID                uint64 `json:"EventId"`
	EventType              string `json:"EventType"`
	DeviceID               uint64 `json:"DeviceId"`
	DeviceFriendlyName     string `json:"DeviceFriendlyName"`
	EnrollmentEmailAddress string `json:"EnrollmentEmailAddress"`
	EventTime              string `json:"EventTime"`
	EnrollmentStatus       string `json:"EnrollmentStatus"`
	CompromisedStatus      string `json:"CompromisedStatus"`
	CompromisedTimeStamp   string `json:"CompromisedTimeStamp"`
	ComplianceStatus       string `json:"ComplianceStatus"`
	Udid                   string `json:"Udid"`
	SerialNumber           string `json:"SerialNumber"`
	MACAddress             string `json:"MACAddress"`
	DeviceIMEI             string `json:"DeviceIMEI"`
	EnrollmentUserID       uint64 `json:"EnrollmentUserId"`
	AssetNumber            string `json:"AssetNumber"`
	Platform               string `json:"Platform"`
	OperatingSystem        string `json:"OperatingSystem"`
	Ownership              string `json:"Ownership"`
	SIMMCC                 string `json:"SIMMCC"`
	CurrentMCC             string `json:"CurrentMCC"`
	OrganizationGroupName  string `json:"OrganizationGroupName"`
}

// //AirWatchEventResponse response object for the airwatch event
// type AirWatchEventResponse struct {
// 	IsLocked bool                           `json:"islocked"`
// 	Profiles []AirWatchEventProfileResponse `json:"profiles"`
// }
// type AirWatchEventProfileResponse struct {
// 	Name             string                            `json:"name"`
// 	Interface        AirWatchEventInterfaceResponse    `json:"inferface"`
// 	Peers            []AirWatchEventPeerResponse       `json:"peers"`
// 	OnDemand         AirWatchEventOnDemandWifiResponse `json:"onDemandWifi"`
// 	OnDemandCellular bool                              `json:"onDemandCellular"`
// }
// type AirWatchEventInterfaceResponse struct {
// 	PrivateKey string   `json:"privateKey"`
// 	Addresses  string   `json:"addresses"`
// 	DNS        []string `json:"dns"`
// }

// type AirWatchEventPeerResponse struct {
// 	PublicKey           string   `json:"publicKey"`
// 	AllowedIPs          []string `json:"allowedIPs"`
// 	EndPoint            []string `json:"endPoint"`
// 	PersistentKeepAlive uint64   `json:"persistentKeepAlive"`
// }
// type AirWatchEventOnDemandWifiResponse struct {
// 	IsEnabled bool     `json:"isEnabled"`
// 	SSIDType  string   `json:"ssidType"`
// 	SSIDS     []string `json:"ssids"`
// }

type UpdateCustomAttributesRequest struct {
	CustomAttributes []CustomAttribute `json:"CustomAttributes"`
}
type CustomAttribute struct {
	Name             string `json:"Name"`
	Value            string `json:"Value"`
	ApplicationGroup string `json:"ApplicationGroup"`
}
