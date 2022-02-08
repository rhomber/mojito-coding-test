package boot

import (
	"mojito-coding-test/common/config"
	"mojito-coding-test/common/core"
)

func Cfg(verbose bool) *config.Config {
	cfg := core.Config
	cfg.Set("verbose", verbose)

	return cfg
}
