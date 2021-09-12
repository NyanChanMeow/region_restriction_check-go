package medias

import (
	"strings"

	"github.com/valyala/fasthttp"
)

func CheckHBOMax(m *Media) (result *CheckResult) {
	m.URL = "https://www.hbomax.com/"
	m.Logger.Infoln("running")

	result = &CheckResult{Media: m.Name, Region: m.Region}

	resp, err := m.Do()
	if err != nil {
		m.Logger.Errorln(err)
		result.Failed(err)
		return
	}
	defer fasthttp.ReleaseResponse(resp)

	redirUrl := string(resp.Header.Peek("location"))
	if strings.Contains(redirUrl, "geo-availability") {
		result.No()
		return
	}
	result.Yes()

	if a := strings.Split(redirUrl, "/"); len(a) >= 4 {
		result.Yes("Region:", strings.ToUpper(a[3]))
	}

	return
}
