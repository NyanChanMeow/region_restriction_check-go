package medias

import (
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

func CheckPCRJP(m *Media) *CheckResult {
	m.Logger.Infoln("running")
	if m.URL == "" {
		m.URL = "https://api-priconne-redive.cygames.jp/"
	}
	if _, ok := m.Headers["User-Agent"]; !ok {
		m.Headers["User-Agent"] = UA_Dalvik
	}
	result := &CheckResult{Media: m.Name, Region: m.Region, Type: "Game"}

	resp, err := m.Do()
	if err != nil {
		m.Logger.Errorln(err)
		result.Failed(err)
		return result
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
	m.Logger.WithFields(log.Fields{
		"status_code": resp.StatusCode(),
		"result":      result.Result,
		"message":     result.Message,
	}).Infoln("done")
	return result
}

func CheckUMAJP(m *Media) *CheckResult {
	if m.URL == "" {
		m.URL = "https://api-umamusume.cygames.jp/"
	}
	return CheckPCRJP(m)
}
