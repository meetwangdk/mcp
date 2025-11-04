package biz

import (
	"github.com/google/wire"
)

// Set is biz providers.
var Set = wire.NewSet(
	NewDemoUseCase,
)
