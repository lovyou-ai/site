package graph

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// TestParseCostDollars verifies cost extraction from post bodies.
func TestParseCostDollars(t *testing.T) {
	cases := []struct {
		body string
		want float64
	}{
		{"Builder shipped the feature. Cost: $0.53. Duration: 3m28s.", 0.53},
		{"No cost here.", 0},
		{"Multiple $1.00 and $2.50 — first wins.", 1.00},
		{"Cost: $0.00", 0},
	}
	for _, c := range cases {
		got := parseCostDollars(c.body)
		if got != c.want {
			t.Errorf("parseCostDollars(%q) = %.2f, want %.2f", c.body, got, c.want)
		}
	}
}

// TestParseDurationStr verifies duration extraction from post bodies.
func TestParseDurationStr(t *testing.T) {
	cases := []struct {
		body string
		want string
	}{
		{"Cost: $0.53. Duration: 3m28s.", "3m28s"},
		{"Duration: 12m", "12m"},
		{"No duration here.", ""},
		{"Duration: 0m5s.", "0m5s"},
	}
	for _, c := range cases {
		got := parseDurationStr(c.body)
		if got != c.want {
			t.Errorf("parseDurationStr(%q) = %q, want %q", c.body, got, c.want)
		}
	}
}

// TestComputeHiveStats verifies aggregate cost metrics.
func TestComputeHiveStats(t *testing.T) {
	posts := []Node{
		{Body: "Shipped X. Cost: $1.00. Duration: 2m."},
		{Body: "Shipped Y. Cost: $0.50."},
		{Body: "No cost here."},
	}
	stats := computeHiveStats(posts)
	if stats.Features != 2 {
		t.Errorf("Features = %d, want 2", stats.Features)
	}
	if stats.TotalCost != 1.50 {
		t.Errorf("TotalCost = %.2f, want 1.50", stats.TotalCost)
	}
	want := 0.75
	if stats.AvgCost != want {
		t.Errorf("AvgCost = %.2f, want %.2f", stats.AvgCost, want)
	}
}

// TestComputePipelineRoles verifies active/idle state and last-active timestamps.
func TestComputePipelineRoles(t *testing.T) {
	now := time.Now()
	recentPost := Node{
		Title:     "[hive:builder] iter 240: shipped feature",
		CreatedAt: now.Add(-5 * time.Minute),
	}
	oldPost := Node{
		Title:     "[hive:scout] iter 238: scouted gap",
		CreatedAt: now.Add(-2 * time.Hour),
	}

	roles := computePipelineRoles([]Node{recentPost, oldPost})

	roleByName := make(map[string]PipelineRole, len(roles))
	for _, r := range roles {
		roleByName[r.Name] = r
	}

	// Builder: recent post — should be active.
	builder, ok := roleByName["Builder"]
	if !ok {
		t.Fatal("Builder role missing from result")
	}
	if !builder.Active {
		t.Error("Builder should be Active (post within activeRoleThreshold)")
	}
	if builder.LastActive.IsZero() {
		t.Error("Builder LastActive should not be zero")
	}

	// Scout: old post — should be inactive.
	scout, ok := roleByName["Scout"]
	if !ok {
		t.Fatal("Scout role missing from result")
	}
	if scout.Active {
		t.Error("Scout should not be Active (post older than activeRoleThreshold)")
	}
	if scout.LastActive.IsZero() {
		t.Error("Scout LastActive should not be zero")
	}

	// Critic: no posts — should be idle with zero LastActive.
	critic, ok := roleByName["Critic"]
	if !ok {
		t.Fatal("Critic role missing from result")
	}
	if critic.Active {
		t.Error("Critic should not be Active (no posts)")
	}
	if !critic.LastActive.IsZero() {
		t.Error("Critic LastActive should be zero (never seen)")
	}
}

// TestGetHive_PublicNoAuth verifies GET /hive returns 200 without an auth cookie.
func TestGetHive_PublicNoAuth(t *testing.T) {
	h, _, _ := testHandlers(t)

	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("GET", "/hive", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("GET /hive without auth: status = %d, want 200; body: %s", w.Code, w.Body.String())
	}
}

// TestGetHive_RendersMetrics verifies the hive page contains stat card text
// when seeded with hive agent posts.
func TestGetHive_RendersMetrics(t *testing.T) {
	h, store, _ := testHandlers(t)

	mux := http.NewServeMux()
	h.Register(mux)

	// Create a public space to house the hive agent posts.
	space, err := store.CreateSpace(t.Context(), "hive-metrics-test", "Hive Metrics Test", "", "owner-hive-metrics", "project", "public")
	if err != nil {
		t.Fatalf("create space: %v", err)
	}
	t.Cleanup(func() { store.DeleteSpace(t.Context(), space.ID) })

	// Seed two posts by a hive agent (author_kind = "agent") with cost metadata.
	for i := range 2 {
		_, err := store.CreateNode(t.Context(), CreateNodeParams{
			SpaceID:    space.ID,
			Kind:       KindPost,
			Title:      fmt.Sprintf("[hive:builder] iter %d: shipped feature", i+230),
			Body:       fmt.Sprintf("Builder shipped the feature. Cost: $0.42. Duration: %dm30s.", i+3),
			Author:     "hive-builder",
			AuthorID:   "hive-agent-metrics-test-id",
			AuthorKind: "agent",
		})
		if err != nil {
			t.Fatalf("create post %d: %v", i, err)
		}
	}

	req := httptest.NewRequest("GET", "/hive", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body: %s", w.Code, w.Body.String())
	}
	body := w.Body.String()
	for _, label := range []string{"Features shipped", "Total autonomous spend", "Avg cost"} {
		if !strings.Contains(body, label) {
			t.Errorf("response body does not contain stat card label %q", label)
		}
	}
}

// TestGetHive_RendersCurrentlyBuilding verifies the "Currently building" section
// shows a task title when an open agent task exists, and "Idle" when none exists.
func TestGetHive_RendersCurrentlyBuilding(t *testing.T) {
	h, store, _ := testHandlers(t)

	mux := http.NewServeMux()
	h.Register(mux)

	// Without any agent tasks, "Idle" should appear.
	req := httptest.NewRequest("GET", "/hive", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	if !strings.Contains(w.Body.String(), "Idle") {
		t.Error("expected 'Idle' when no agent tasks exist")
	}

	// Create a space and seed an open agent task.
	space, err := store.CreateSpace(t.Context(), "hive-task-test", "Hive Task Test", "", "owner-hive-task", "project", "public")
	if err != nil {
		t.Fatalf("create space: %v", err)
	}
	t.Cleanup(func() { store.DeleteSpace(t.Context(), space.ID) })

	_, err = store.CreateNode(t.Context(), CreateNodeParams{
		SpaceID:    space.ID,
		Kind:       KindTask,
		Title:      "Build the knowledge layer",
		Body:       "Layer 6 implementation.",
		Author:     "hive-strategist",
		AuthorID:   "hive-agent-task-test-id",
		AuthorKind: "agent",
	})
	if err != nil {
		t.Fatalf("create task: %v", err)
	}

	req2 := httptest.NewRequest("GET", "/hive", nil)
	w2 := httptest.NewRecorder()
	mux.ServeHTTP(w2, req2)
	if w2.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body: %s", w2.Code, w2.Body.String())
	}
	if !strings.Contains(w2.Body.String(), "Build the knowledge layer") {
		t.Error("expected task title in 'Currently building' section")
	}
}

// TestGetHiveStats_Partial verifies GET /hive/stats returns 200 with stats bar HTML.
func TestGetHiveStats_Partial(t *testing.T) {
	h, _, _ := testHandlers(t)

	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("GET", "/hive/stats", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("GET /hive/stats: status = %d, want 200; body: %s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "total ops") {
		t.Error("expected 'total ops' in /hive/stats response")
	}
}
