package metrics

import (
	"runtime/debug"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/version"
)

// SetBuildInfo sets the build information.
func SetBuildInfo(v, revision, branch, buildUser, buildDate string) {
	version.Version = v
	version.Revision = revision
	version.Branch = branch
	version.BuildUser = buildUser
	version.BuildDate = buildDate
	PromRegistry.MustRegister(version.NewCollector("easyai"))
}

var (
	metricEasyaiVersion = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "easyai",
		Subsystem: "easyai",
		Name:      "version",
		Help:      "easyai platform version.",
	}, []string{"name"})
)

// nolint: gochecknoinits
func init() {
	var version string
	if version = buildVersion("easyai-platform"); version == "" {
		version = buildVersion("github.com/easyai-io/easyai-platform")
	}
	Register(metricEasyaiVersion)
	metricEasyaiVersion.WithLabelValues(version).Inc()
}

func buildVersion(path string) string {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		return ""
	}
	for _, d := range buildInfo.Deps {
		if d.Path == path {
			if d.Replace != nil {
				return d.Replace.Version
			}
			return d.Version
		}
	}
	return ""
}
