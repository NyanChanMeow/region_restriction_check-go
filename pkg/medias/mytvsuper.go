package medias

import (
	"bytes"

	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

func CheckMyTVSuper(m *Media) *CheckResult {
	m.Logger.Infoln("running")
	if m.URL == "" {
		m.URL = "https://www.mytvsuper.com/iptest.php"
	}
	if _, ok := m.Headers["User-Agent"]; !ok {
		m.Headers["User-Agent"] = UA_Browser
	}
	result := &CheckResult{Media: m.Name, Region: m.Region}

	resp, err := m.Do()
	if err != nil {
		m.Logger.Errorln(err)
		result.Failed(err)
		return result
	}
	defer fasthttp.ReleaseResponse(resp)

	if bytes.Contains([]byte("HK"), resp.Body()) {
		result.Yes()
	} else {
		result.No()
	}

	m.Logger.WithFields(log.Fields{
		"status_code": resp.StatusCode(),
		"result":      result.Result,
		"message":     result.Message,
	}).Infoln("done")
	return result
}
