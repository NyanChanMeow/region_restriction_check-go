package medias

import (
	"encoding/base64"
	"encoding/json"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

const (
	CheckResultYes         = "yes"
	CheckResultNo          = "no"
	CheckResultUnexpected  = "unexpected"
	CheckResultFailed      = "failed"
	CheckResultOverseaOnly = "oversea_only"
)

type CheckResult struct {
	Result string
	Media  string
	Region string
	Error  error
}

type Media struct {
	Enabled  bool              `json:"enabled"`
	URL      string            `json:"url"`
	Method   string            `json:"method"`
	Headers  map[string]string `json:"headers"`
	Body     string            `json:"body"`
	DNS      string            `json:"dns"`
	Timeout  int               `json:"timeout"`
	Interval int               `json:"interval"`
	Name     string            `json:"-"`
	Region   string            `json:"-"`
	Logger   *log.Entry        `json:"-"`
}

func NewMediaConf() *Media {
	m := Media{}
	m.Headers = make(map[string]string)
	return &m
}

func (m *Media) UnmarshalJSON(data []byte) error {
	m.Enabled = true

	var result map[string]json.RawMessage
	err := json.Unmarshal(data, &result)
	if err != nil {
		return err
	}
	m.Headers = make(map[string]string)

	for k, v := range result {
		switch k {
		case "enabled":
			err = json.Unmarshal(v, &m.Enabled)
		case "url":
			err = json.Unmarshal(v, &m.URL)
		case "method":
			err = json.Unmarshal(v, &m.Method)
		case "body":
			err = json.Unmarshal(v, &m.Body)
		case "timeout":
			err = json.Unmarshal(v, &m.Timeout)
		case "headers":
			err = json.Unmarshal(v, &m.Headers)
		case "interval":
			err = json.Unmarshal(v, &m.Interval)
		}
		if err != nil {
			return err
		}
	}

	if m.Timeout == 0 {
		m.Timeout = 10
	}
	if m.Method == "" {
		m.Method = "GET"
	}
	return nil
}

func (m *Media) Do() (*fasthttp.Response, error) {
	client := fasthttp.Client{}
	client.NoDefaultUserAgentHeader = true

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.SetRequestURI(m.URL)
	req.Header.SetMethod(m.Method)
	for k, v := range m.Headers {
		req.Header.Set(k, v)
	}
	req.SetBodyString(m.Body)

	m.Logger = m.Logger.WithFields(log.Fields{
		"url":         string(req.URI().FullURI()),
		"method":      string(req.Header.Method()),
		"req_body":    base64.StdEncoding.EncodeToString(req.Body()),
		"user_agent":  string(req.Header.UserAgent()),
		"timeout":     m.Timeout,
		"status_code": 0,
	})

	resp := fasthttp.AcquireResponse()
	if err := client.DoDeadline(req, resp, time.Now().Add(time.Duration(m.Timeout)*time.Second)); err != nil {
		fasthttp.ReleaseResponse(resp)
		return nil, err
	}
	return resp, nil
}
