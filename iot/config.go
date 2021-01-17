package iot

const (
	// WsnToken is the token of wsn iot cloud
	WsnToken = "your_wsn_token"
	// WsnNumericalAPI is the api of wsn iot cloud for pushing numerical datapoints
	WsnNumericalAPI = "http://www.wsncloud.com/api/data/v1/numerical/insert"
	// WsnGenericAPI is the api of wsn iot cloud for pushing generic datapoints
	WsnGenericAPI = "http://www.wsncloud.com/api/data/v1/generic/insert"
)

const (
	// OneNetToken is the token of OneNet iot cloud
	OneNetToken = "your_onenet_token"
	// OneNetAPI is the api of OneNet iot cloud for pushing datapoints
	OneNetAPI = "http://api.heclouds.com/devices/540381180/datapoints"
)

// WsnConfig ...
type WsnConfig struct {
	Token string `json:"token"`
	API   string `json:"api"`
}

// OneNetConfig ...
type OneNetConfig struct {
	Token string `json:"token"`
	API   string `json:"api"`
}
