package common

import "go.uber.org/fx"

type Params struct {
	fx.In

	enableFeature1 bool
	enableFeature2 bool
}
