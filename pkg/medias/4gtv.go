package medias

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

func Check4GTV(m *Media) *CheckResult {
	m.Method = "POST"
	m.Headers[fasthttp.HeaderContentType] = ContentTypeJSON
	m.Logger.Infoln("running")
	if m.URL == "" {
		m.URL = "https://api2.4gtv.tv//Vod/GetVodUrl3"
	}
	if _, ok := m.Headers["User-Agent"]; !ok {
		m.Headers["User-Agent"] = UA_Browser
	}
	if m.Body == "" {
		m.Body = `value=D33jXJ0JVFkBqV%2BZSi1mhPltbejAbPYbDnyI9hmfqjKaQwRQdj7ZKZRAdb16%2FRUrE8vGXLFfNKBLKJv%2BfDSiD%2BZJlUa5Msps2P4IWuTrUP1%2BCnS255YfRadf%2BKLUhIPj`
	}
	result := &CheckResult{Media: m.Name, Region: m.Region}

	resp, err := m.Do()
	if err != nil {
		m.Logger.Errorln(err)
		result.Failed(err.Error())
		return result
	}
	defer fasthttp.ReleaseResponse(resp)

	if resp.StatusCode() == fasthttp.StatusAccepted {

		r := make(map[string]interface{})
		err = json.Unmarshal(resp.Body(), &r)
		if err != nil {
			m.Logger.Errorln(err)
			result.Failed(err)
			return result
		}

		if rr, ok := r["Success"]; ok {
			if rr.(bool) == true {
				result.Yes()
			} else if rr.(bool) == false {
				result.No()
			} else {
				result.Unexpected(fmt.Sprintf("%+v", rr))
			}
		} else {
			result.Failed(`key "Success" not found`)
		}
	} else {
		result.UnexpectedStatusCode(resp.StatusCode())
	}

	m.Logger.WithFields(log.Fields{
		"status_code": resp.StatusCode(),
		"result":      result.Result,
		"message":     result.Message,
	}).Infoln("done")
	return result
}
