package flags_test

import (
	"context"
	"errors"
	"goflagsmith/internal/service/flags"
	"goflagsmith/internal/state"
	"testing"
	"time"

	"github.com/Flagsmith/flagsmith-go-client/v4"
)

type mockFlagsmithSdk struct {
	flags      []flagsmith.Flag
	shouldFail bool
}

func (m *mockFlagsmithSdk) GetEnvironmentFlags(ctx context.Context) (flagsmith.Flags, error) {
	if m.shouldFail {
		return flagsmith.Flags{}, errors.New("network error")
	}

	return flagsmith.Flags{}, nil
}

func TestNewClient_EmptyApiKey(t *testing.T) {
	_, err := flags.NewClient(t.Context(), "")
	if err == nil {
		t.Errorf("Expected error when apiKey is empty")
	}
}

func TestMonitorFlagsReady(t *testing.T) {

	mockSdk := &mockFlagsmithSdk{shouldFail: false}
	client := flags.NewClientWithSDK(mockSdk)

	appState := state.NewState()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client.MonitorFlagsReady(ctx, appState, 100*time.Millisecond)

	time.Sleep(300 * time.Millisecond)

	if !appState.Snapshot().Features {
		t.Errorf("Expected Features to be true after monitor execution")
	}
}

func TestIsFeatureEnabled_GracefulFallback(t *testing.T) {
	mockSdk := &mockFlagsmithSdk{shouldFail: true}
	client := flags.NewClientWithSDK(mockSdk)

	enabled := client.IsFeatureEnabled(context.Background(), "test_flag")
	if enabled {
		t.Errorf("Expected false (graceful fallback) on SDK error, got true")
	}
}
