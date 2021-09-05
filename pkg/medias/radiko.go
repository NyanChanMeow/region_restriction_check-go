package medias

import (
	"bytes"

	"github.com/valyala/fasthttp"
)

func CheckRadiko(m *Media) (result *CheckResult) {
	m.URL = "https://radiko.jp/area?_=1625406539531"
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

	switch resp.StatusCode() {
	case fasthttp.StatusOK:
		if bytes.Contains(resp.Body(), []byte("JAPAN")) {
			result.Yes()
		} else if bytes.Contains(resp.Body(), []byte(`class="OUT"`)) {
			result.No()
		} else {
			result.Unexpected("body unsupported")
		}
	default:
		result.UnexpectedStatusCode(resp.StatusCode())
	}

	return result
}
