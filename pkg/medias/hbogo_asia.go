package medias

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

func CheckHBOGoAsia(m *Media) *CheckResult {
	m.Logger.Infoln("running")
	if m.URL == "" {
		m.URL = "https://api2.hbogoasia.com/v1/geog?lang=undefined&version=0&bundleId=www.hbogoasia.com"
	}
	if _, ok := m.Headers["User-Agent"]; !ok {
		m.Headers["User-Agent"] = UA_Browser
	}
	result := &CheckResult{Media: m.Name, Region: "Asia"}

	resp, err := m.Do()
	if err != nil {
		m.Logger.Errorln(err)
		result.Failed(err)
		return result
	}
	defer fasthttp.ReleaseResponse(resp)

	if resp.StatusCode() == fasthttp.StatusOK {
		r := make(map[string]interface{})
		err = json.Unmarshal(resp.Body(), &r)
		if err != nil {
			result.Failed(err)
		} else {
			if c, ok := r["country"]; ok {
				result.Yes()
				result.Message = "Region: " + c.(string)
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
