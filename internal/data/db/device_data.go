package db

import (
	"encoding/json"
	"net/http"

	"github.com/cifra-city/comtools/httpkit"
)

type DeviceData struct {
	UserAgent   string `json:"user_agent"`
	IP          string `json:"ip"`
	Fingerprint string `json:"fingerprint"`
}

func NewDeviceData(r *http.Request) ([]byte, error) {
	ip := httpkit.GetClientIP(r)
	userAgent := httpkit.GetUserAgent(r)
	fingerprint := httpkit.GenerateFingerprint(r)

	deviceData := DeviceData{
		UserAgent:   userAgent,
		IP:          ip,
		Fingerprint: fingerprint,
	}

	deviceDataJSON, err := json.Marshal(deviceData)
	if err != nil {
		return nil, err
	}

	return deviceDataJSON, nil
}
