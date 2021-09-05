package medias

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"text/tabwriter"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

const (
	CheckResultYes           = "yes"
	CheckResultNo            = "no"
	CheckResultUnexpected    = "unexpected"
	CheckResultFailed        = "failed"
	CheckResultOverseaOnly   = "oversea_only"
	CheckResultOriginalsOnly = "originals_only"

	printPadding = 48

	ContentTypeJSON = "application/json"
)

type CheckResult struct {
	Result  string
	Media   string
	Region  string
	Type    string
	Message string
}

func (c *CheckResult) Yes(intr ...interface{}) {
	c.Result = CheckResultYes
	if len(intr) > 0 {
		c.Message = fmt.Sprint(intr...)
	}
}

func (c *CheckResult) No() {
	c.Result = CheckResultNo
}

func (c *CheckResult) Oversea() {
	c.Result = CheckResultOverseaOnly
}

func (c *CheckResult) OriginalsOnly() {
	c.Result = CheckResultOriginalsOnly
}

func (c *CheckResult) Unexpected(msg interface{}) {
	c.Result = CheckResultUnexpected
	c.Message = fmt.Sprint(msg)
}

func (c *CheckResult) UnexpectedStatusCode(code interface{}) {
	c.Result = CheckResultUnexpected
	c.Message = fmt.Sprintf("status code: %s", fmt.Sprint(code))
}

func (c *CheckResult) Failed(msg interface{}) {
	c.Result = CheckResultFailed
	c.Message = fmt.Sprint(msg)
}

type CheckResultSlice []*CheckResult

func (c *CheckResultSlice) Len() int {
	return len(*c)
}

func (c *CheckResultSlice) Swap(i, j int) {
	(*c)[i], (*c)[j] = (*c)[j], (*c)[i]
}

func (c *CheckResultSlice) Less(i, j int) bool {
	if (*c)[i].Region < (*c)[j].Region {
		return true
	} else if (*c)[i].Region == (*c)[j].Region {
		if (*c)[i].Type < (*c)[j].Type {
			return true
		} else if (*c)[i].Type == (*c)[j].Type {
			return (*c)[i].Media < (*c)[j].Media
		}
	}
	return false
}

func (c *CheckResultSlice) PrintTo(writer io.Writer) {
	w := tabwriter.NewWriter(writer, 8, 8, 0, ' ', 0)
	lastRegion := ""
	lastOttType := ""
	for _, res := range *c {
		if lastRegion != res.Region {
			w.Flush()
			fmt.Fprintf(writer, "\n==========[ %s ]==========\n", res.Region)
			lastRegion = res.Region
		}
		if lastOttType != res.Type {
			lastOttType = res.Type
			if lastOttType != "" {
				w.Flush()
				fmt.Fprintf(writer, "\n------< %s - %s >------\n", res.Region, res.Type)
			}
		}

		s := HumanReadableNames[res.Media]
		pad := printPadding - len(s)
		for i := 0; i < pad; i++ {
			s += " "
		}
		s += "\t"
		s += strings.ToUpper(res.Result)

		if res.Message != "" {
			s += fmt.Sprintf(" (%s)", res.Message)
		}
		fmt.Fprintln(w, s)
	}
	w.Flush()
}

type Media struct {
	Enabled  bool              `json:"enabled"`
	URL      string            `json:"-"`
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
	dialer := fasthttp.TCPDialer{
		Resolver: &net.Resolver{
			PreferGo:     true,
			StrictErrors: false,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				d := net.Dialer{}
				m.Logger.Debugln("connecting to dns")
				return d.DialContext(ctx, network, m.DNS)
			},
		},
	}
	client := fasthttp.Client{Dial: dialer.Dial}
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
		"dns":         m.DNS,
	})

	resp := fasthttp.AcquireResponse()
	if err := client.DoDeadline(req, resp, time.Now().Add(time.Duration(m.Timeout)*time.Second)); err != nil {
		fasthttp.ReleaseResponse(resp)
		return nil, err
	}
	m.Logger = m.Logger.WithField("status_code", resp.StatusCode())
	return resp, nil
}

func (m *Media) DoRedirects() (*fasthttp.Response, error) {
	cnt := 0
	for {
		resp, err := m.Do()
		if err != nil {
			return nil, err
		}
		status := resp.StatusCode()

		if status == fasthttp.StatusFound || status == fasthttp.StatusMovedPermanently {
			cnt += 1
			if cnt > 50 {
				return nil, errors.New("too many redirects")
			}
			m.URL = string(resp.Header.Peek("location"))
			fasthttp.ReleaseResponse(resp)
			continue
		} else {
			return resp, nil
		}
	}
}
