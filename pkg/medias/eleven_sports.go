package medias

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

func CheckElevenSports(m *Media) *CheckResult {
	m.Logger.Infoln("running")
	if m.URL == "" {
		m.URL = "https://apis.v-saas.com:9501/member/api/viewAuthorization?contentId=1&memberId=384030&menuId=3&platform=5&imei=c959b475-f846-4a86-8e9b-508048372508"
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

	if resp.StatusCode() != fasthttp.StatusOK {
		result.UnexpectedStatusCode(resp.StatusCode())
	} else {
		var r struct {
			Data struct {
				QQ             string `json:"qq"`
				ST             string `json:"st"`
				BoostStreamURL string `json:"boostStreamUrl"`
			} `json:"data"`
		}
		err = json.Unmarshal(resp.Body(), &r)
		if err == nil {
			oURL := m.URL
			defer func() { m.URL = oURL }()

			m.URL = r.Data.BoostStreamURL + fmt.Sprintf("?st=%s&qq=%s", r.Data.ST, r.Data.QQ)
			resp2, err := m.Do()
			if err != nil {
				m.Logger.Errorln(err)
				result.Failed(err)
			} else {
				defer fasthttp.ReleaseResponse(resp2)

				switch resp.StatusCode() {
				case fasthttp.StatusOK:
					result.Yes()
				case fasthttp.StatusForbidden:
					result.No()
				default:
					result.UnexpectedStatusCode(resp.StatusCode())
				}
			}
		} else {
			result.Failed(err)
		}
	}

	m.Logger.WithFields(log.Fields{
		"status_code": resp.StatusCode(),
		"result":      result.Result,
		"message":     result.Message,
	}).Infoln("done")
	return result
}
