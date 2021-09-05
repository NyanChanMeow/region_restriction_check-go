package medias

import (
	"github.com/valyala/fasthttp"
)

func CheckKaraokeDAM(m *Media) (result *CheckResult) {
	m.URL = "http://cds1.clubdam.com/vhls-cds1/site/xbox/sample_1.mp4.m3u8"
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
		result.Yes()
	case fasthttp.StatusForbidden:
		result.No()
	default:
		result.UnexpectedStatusCode(resp.StatusCode())
	}

	return result
}
