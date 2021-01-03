package wgairwatch

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

// VerifyWGManager call WGManager and verify connectivity
func (w *WGConfigAirwatch) VerifyWGManager() (res *resty.Response, err error) {
	client := resty.New()

	if w.WGManagerUseTLS {
		res, err = client.R().Get(fmt.Sprintf("https://%s:%d/api/client", w.WGManagerAddressIP, w.WGManagerPort))
		if err != nil {
			return nil, err
		}
	} else {
		res, err = client.R().Get(fmt.Sprintf("http://%s:%d/api/client", w.WGManagerAddressIP, w.WGManagerPort))
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

// AllocateClient Allocate New Client
func (w *WGConfigAirwatch) AllocateClient(instanceName string) error {
	return nil
}

// RevokeClient Get All Instance in WGManager
func (w *WGConfigAirwatch) RevokeClient(instanceName string) error {
	return nil
}
