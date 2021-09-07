package exporter

import (
	"github.com/NyanChanMeow/region_restriction_check-go/pkg/medias"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	StatusMapping = map[string]int{
		medias.CheckResultYes:           10,
		medias.CheckResultNo:            0,
		medias.CheckResultOriginalsOnly: -1,
		medias.CheckResultOverseaOnly:   -2,
		medias.CheckResultUnexpected:    -3,
		medias.CheckResultFailed:        -4,
	}

	rrcStatus = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "rrc_unblock_status",
		Help: "Region Restriction Check Status",
	}, []string{"region", "media", "media_readable", "task"})
)
