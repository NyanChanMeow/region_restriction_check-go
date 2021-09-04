package medias

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

func CheckViuTV(m *Media) *CheckResult {
	m.Method = "POST"
	m.Headers[fasthttp.HeaderContentType] = ContentTypeJSON
	m.Logger.Infoln("running")
	if m.URL == "" {
		m.URL = "https://api.viu.now.com/p8/3/getLiveURL"
	}
	if _, ok := m.Headers["User-Agent"]; !ok {
		m.Headers["User-Agent"] = UA_Browser
	}
	if m.Body == "" {
		m.Body = `{"callerReferenceNo":"20210726112323","contentId":"099","contentType":"Channel","channelno":"099","mode":"prod","deviceId":"29b3cb117a635d5b56","deviceType":"ANDROID_WEB"}`
	}
	result := &CheckResult{Media: m.Name, Region: m.Region}

	resp, err := m.Do()
	if err != nil {
		m.Logger.Errorln(err)
		result.Message = err.Error()
		result.Result = CheckResultFailed
		return result
	}
	defer fasthttp.ReleaseResponse(resp)

	result.Result = CheckResultUnexpected
	if resp.StatusCode() == fasthttp.StatusOK {

		r := make(map[string]interface{})
		err = json.Unmarshal(resp.Body(), &r)
		if err != nil {
			m.Logger.Errorln(err)
			result.Message = err.Error()
			result.Result = CheckResultFailed
			return result
		}

		if rr, ok := r["responseCode"]; ok {
			switch rr {
			case "SUCCESS", "PRODUCT_INFORMATION_INCOMPLETE":
				result.Result = CheckResultYes
			case "GEO_CHECK_FAIL":
				result.Result = CheckResultNo
			default:
				result.Message = fmt.Sprintf("result: %s", rr)
			}
		} else {
			result.Message = fmt.Sprintf("responseCode not found")
		}
	} else {
		result.Message = fmt.Sprintf("status code: %d", resp.StatusCode())
	}

	m.Logger.WithFields(log.Fields{
		"status_code": resp.StatusCode(),
		"result":      result.Result,
		"message":     result.Message,
	}).Infoln("done")
	return result
}
