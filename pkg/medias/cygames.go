package medias

import (
	"github.com/valyala/fasthttp"
)

func CheckPCRJP(m *Media) (result *CheckResult) {
	m.Logger.Infoln("running")
	if m.URL == "" {
		m.URL = "https://api-priconne-redive.cygames.jp/"
	}
	if _, ok := m.Headers["User-Agent"]; !ok {
		m.Headers["User-Agent"] = UA_Dalvik
	}
	result = &CheckResult{Media: m.Name, Region: m.Region, Type: "Game"}

	resp, err := m.Do()
	if err != nil {
		m.Logger.Errorln(err)
		result.Failed(err)
		return
	}
	defer fasthttp.ReleaseResponse(resp)

	switch resp.StatusCode() {
	case fasthttp.StatusNotFound:
		result.Yes()
	case fasthttp.StatusForbidden:
		result.No()
	default:
		result.UnexpectedStatusCode(resp.StatusCode())
	}

	return
}

func CheckUMAJP(m *Media) *CheckResult {
	m.URL = "https://api-umamusume.cygames.jp/"
	return CheckPCRJP(m)
}
