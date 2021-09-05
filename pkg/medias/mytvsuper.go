package medias

import (
	"bytes"

	"github.com/valyala/fasthttp"
)

func CheckMyTVSuper(m *Media) (result *CheckResult) {
	m.URL = "https://www.mytvsuper.com/iptest.php"
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

	if bytes.Contains(resp.Body(), []byte("HK")) {
		result.Yes()
	} else {
		result.No()
	}

	return result
}
