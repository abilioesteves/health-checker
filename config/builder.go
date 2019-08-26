package config

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	port            = "port"
	logLevel        = "log-level"
	targetHealthURL = "target-health-url"
	targetName      = "target-name"
)

// Flags define the fields that will be passed via cmd
type Flags struct {
	Port            string
	LogLevel        string
	TargetHealthURL string
	TargetName      string
}

// Builder defines the parametric information of a whisper server instance
type Builder struct {
	*Flags
	HealthMetric *prometheus.GaugeVec
}

// AddFlags adds flags for Builder.
func AddFlags(flags *pflag.FlagSet) {
	flags.String(port, "37441", "[optional] Custom port for accessing Whisper's services. Defaults to 7070")
	flags.String(logLevel, "info", "[optional] Sets the Log Level to one of seven (trace, debug, info, warn, error, fatal, panic). Defaults to info")
	flags.String(targetHealthURL, "", "Determines the url for the health-checker to consume for health statistics")
	flags.String(targetName, "", "Determines the name of the target that this health-checker instance will be watching")
}

// InitFromViper initializes the web server builder with properties retrieved from Viper.
func (b *Builder) InitFromViper(v *viper.Viper) *Builder {
	flags := new(Flags)
	flags.Port = v.GetString(port)
	flags.LogLevel = v.GetString(logLevel)
	flags.TargetHealthURL = v.GetString(targetHealthURL)
	flags.TargetName = v.GetString(targetName)

	flags.check()

	b.Flags = flags
	b.HealthMetric = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "dependency_up_or_down",
		Help: "Records the status of a dependency",
	}, []string{"servicename", "dependency", "err"})
	prometheus.MustRegister(b.HealthMetric)

	logrus.Infof("Flags: '%v'", b.Flags)
	return b
}

func (flags *Flags) check() {

	if flags.TargetHealthURL == "" && flags.TargetName == "" {
		panic("target-health-url nd target-name cannot be empty")
	}

}
