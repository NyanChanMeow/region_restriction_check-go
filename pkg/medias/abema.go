package medias

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

// {"isoCountryCode":"JP","timeZone":"Asia/Tokyo","utcOffset":"+09:00","cdnRegionUrl":"https://ds-linear-abematv.akamaized.net/region"}
func CheckAbemaTV(m *Media) *CheckResult {
	m.Logger.Infoln("running")
	if m.URL == "" {
		m.URL = "https://api.abema.io/v1/ip/check?device=android"
	}
	if _, ok := m.Headers["User-Agent"]; !ok {
		m.Headers["User-Agent"] = UA_Dalvik
	}
	result := &CheckResult{Media: m.Name, Region: m.Region}

	resp, err := m.Do()
	if err != nil {
		m.Logger.Errorln(err)
		result.Failed(err)
		return result
	}
	defer fasthttp.ReleaseResponse(resp)

	result.No()
	if resp.StatusCode() != fasthttp.StatusForbidden {
		r := make(map[string]string)
		err = json.Unmarshal(resp.Body(), &r)
		if err != nil {
			result.Unexpected(err)
		} else {
			if reg, ok := r["isoCountryCode"]; ok {
				if reg == "JP" {
					result.Yes()
				} else {
					result.No()
				}
			}
		}
	}

	m.Logger.WithFields(log.Fields{
		"status_code": resp.StatusCode(),
		"result":      result.Result,
		"message":     result.Message,
	}).Infoln("done")
	return result
}
