package boot

import (
	"github.com/sirupsen/logrus"
	"mojito-coding-test/common/core"
)

func Logger() func() {
	closeFn := func() {}

	if "json" == core.Config.GetString("logger.format") {
		core.Logrus.Formatter = &logrus.JSONFormatter{}
	} else {
		core.Logrus.Formatter = &logrus.TextFormatter{
			ForceColors: true,
		}
	}
	if core.Config.GetBool("verbose") {
		core.Logrus.SetLevel(logrus.TraceLevel)
	} else {
		core.Logrus.SetLevel(logrus.InfoLevel)
	}

	// Logging fields

	core.Logger = core.Logrus.WithFields(logrus.Fields{
		"env": core.Config.GetEnv(),
	})

	// Kubernetes

	if podName := core.Config.GetString("pod.name"); podName != "" {
		core.Logger = core.Logger.WithFields(logrus.Fields{
			"pod_name":      podName,
			"pod_namespace": core.Config.GetString("pod.namespace"),
			"pod_ip":        core.Config.GetString("pod.ip"),
			"node_name":     core.Config.GetString("node.name"),
		})
	}

	return closeFn
}
