package medias

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

func CheckLineTV(m *Media) *CheckResult {
	m.Logger.Infoln("running")
	if m.URL == "" {
		m.URL = "https://www.linetv.tw/api/part/11829/eps/1/part?chocomemberId="
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
		r := make(map[string]interface{})
		err = json.Unmarshal(resp.Body(), &r)
		if err != nil {
			result.Failed(err)
		} else {
			if c, ok := r["countryCode"]; ok {
				cc := fmt.Sprintf("%v", c)
				if cc == "228" {
					result.Yes()
				} else if cc == "114" {
					result.No()
				} else {
					result.Unexpected(fmt.Sprintf("country code: %s", cc))
				}
			} else {
				result.Unexpected(`key "countryCode" not found`)
			}
		}
	} else if resp.StatusCode() == fasthttp.StatusForbidden {
		result.No()
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
