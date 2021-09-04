package medias

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

func CheckCatchplay(m *Media) *CheckResult {
	m.Logger.Infoln("running")
	if m.URL == "" {
		m.URL = "https://sunapi.catchplay.com/geo"
	}
	if _, ok := m.Headers["User-Agent"]; !ok {
		m.Headers["User-Agent"] = UA_Browser
	}
	if _, ok := m.Headers["Authorization"]; !ok {
		m.Headers["Authorization"] = "Basic NTQ3MzM0NDgtYTU3Yi00MjU2LWE4MTEtMzdlYzNkNjJmM2E0Ok90QzR3elJRR2hLQ01sSDc2VEoy"
	}
	result := &CheckResult{Media: m.Name, Region: m.Region}

	resp, err := m.Do()
	if err != nil {
		m.Logger.Errorln(err)
		result.Message = err.Error()
		result.Result = CheckResultFailed
		return result
	}
	defer fasthttp.ReleaseResponse(resp)

	result.Result = CheckResultUnexpected
	if resp.StatusCode() == fasthttp.StatusOK {

		r := make(map[string]interface{})
		err = json.Unmarshal(resp.Body(), &r)
		if err != nil {
			m.Logger.Errorln(err)
			result.Message = err.Error()
			result.Result = CheckResultFailed
			return result
		}

		if rr, ok := r["code"]; ok {
			if rr.(string) == "0" {
				result.Result = CheckResultYes
			} else if rr.(string) == "100016" {
				result.Result = CheckResultNo
			} else {
				result.Message = fmt.Sprintf("code: %+v", rr)
			}
		} else {
			result.Message = fmt.Sprintf("code not found")
		}
	} else {
		result.Message = fmt.Sprintf("status code: %d", resp.StatusCode())
	}

	m.Logger.WithFields(log.Fields{
		"status_code": resp.StatusCode(),
		"result":      result.Result,
		"message":     result.Message,
	}).Infoln("done")
	return result
}
