package flags

import (
	"context"
	"goflagsmith/internal/domain"
	"goflagsmith/internal/state"
	"time"
)

type Reader interface {
	IsFeatureEnabled(ctx context.Context, featureName string) bool
	GetJSONConfig(ctx context.Context, configName string) (string, error)
	GetCanaryRules(ctx context.Context) (*domain.CanaryRoutingRules, error)
}

type Service interface {
	Reader

	MonitorFlagsReady(
		ctx context.Context, appState *state.State, interval time.Duration,
	)
	StartRulesSync(
		ctx context.Context, interval time.Duration,
	)
}
