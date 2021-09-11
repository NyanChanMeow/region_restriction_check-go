package medias

import (
	"strings"

	"github.com/valyala/fasthttp"
)

func CheckViuCom(m *Media) (result *CheckResult) {
	m.URL = "https://www.viu.com/"
	m.Logger.Infoln("running")

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

	redirUrl := string(resp.Header.Peek("location"))

	if a := strings.Split(redirUrl, "/"); len(a) >= 5 {
		if a[4] == "no-service" {
			result.No()
		} else {
			result.Yes("Region:", strings.ToUpper(a[4]))
		}
	} else {
		result.Unexpected("URL:", redirUrl)
	}

	return
}
