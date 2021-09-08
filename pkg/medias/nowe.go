package medias

import (
	"encoding/json"
	"fmt"

	"github.com/valyala/fasthttp"
)

func CheckNowE(m *Media) (result *CheckResult) {
	m.URL = "https://webtvapi.nowe.com/16/1/getVodURL"
	m.Method = "POST"
	m.Headers[fasthttp.HeaderContentType] = ContentTypeJSON
	m.Logger.Infoln("running")

	if _, ok := m.Headers["User-Agent"]; !ok {
		m.Headers["User-Agent"] = UA_Browser
	}
	if m.Body == "" {
		m.Body = `{"contentId":"202105121370235","contentType":"Vod","pin":"","deviceId":"W-60b8d30a-9294-d251-617b-c12f9d0c","deviceType":"WEB"}`
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

	r := make(map[string]interface{})
	err = json.Unmarshal(resp.Body(), &r)
	if err != nil {
		m.Logger.Errorln(err)
		result.Failed(err)
		return
	}

	if rr, ok := r["responseCode"]; ok {
		switch rr {
		case "SUCCESS", "PRODUCT_INFORMATION_INCOMPLETE":
			result.Yes()
		case "GEO_CHECK_FAIL":
			result.No()
		default:
			result.Unexpected(fmt.Sprintf("result: %s", rr))
		}
	} else {
		result.Unexpected(`key "responseCode" not found`)
	}

	return
}
