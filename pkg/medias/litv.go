package medias

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

func CheckLiTV(m *Media) *CheckResult {
	m.Method = "POST"
	m.Headers[fasthttp.HeaderContentType] = ContentTypeJSON
	m.Logger.Infoln("running")
	if m.URL == "" {
		m.URL = "https://www.litv.tv/vod/ajax/getUrl"
	}
	if _, ok := m.Headers["User-Agent"]; !ok {
		m.Headers["User-Agent"] = UA_Browser
	}
	if m.Body == "" {
		m.Body = `{"type":"noauth","assetId":"vod44868-010001M001_800K","puid":"6bc49a81-aad2-425c-8124-5b16e9e01337"}`
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
			m.Logger.Errorln(err)
			result.Failed(err)
			return result
		}

		if rr, ok := r["errorMessage"]; ok {
			if rr == nil {
				result.Yes()
			} else if rr.(string) == "vod.error.outsideregionerror" {
				result.No()
			} else {
				result.Unexpected(fmt.Sprintf("%+v", rr))
			}
		} else {
			result.Unexpected(`key "errorMessage" not found`)
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
