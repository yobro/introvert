package probe

import (
	"context"
)

// Probe prober interface
type Probe interface {
	// Run starts a probe
	Run(ctx context.Context) error
}
