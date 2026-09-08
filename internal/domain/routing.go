package domain

import "time"

// CanaryRoutingRules defines the dynamic configuration structure
// retrieved via Remote Config.
type CanaryRoutingRules struct {
	Percentage int      `json:"canary_percentage" binding:"required"`
	Countries  []string `json:"allowed_countries" binding:"required"`
}

func NewCanaryRules(percentage int, coutries []string) *CanaryRoutingRules {
	return &CanaryRoutingRules{percentage, coutries}
}

// UserContext defines the client context data required
// to make a routing decision.
type UserContext struct {
	UserID     string `json:"user_id" binding:"required"`
	Country    string `json:"country" binding:"required"`
	AppVersion string `json:"app_version" binding:"omitempty"`
}

// TelemetryData encapsulates telemetry metadata about the routing evaluation.
type TelemetryData struct {
	EvaluatedAt   time.Time `json:"evaluated_at" binding:"required"`
	CacheHydrated bool      `json:"cache_hydrated" binding:"omitempty"`
}

// RouteDecision represents the data containing the final routing decision.
type RouteDecision struct {
	Target    string        `json:"target" binding:"required"`
	Reason    string        `json:"reason" binding:"required"`
	Telemetry TelemetryData `json:"telemetry" binding:"required"`
}
