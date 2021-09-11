package medias

import (
	"encoding/json"
	"fmt"

	"github.com/valyala/fasthttp"
)

func CheckTVBAnywhere(m *Media) (result *CheckResult) {
	m.URL = "https://uapisfm.tvbanywhere.com.sg/geoip/check/platform/android"
	m.Logger.Infoln("running")

	if _, ok := m.Headers["User-Agent"]; !ok {
		m.Headers["User-Agent"] = UA_Browser
	}
	result = &CheckResult{Media: m.Name, Region: m.Region}

	resp, err := m.Do()
	if err != nil {
		m.Logger.Errorln(err)
		result.Failed(err)
		return result
	}
	defer fasthttp.ReleaseResponse(resp)

	if resp.StatusCode() != fasthttp.StatusOK {
		result.UnexpectedStatusCode(resp.StatusCode())
	}
	r := make(map[string]interface{})
	err = json.Unmarshal(resp.Body(), &r)
	if err != nil {
		m.Logger.Errorln(err)
		result.Failed(err)
		return
	}

	if c, ok := r["allow_in_this_country"]; ok {
		if c.(bool) == true {
			result.Yes()
		} else if c.(bool) == false {
			result.No()
		} else {
			result.Unexpected(fmt.Sprintf("allow_in_this_country: %+v", c))
		}
	} else {
		result.Unexpected(`key "allow_in_this_country" not found`)
	}

	return
}
