package medias

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

func CheckKKTV(m *Media) *CheckResult {
	m.Logger.Infoln("running")
	if m.URL == "" {
		m.URL = "https://api.kktv.me/v3/ipcheck"
	}
	if _, ok := m.Headers["User-Agent"]; !ok {
		m.Headers["User-Agent"] = UA_Browser
	}
	result := &CheckResult{Media: m.Name, Region: m.Region}

	resp, err := m.Do()
	if err != nil {
		m.Logger.Errorln(err)
		result.Failed(err)
		return result
	}
	defer fasthttp.ReleaseResponse(resp)

	if resp.StatusCode() == fasthttp.StatusOK {
		var r struct {
			Data struct {
				Country string `json:"country"`
			} `json:"data"`
		}
		err = json.Unmarshal(resp.Body(), &r)
		if err != nil {
			result.Failed(err)
		} else {
			if r.Data.Country == "TW" {
				result.Yes()
			} else {
				result.No()
			}
		}
	} else {
		result.UnexpectedStatusCode(resp.StatusCode())
	}

	m.Logger.WithFields(log.Fields{
		"status_code": resp.StatusCode(),
		"result":      result.Result,
		"message":     result.Message,
	}).Infoln("done")
	return result
}
