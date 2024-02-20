package version

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/version"
)

func NewVersionCollector(program string, ver string) prometheus.Collector {
	return prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Namespace: program,
			Name:      "build_info",
			Help: fmt.Sprintf(
				"A metric with a constant '1' value labeled by version, revision, branch, goversion from which %s was build, and the goos and goarch for the build.", program,
			),
			ConstLabels: prometheus.Labels{
				"version":   ver,
				"revision":  version.GetRevision(),
				"branch":    version.Branch,
				"goversion": version.GoVersion,
				"goos":      version.GoOS,
				"goarch":    version.GoArch,
				"tags":      version.GetTags(),
			},
		},
		func() float64 { return 1 },
	)
}
