package medias

import (
	"encoding/json"

	"github.com/valyala/fasthttp"
)

func CheckKKTV(m *Media) (result *CheckResult) {
	m.URL = "https://api.kktv.me/v3/ipcheck"
	m.Logger.Infoln("running")

	if _, ok := m.Headers["User-Agent"]; !ok {
		m.Headers["User-Agent"] = UA_Browser
	}
	result = &CheckResult{Media: m.Name, Region: m.Region}

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

	var r struct {
		Data struct {
			Country string `json:"country"`
		} `json:"data"`
	}
	err = json.Unmarshal(resp.Body(), &r)
	if err != nil {
		m.Logger.Errorln(err)
		result.Failed(err)
		return
	}

	if r.Data.Country == "TW" {
		result.Yes()
	} else {
		result.No()
	}

	return result
}
