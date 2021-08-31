package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"time"

	"github.com/NyanChanMeow/region_restriction_check-go/pkg/medias"
	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	version     = "0.0.1"
	modeChecker = "checker"
	modeMonitor = "monitor"
)

var (
	flags struct {
		ConfigFile string `json:"-"`
		Version    bool   `json:"-"`
		// Checker
		Mode    string      `json:"-"`
		Regions regionArray `json:"-"`

		// Monitor
		Interval int             `json:"interval"`
		Log      Log             `json:"log"`
		Tasks    map[string]Task `json:"tasks"`
	}
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.StringVar(&flags.ConfigFile, "config.file", "config.json", "only working with monitor mode")
	flag.StringVar(&flags.Mode, "mode", modeChecker, "[checker, monitor]")
	flag.Var(&flags.Regions, "region", "available regions: [all, JP]")
	flag.BoolVar(&flags.Version, "version", false, "display version and exit")
	flag.Parse()
}

func main() {
	log.SetLevel(log.InfoLevel)
	// Timestamp formatted as RFC3339Nano
	log.SetFormatter(&log.TextFormatter{
		DisableColors:   true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05.000000000Z07:00",
	})
	log.SetOutput(os.Stdout)

	log.WithField("version", version).Infoln("Region Restriction Check")
	if flags.Version {
		return
	}

	switch flags.Mode {
	case modeChecker:
		runChecker()
	case modeMonitor:
		runMonitor()
	default:
		log.WithField("mode", flags.Mode).Fatalln("mode not found")
	}
}

func runChecker() {
	flags.Regions.CheckAll()
	result := make(chan *medias.CheckResult)

	cnt := 0
	for _, region := range flags.Regions {
		if len(region) == 0 {
			continue
		}

		if mediaFuncsRegion, ok := medias.MediaFuncs[region]; ok {
			for mediaName, mediaFunc := range mediaFuncsRegion {
				cnt++
				go func(mn, rg string, f func(*medias.Media) *medias.CheckResult) {
					mc := medias.NewMediaConf()
					mc.Name = mn
					mc.Region = rg
					mc.Logger = log.WithFields(log.Fields{
						"region": rg,
						"media":  mn,
					})
					mc.Timeout = 10
					result <- f(mc)
				}(mediaName, region, mediaFunc)
			}
		} else {
			log.WithField("region", region).Errorln("region not found")
		}
	}

	for ; cnt > 0; cnt-- {
		res := <-result
		fmt.Printf("%+v\n", res)
	}
}

func runMonitor() {
	buf, err := ioutil.ReadFile(flags.ConfigFile)
	if err != nil {
		log.Fatalln(err)
	}

	err = json.Unmarshal(buf, &flags)
	if err != nil {
		log.Fatalln(err)
	}

	log.WithField("filename", flags.Log.Filename).Infoln("redirect log to file")
	log.SetLevel(flags.Log.ParseLevel())
	log.SetOutput(&lumberjack.Logger{
		Filename:   flags.Log.Filename,
		MaxSize:    flags.Log.MaxSize,
		MaxBackups: flags.Log.MaxBackups,
		MaxAge:     flags.Log.MaxAge,
		Compress:   false,
	})
	log.WithField("version", version).Infoln("Region Restriction Check")

	result := make(chan *medias.CheckResult)

	for taskName, task := range flags.Tasks {
		if !task.Enabled {
			continue
		}

		for mediaName, mediaConf := range task.Medias {
			if !mediaConf.Enabled {
				continue
			}
			if mediaConf.Interval == 0 {
				mediaConf.Interval = task.Interval
			}
			mediaConf.Name = mediaName
			mlog := log.WithFields(log.Fields{
				"task":  taskName,
				"media": mediaName,
			})

			found := false
			for region, mediaFuncsRegion := range medias.MediaFuncs {
				mlog = mlog.WithFields(log.Fields{
					"region": region,
				})

				if mediaFunc, ok := mediaFuncsRegion[mediaName]; ok {
					go func(mc *medias.Media, logger *log.Entry) {
						mc.Logger = logger
						for {
							result <- mediaFunc(mc)
							logger.WithField("interval", mc.Interval).Infoln("waiting")
							time.Sleep(time.Duration(mc.Interval) * time.Second)
						}
					}(mediaConf, mlog)
					found = true
					break
				}
			}
			if !found {
				mlog.Errorln("media not found")
			}
		}
	}

	for {
		res := <-result
		fmt.Printf("%+v\n", res)
	}
}