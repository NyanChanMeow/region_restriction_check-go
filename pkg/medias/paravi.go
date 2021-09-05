package medias

import (
	"github.com/valyala/fasthttp"
)

func CheckParavi(m *Media) (result *CheckResult) {
	m.URL = "https://api.paravi.jp/api/v1/playback/auth"
	m.Logger.Infoln("running")

	if _, ok := m.Headers["User-Agent"]; !ok {
		m.Headers["User-Agent"] = UA_Browser
	}
	if _, ok := m.Headers[fasthttp.HeaderContentType]; !ok {
		m.Headers[fasthttp.HeaderContentType] = ContentTypeJSON
	}
	if m.Body == "" {
		m.Body = `{"meta_id":17414,"vuid":"3b64a775a4e38d90cc43ea4c7214702b","device_code":1,"app_id":1}`
	}
	result = &CheckResult{Media: m.Name, Region: m.Region}

	resp, err := m.Do()
	if err != nil {
		m.Logger.Errorln(err)
		result.Failed(err)
		return
	}
	defer fasthttp.ReleaseResponse(resp)

	switch resp.StatusCode() {
	case fasthttp.StatusUnauthorized:
		result.Yes()
	case fasthttp.StatusForbidden:
		result.No()
	default:
		result.UnexpectedStatusCode(resp.StatusCode())
	}

	return result
}
