package medias

import (
	"encoding/json"

	"github.com/valyala/fasthttp"
)

func CheckHBOGoAsia(m *Media) (result *CheckResult) {
	m.URL = "https://api2.hbogoasia.com/v1/geog?lang=undefined&version=0&bundleId=www.hbogoasia.com"
	m.Logger.Infoln("running")

	if _, ok := m.Headers["User-Agent"]; !ok {
		m.Headers["User-Agent"] = UA_Browser
	}
	result = &CheckResult{Media: m.Name, Region: "Asia"}

	resp, err := m.Do()
	if err != nil {
		m.Logger.Errorln(err)
		result.Failed(err)
		return
	}
	defer fasthttp.ReleaseResponse(resp)

	if resp.StatusCode() != fasthttp.StatusOK {
		result.UnexpectedStatusCode(resp.StatusCode())
		return
	}

	r := make(map[string]interface{})
	err = json.Unmarshal(resp.Body(), &r)
	if err != nil {
		m.Logger.Errorln(err)
		result.Failed(err)
		return
	}

	if c, ok := r["country"]; ok {
		result.Yes("Region:", c.(string))
	} else {
		result.No()
	}

	return
}
