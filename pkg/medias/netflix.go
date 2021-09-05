package medias

import (
	"strings"

	"github.com/valyala/fasthttp"
)

func CheckNetflix(m *Media) (result *CheckResult) {
	m.Logger.Infoln("running")
	m.URL = "https://www.netflix.com/title/81215567"
	if _, ok := m.Headers["User-Agent"]; !ok {
		m.Headers["User-Agent"] = UA_Browser
	}
	result = &CheckResult{Media: m.Name, Region: m.Region}

	resp, err := m.DoRedirects()
	if err != nil {
		m.Logger.Errorln(err)
		result.Failed(err)
		return
	}
	defer fasthttp.ReleaseResponse(resp)

	switch resp.StatusCode() {
	case fasthttp.StatusNotFound:
		result.OriginalsOnly()
		return
	case fasthttp.StatusForbidden:
		result.No()
		return
	case fasthttp.StatusOK:
		break
	default:
		result.UnexpectedStatusCode(resp.StatusCode())
		return
	}

	m.URL = "https://www.netflix.com/title/80018499"
	resp2, err := m.Do()
	if err != nil {
		m.Logger.Errorln(err)
		result.Failed(err)
		return
	}
	defer fasthttp.ReleaseResponse(resp2)

	redirUrl := resp2.Header.Peek("location")

	result.Yes("Region: US")
	if a := strings.Split(string(redirUrl), "/"); len(a) >= 4 {
		if aa := strings.Split(a[3], "-"); len(aa) > 1 {
			result.Yes("Region:", strings.ToUpper(aa[0]))
		}
	}

	return
}
