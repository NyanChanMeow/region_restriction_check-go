package medias

import (
	"encoding/json"
	"fmt"

	"github.com/valyala/fasthttp"
)

func CheckCatchplay(m *Media) (result *CheckResult) {
	m.URL = "https://sunapi.catchplay.com/geo"
	m.Logger.Infoln("running")

	if _, ok := m.Headers["User-Agent"]; !ok {
		m.Headers["User-Agent"] = UA_Browser
	}
	if _, ok := m.Headers["Authorization"]; !ok {
		m.Headers["Authorization"] = "Basic NTQ3MzM0NDgtYTU3Yi00MjU2LWE4MTEtMzdlYzNkNjJmM2E0Ok90QzR3elJRR2hLQ01sSDc2VEoy"
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
	r := make(map[string]interface{})
	err = json.Unmarshal(resp.Body(), &r)
	if err != nil {
		m.Logger.Errorln(err)
		result.Failed(err)
		return
	}

	if rr, ok := r["code"]; ok {
		if rr.(string) == "0" {
			result.Yes()
		} else if rr.(string) == "100016" {
			result.No()
		} else {
			result.Unexpected(fmt.Sprintf("code: %+v", rr))
		}
	} else {
		result.Unexpected(fmt.Sprintf(`key "code" not found`))
	}

	return
}
