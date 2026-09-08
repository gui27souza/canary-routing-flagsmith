package handlers_test

import (
	"encoding/json"
	"goflagsmith/internal/handlers"
	"goflagsmith/internal/state"
	"goflagsmith/internal/testutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupMockRouterReadyz(handler *handlers.AppHandler) *gin.Engine {
	return testutil.MockRouter(
		http.MethodGet, "/readyz", handler.Readyz,
	)
}

func TestReadyz_NotReady(t *testing.T) {

	appState := state.NewState()
	mockFlags := testutil.NewMockService(
		testutil.NewMockReader(true, "", nil, nil, nil),
	)

	h := handlers.NewAppHandler(appState, mockFlags)

	router := setupMockRouterReadyz(h)

	// Simulates a GET request to /readyz
	req, _ := http.NewRequest(http.MethodGet, "/readyz", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	// State was created, not ready yet
	if resp.Code != http.StatusServiceUnavailable {
		t.Errorf("expected status 503, got %d", resp.Code)
	}

	var body state.ReadinessStatus
	if err := json.Unmarshal(resp.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode response JSON: %v", err)
	}

	if body.Clients || body.Features {
		t.Errorf("expected clients and features to be false in payload")
	}
}

func TestReadyz_Ready(t *testing.T) {

	appState := state.NewState()

	// Sets State to be fully ready
	appState.SetClientsReady()
	appState.SetFeaturesReady()

	mockFlags := testutil.NewMockService(
		testutil.NewMockReader(true, "", nil, nil, nil),
	)
	h := handlers.NewAppHandler(appState, mockFlags)
	router := setupMockRouterReadyz(h)

	// Simulates a GET request to /readyz
	req, _ := http.NewRequest(http.MethodGet, "/readyz", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	// State fully ready
	if resp.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.Code)
	}
}
