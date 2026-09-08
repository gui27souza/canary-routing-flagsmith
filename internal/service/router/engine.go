package router

import (
	"context"
	"goflagsmith/internal/domain"
	"goflagsmith/internal/service/flags"
	"goflagsmith/internal/state"
	"strings"
	"time"
)

type BucketCalculator func(expression string) int
type Now func() time.Time

// Engine is responsible for making traffic routing decisions
// based on feature flags and dynamic canary rules.
type Engine struct {
	fr  flags.Reader
	s   *state.State
	bc  BucketCalculator
	now Now
}

func NewEngine(
	fr flags.Reader, s *state.State, bc BucketCalculator, now Now,
) *Engine {
	return &Engine{
		fr: fr, s: s, bc: bc, now: now,
	}
}

// Route evaluates the user context and deterministically decides whether the
// traffic should be routed to the stable version ("v1") or the canary version ("v2").
func (e *Engine) Route(
	ctx context.Context, uc domain.UserContext,
) domain.RouteDecision {

	if !e.fr.IsFeatureEnabled(ctx, "enable_v2_routing") {
		return e.deniedRouting("v2 routing not enabled")
	}

	rules, err := e.fr.GetCanaryRules(ctx)
	if err != nil {
		return e.deniedRouting("internal error on canary rules fetching")
	}

	var allowCountry bool
	for _, country := range rules.Countries {
		if strings.EqualFold(country, uc.Country) {
			allowCountry = true
			break
		}
	}
	if !allowCountry {
		return e.deniedRouting("unavailable for user country")
	}

	usrBucket := e.bc(uc.UserID)

	target := "v1"
	if usrBucket < rules.Percentage {
		target = "v2"
	}

	return domain.RouteDecision{
		Target:    target,
		Reason:    "canary sorting rules",
		Telemetry: e.getTelemetryData(),
	}
}

func (e *Engine) getTelemetryData() domain.TelemetryData {
	snapshot := e.s.Snapshot()
	return domain.TelemetryData{
		EvaluatedAt:   e.now(),
		CacheHydrated: snapshot.Features,
	}
}

func (e *Engine) deniedRouting(msg string) domain.RouteDecision {
	return domain.RouteDecision{
		Target:    "v1",
		Reason:    msg,
		Telemetry: e.getTelemetryData(),
	}
}
