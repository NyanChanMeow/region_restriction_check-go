package medias

import (
	"github.com/valyala/fasthttp"
)

func CheckKonosubaFD(m *Media) (result *CheckResult) {
	m.URL = "https://api.konosubafd.jp/api/masterlist"
	m.Logger.Infoln("running")

	if _, ok := m.Headers["User-Agent"]; !ok {
		m.Headers["User-Agent"] = "pj0007/212 CFNetwork/1240.0.4 Darwin/20.6.0"
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
	case fasthttp.StatusInternalServerError:
		result.Yes()
	case fasthttp.StatusForbidden:
		result.No()
	default:
		result.UnexpectedStatusCode(resp.StatusCode())
	}

	return result
}
