package medias

import (
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

func CheckKancolle(m *Media) *CheckResult {
	m.Logger.Infoln("running")
	if m.URL == "" {
		m.URL = "http://203.104.209.7/kcscontents/"
	}
	if _, ok := m.Headers["User-Agent"]; !ok {
		m.Headers["User-Agent"] = UA_Dalvik
	}
	result := &CheckResult{Media: m.Name, Region: m.Region}

	resp, err := m.Do()
	if err != nil {
		m.Logger.Errorln(err)
		result.Error = err
		result.Result = CheckResultFailed
		return result
	}
	defer fasthttp.ReleaseResponse(resp)

	switch resp.StatusCode() {
	case fasthttp.StatusOK:
		result.Result = CheckResultYes
	case fasthttp.StatusForbidden:
		result.Result = CheckResultNo
	default:
		result.Result = CheckResultUnexpected
	}

	m.Logger.WithFields(log.Fields{
		"status_code": resp.StatusCode(),
		"result":      result.Result,
		"error":       result.Error,
	}).Infoln("done")
	return result
}
