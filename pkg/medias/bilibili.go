package medias

import (
	"crypto/md5"
	"fmt"

	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
)

func CheckBilibiliTW(m *Media) *CheckResult {
	m.Logger.Infoln("running")
	result := &CheckResult{Media: m.Name, Region: m.Region}

	s := randomSession()
	m.URL = fmt.Sprintf("https://api.bilibili.com/pgc/player/web/playurl?avid=50762638&cid=100279344&qn=0&type=&otype=json&ep_id=268176&fourk=1&fnver=0&fnval=16&session=%s&module=bangumi", s)

	checkBilibili(m, result)
	return result
}

func CheckBilibiliHKMCTW(m *Media) *CheckResult {
	m.Logger.Infoln("running")
	result := &CheckResult{Media: m.Name, Region: "HK_MC_TW"}

	s := randomSession()
	m.URL = fmt.Sprintf("https://api.bilibili.com/pgc/player/web/playurl?avid=18281381&cid=29892777&qn=0&type=&otype=json&ep_id=183799&fourk=1&fnver=0&fnval=16&session=%s&module=bangumi", s)

	checkBilibili(m, result)
	return result
}

func randomSession() string {
	u := uuid.New().String()
	d := md5.New()
	d.Write([]byte(u))
	return fmt.Sprintf("%x", d.Sum(nil))
}

func checkBilibili(m *Media, result *CheckResult) {
	if _, ok := m.Headers["User-Agent"]; !ok {
		m.Headers["User-Agent"] = UA_Browser
	}

	resp, err := m.Do()
	if err != nil {
		m.Logger.Errorln(err)
		result.Failed(err)
		return
	}
	defer fasthttp.ReleaseResponse(resp)

	switch resp.StatusCode() {
	case fasthttp.StatusOK:
		result.Yes()
	case fasthttp.StatusForbidden:
		result.No()
	default:
		result.UnexpectedStatusCode(resp.StatusCode())
	}
}
