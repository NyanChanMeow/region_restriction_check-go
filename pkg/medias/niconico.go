package medias

import (
	"bytes"

	"github.com/valyala/fasthttp"
)

func CheckNiconico(m *Media) (result *CheckResult) {
	m.URL = "https://www.nicovideo.jp/watch/so23017073"
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

	if resp.StatusCode() != fasthttp.StatusOK {
		result.UnexpectedStatusCode(resp.StatusCode())
		return
	}

	if bytes.Contains(resp.Body(), []byte("同じ地域")) {
		result.No()
	} else {
		result.Yes()
	}

	return
}
