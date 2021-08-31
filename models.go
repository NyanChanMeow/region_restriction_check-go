package main

import (
	"encoding/json"
	"sort"

	"github.com/NyanChanMeow/region_restriction_check-go/pkg/medias"

	log "github.com/sirupsen/logrus"
)

type Log struct {
	Level      string `json:"level"`
	Filename   string `json:"filename"`
	MaxSize    int    `json:"max_size"`
	MaxAge     int    `json:"max_age"`
	MaxBackups int    `json:"max_backups"`
}

func (l *Log) ParseLevel() log.Level {
	switch l.Level {
	case "trace":
		return log.TraceLevel
	case "debug":
		return log.DebugLevel
	case "info":
		return log.InfoLevel
	case "warning":
		return log.WarnLevel
	case "error":
		return log.ErrorLevel
	case "fatal":
		return log.FatalLevel
	default:
		return log.InfoLevel
	}
}

type Task struct {
	Enabled  bool                     `json:"enabled"`
	DNS      string                   `json:"dns"`
	Medias   map[string]*medias.Media `json:"medias"`
	Interval int                      `json:"interval"`
}

func (t *Task) UnmarshalJSON(data []byte) error {
	t.Enabled = true
	t.Interval = flags.Interval

	var result map[string]json.RawMessage
	err := json.Unmarshal(data, &result)
	if err != nil {
		return err
	}

	for k, v := range result {
		switch k {
		case "enabled":
			err = json.Unmarshal(v, &t.Enabled)
		case "dns":
			err = json.Unmarshal(v, &t.DNS)
		case "medias":
			err = json.Unmarshal(v, &t.Medias)
		case "interval":
			err = json.Unmarshal(v, &t.Interval)
		}
		if err != nil {
			return err
		}
	}

	if t.DNS == "" {
		t.DNS = "1.1.1.1:53"
	}

	for _, v := range t.Medias {
		if v.DNS == "" {
			v.DNS = t.DNS
		}
	}

	return nil
}

type regionArray []string

func (r *regionArray) String() string {
	return ""
}

func (r *regionArray) Set(v string) error {
	*r = append(*r, v)
	return nil
}

func (r *regionArray) CheckAll() {
	if len(*r) == 0 {
		*r = append(*r, "all")
	}
	sort.Strings(*r)
	if i := sort.SearchStrings(*r, "all"); i < len(*r) {
		var regions regionArray
		for k := range medias.MediaFuncs {
			regions = append(regions, k)
		}
		*r = regions
	}
}
