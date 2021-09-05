package medias

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

func CheckHamiVideo(m *Media) *CheckResult {
	m.Logger.Infoln("running")
	if m.URL == "" {
		m.URL = "https://hamivideo.hinet.net/api/play.do?id=OTT_VOD_0000249064&freeProduct=1"
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
		r := make(map[string]string)
		err = json.Unmarshal(resp.Body(), &r)
		if err != nil {
			result.Failed(err)
		} else {
			if c, ok := r["code"]; ok {
				if c == "06001-107" {
					result.Yes()
				} else if c == "06001-106" {
					result.No()
				} else {
					result.Unexpected(fmt.Sprintf("code: %s", c))
				}
			} else {
				result.Unexpected(`key "code" not found`)
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
