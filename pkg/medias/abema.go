package medias

import (
	"encoding/json"

	"github.com/valyala/fasthttp"
)

// {"isoCountryCode":"JP","timeZone":"Asia/Tokyo","utcOffset":"+09:00","cdnRegionUrl":"https://ds-linear-abematv.akamaized.net/region"}
func CheckAbemaTV(m *Media) (result *CheckResult) {
	m.URL = "https://api.abema.io/v1/ip/check?device=android"
	m.Logger.Infoln("running")

	if _, ok := m.Headers["User-Agent"]; !ok {
		m.Headers["User-Agent"] = UA_Dalvik
	}
	result = &CheckResult{Media: m.Name, Region: m.Region}

	resp, err := m.Do()
	if err != nil {
		m.Logger.Errorln(err)
		result.Failed(err)
		return
	}
	defer fasthttp.ReleaseResponse(resp)

	if resp.StatusCode() == fasthttp.StatusForbidden {
		result.No()
		return
	}

	r := make(map[string]string)
	err = json.Unmarshal(resp.Body(), &r)
	if err != nil {
		m.Logger.Errorln(err)
		result.Failed(err)
		return
	}

	if reg, ok := r["isoCountryCode"]; ok {
		if reg == "JP" {
			result.Yes()
		} else {
			result.Oversea()
		}
	} else {
		result.No()
	}
	return
}
