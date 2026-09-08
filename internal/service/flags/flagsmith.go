package flags

import (
	"context"
	"encoding/json"
	"errors"
	"goflagsmith/internal/domain"
	"goflagsmith/internal/state"
	"log"
	"sync/atomic"
	"time"

	"github.com/Flagsmith/flagsmith-go-client/v4"
)

// flagsmithSDK defines the minimum interface required from the official
// Flagsmith SDK to support deterministic mocking during unit testing.
type flagsmithSDK interface {
	GetEnvironmentFlags(ctx context.Context) (flagsmith.Flags, error)
}

// Client implements the Service interface using the Flagsmith SDK.
// It encapsulates background polling and local evaluation logic.
type Client struct {
	sdk   flagsmithSDK
	rules atomic.Pointer[domain.CanaryRoutingRules]
}

// NewClient initializes and returns a concrete Client as a Service,
// bootstrapping local evaluation and background polling configurations.
func NewClient(ctx context.Context, apiKey string) (Service, error) {

	if apiKey == "" {
		return nil, errors.New("FLAGSMITH_API_KEY environment variable is required")
	}

	sdk := flagsmith.NewClient(
		apiKey,
		flagsmith.WithLocalEvaluation(ctx),
		flagsmith.WithEnvironmentRefreshInterval(60*time.Second),
		flagsmith.WithRequestTimeout(5*time.Second),
	)

	return &Client{sdk: sdk}, nil
}

// Testable Constructor
func NewClientWithSDK(sdk flagsmithSDK) Service {
	return &Client{sdk: sdk}
}

// MonitorFlagsReady starts an asynchronous background goroutine that polls
// the Flagsmith SDK until the first successful cache hydration occurs,
// updating the global application state once completed.
func (c *Client) MonitorFlagsReady(
	ctx context.Context, appState *state.State, interval time.Duration,
) {

	go func() {

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if appState.Snapshot().Features {
					return
				}

				ctxCheck, cancel := context.WithTimeout(ctx, 1*time.Second)
				_, err := c.sdk.GetEnvironmentFlags(ctxCheck)
				cancel()

				if err == nil {
					appState.SetFeaturesReady()
					log.Println("INFO: Flags cache ready, service is now fully ready for traffic")
					return
				}
			}
		}
	}()
}

// IsFeatureEnabled queries the locally cached environment configurations and
// returns whether a given boolean feature flag is currently active.
func (c *Client) IsFeatureEnabled(ctx context.Context, featureName string) bool {

	flags, err := c.sdk.GetEnvironmentFlags(ctx)
	if err != nil {
		return false
	}

	enabled, err := flags.IsFeatureEnabled(featureName)
	if err != nil {
		return false
	}
	return enabled
}

// GetJSONConfig retrieves a dynamic Remote Config value from the local cache
// and returns it as a raw string. If the value is not a string, it returns
// an empty string safely without panicking.
func (c *Client) GetJSONConfig(ctx context.Context, configName string) (string, error) {

	flags, err := c.sdk.GetEnvironmentFlags(ctx)
	if err != nil {
		return "", err
	}

	value, err := flags.GetFeatureValue(configName)
	if err != nil {
		return "", err
	}

	strValue, ok := value.(string)
	if !ok {
		return "", nil
	}

	return strValue, nil
}

func (c *Client) GetCanaryRules(ctx context.Context) (*domain.CanaryRoutingRules, error) {

	currentRules := c.rules.Load()
	if currentRules == nil {
		return nil, errors.New("canary rules are empty")
	}

	return currentRules, nil
}

func (c *Client) syncRules(ctx context.Context) {

	// Timeout defensive ctx so the operation doesn't lock during execution
	ctxSync, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	rulesJSON, err := c.GetJSONConfig(ctxSync, "canary_routing_rules")
	if err != nil {
		// TODO - log fetching error on new rules
		return
	}
	if rulesJSON == "" {
		return
	}

	var rules domain.CanaryRoutingRules
	if err := json.Unmarshal([]byte(rulesJSON), &rules); err != nil {
		// TODO - log syntax error on new rules
		// Fallback - c.rules keeps last pointer
		return
	}

	// Atomic lock-free update
	c.rules.Store(&rules)
}

func (c *Client) StartRulesSync(
	ctx context.Context, interval time.Duration,
) {
	go func() {

		c.syncRules(ctx)

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				c.syncRules(ctx)
			}
		}
	}()
}
