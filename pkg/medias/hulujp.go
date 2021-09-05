package medias

import (
	"bytes"

	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

func CheckHuluJP(m *Media) *CheckResult {
	m.Logger.Infoln("running")
	if m.URL == "" {
		m.URL = "https://id.hulu.jp"
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

	switch resp.StatusCode() {
	case fasthttp.StatusFound:
		redirUrl := resp.Header.Peek("location")
		if bytes.Contains(redirUrl, []byte("login")) {
			result.Yes()
		} else if bytes.Contains(redirUrl, []byte("restrict")) {
			result.No()
		} else {
			result.Unexpected("location: " + string(redirUrl))
		}
	default:
		result.UnexpectedStatusCode(resp.StatusCode())
	}

	m.Logger.WithFields(log.Fields{
		"status_code": resp.StatusCode(),
		"result":      result.Result,
		"message":     result.Message,
	}).Infoln("done")
	return result
}
