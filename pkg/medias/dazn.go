package medias

import (
	"encoding/json"
	"strings"

	"github.com/valyala/fasthttp"
)

func CheckDazn(m *Media) (result *CheckResult) {
	m.URL = "https://startup.core.indazn.com/misl/v5/Startup"
	m.Method = "POST"
	m.Headers[fasthttp.HeaderContentType] = ContentTypeJSON
	m.Logger.Infoln("running")

	if _, ok := m.Headers["User-Agent"]; !ok {
		m.Headers["User-Agent"] = UA_Browser
	}
	if _, ok := m.Headers[fasthttp.HeaderContentType]; !ok {
		m.Headers[fasthttp.HeaderContentType] = ContentTypeJSON
	}
	if m.Body == "" {
		m.Body = `{"LandingPageKey":"generic","Languages":"zh-CN,zh,en","Platform":"web","PlatformAttributes":{},"Manufacturer":"","PromoCode":"","Version":"2"}`
	}
	result = &CheckResult{Media: m.Name, Region: m.Region}

	resp, err := m.Do()
	if err != nil {
		m.Logger.Errorln(err)
		result.Failed(err.Error())
		return
	}
	defer fasthttp.ReleaseResponse(resp)

	if resp.StatusCode() != fasthttp.StatusOK {
		result.UnexpectedStatusCode(resp.StatusCode())
		return
	}

	var r struct {
		Region struct {
			IsAllowed bool   `json:"isAllowed"`
			Country   string `json:"GeolocatedCountry"`
		} `json:"Region"`
	}
	err = json.Unmarshal(resp.Body(), &r)
	if err != nil {
		m.Logger.Errorln(err)
		result.Failed(err)
		return
	}

	if r.Region.IsAllowed {
		result.Yes("Region:", strings.ToUpper(r.Region.Country))
	} else {
		result.No()
	}

	return
}
