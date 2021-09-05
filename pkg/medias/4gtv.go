package medias

import (
	"encoding/json"
	"fmt"

	"github.com/valyala/fasthttp"
)

func Check4GTV(m *Media) (result *CheckResult) {
	m.Method = "POST"
	m.URL = "https://api2.4gtv.tv//Vod/GetVodUrl3"
	m.Headers[fasthttp.HeaderContentType] = ContentTypeJSON
	m.Logger.Infoln("running")

	if _, ok := m.Headers["User-Agent"]; !ok {
		m.Headers["User-Agent"] = UA_Browser
	}
	if m.Body == "" {
		m.Body = `value=D33jXJ0JVFkBqV%2BZSi1mhPltbejAbPYbDnyI9hmfqjKaQwRQdj7ZKZRAdb16%2FRUrE8vGXLFfNKBLKJv%2BfDSiD%2BZJlUa5Msps2P4IWuTrUP1%2BCnS255YfRadf%2BKLUhIPj`
	}
	result = &CheckResult{Media: m.Name, Region: m.Region}

	resp, err := m.Do()
	if err != nil {
		m.Logger.Errorln(err)
		result.Failed(err.Error())
		return
	}
	defer fasthttp.ReleaseResponse(resp)

	if resp.StatusCode() != fasthttp.StatusAccepted {
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

	return
}
