package medias

import (
	"bytes"

	"github.com/valyala/fasthttp"
)

func CheckHuluJP(m *Media) (result *CheckResult) {
	m.URL = "https://id.hulu.jp"
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

	if resp.StatusCode() != fasthttp.StatusFound {
		result.UnexpectedStatusCode(resp.StatusCode())
		return
	}

	redirUrl := resp.Header.Peek("location")
	if bytes.Contains(redirUrl, []byte("login")) {
		result.Yes()
	} else if bytes.Contains(redirUrl, []byte("restrict")) {
		result.No()
	} else {
		result.Unexpected("location: " + string(redirUrl))
	}

	return result
}
