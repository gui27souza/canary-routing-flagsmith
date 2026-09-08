package testutil

import (
	"context"
	"goflagsmith/internal/domain"
	"goflagsmith/internal/state"
	"time"
)

type MockReader struct {
	FeatureEnabled bool
	JSONConfig     string
	ErrJSONConfig  error
	CanaryRules    *domain.CanaryRoutingRules
	ErrCanaryRules error
}

func (m *MockReader) IsFeatureEnabled(ctx context.Context, featureName string) bool {
	return m.FeatureEnabled
}
func (m *MockReader) GetJSONConfig(ctx context.Context, configName string) (string, error) {
	return m.JSONConfig, m.ErrJSONConfig
}
func (m *MockReader) GetCanaryRules(ctx context.Context) (*domain.CanaryRoutingRules, error) {
	return m.CanaryRules, m.ErrCanaryRules
}

func NewMockReader(
	featEnabled bool,
	jsonConfig string, errJson error,
	canaryRules *domain.CanaryRoutingRules, errCanaryRules error,
) *MockReader {
	return &MockReader{
		FeatureEnabled: featEnabled,
		JSONConfig:     jsonConfig,
		ErrJSONConfig:  errJson,
		CanaryRules:    canaryRules,
		ErrCanaryRules: errCanaryRules,
	}
}

type MockService struct {
	*MockReader
}

func NewMockService(reader *MockReader) *MockService {
	return &MockService{MockReader: reader}
}

func (m *MockService) MonitorFlagsReady(
	ctx context.Context, appState *state.State, interval time.Duration,
) {
}

func (m *MockService) StartRulesSync(
	ctx context.Context, interval time.Duration,
) {
}
