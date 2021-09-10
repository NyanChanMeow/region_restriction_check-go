package exporter

import (
	"net/http"

	"github.com/NyanChanMeow/region_restriction_check-go/pkg/medias"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func HandleStatusUpdate(result *medias.CheckResult) {
	rrcStatus.WithLabelValues(
		result.Media,
		medias.HumanReadableNames[result.Media],
		result.Task,
	).Set(float64(StatusMapping[result.Result]))
}

func ServeExporter(listen string) error {
	http.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(listen, nil)
}
