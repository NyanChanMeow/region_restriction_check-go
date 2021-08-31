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
		m.Headers["User-Agent"] = "Dalvik/2.1.0 (Linux; U; Android 9; ALP-AL00 Build/HUAWEIALP-AL00)"
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
	case fasthttp.StatusNotFound:
		result.Result = CheckResultYes
	case fasthttp.StatusForbidden:
		result.Result = CheckResultNo
	default:
		result.Result = CheckResultUnexpected
	}
	m.Logger.WithFields(log.Fields{
		"status_code": resp.StatusCode(),
		"result":      result.Result,
	}).Infoln("done")
	return result
}

func CheckUMAJP(m *Media) *CheckResult {
	if m.URL == "" {
		m.URL = "https://api-umamusume.cygames.jp/"
	}
	return CheckPCRJP(m)
}
