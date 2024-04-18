package cmd

import (
	cli "github.com/Foxcapades/Argonaut"
	"github.com/Foxcapades/Argonaut/pkg/argo"
	"github.com/sirupsen/logrus"
)

type Opts struct {
	LogLevel logrus.Level
}

func RegisterOpts(builder argo.CommandTreeBuilder, opts *Opts) {
	opts.LogLevel = logrus.WarnLevel

	builder.WithFlag(cli.ShortFlag('V').
		WithDescription("Verbose logging.  Specify multiple times to enable more granular logging.\n1x INFO, 2x DEBUG, 3x TRACE.").
		WithCallback(func(flag argo.Flag) {
			switch flag.HitCount() {
			case 1:
				opts.LogLevel = logrus.InfoLevel
			case 2:
				opts.LogLevel = logrus.DebugLevel
			default:
				opts.LogLevel = logrus.TraceLevel
			}
		}))
}
