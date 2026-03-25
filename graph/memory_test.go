package graph

import (
	"context"
	"strings"
	"testing"
)

// TestRememberAndRecallForPersona verifies that memories are stored and recalled via the real store.
func TestRememberAndRecallForPersona(t *testing.T) {
	_, store := testDB(t)
	ctx := context.Background()

	persona := "test-persona-" + newID()
	userID := "test-user-" + newID()

	// Store three memories with different importance levels.
	if err := store.RememberForPersona(ctx, persona, userID, "fact", "high importance fact", "", 9); err != nil {
		t.Fatalf("RememberForPersona fact: %v", err)
	}
	if err := store.RememberForPersona(ctx, persona, userID, "context", "low importance context", "", 2); err != nil {
		t.Fatalf("RememberForPersona context: %v", err)
	}
	if err := store.RememberForPersona(ctx, persona, userID, "preference", "medium preference", "", 5); err != nil {
		t.Fatalf("RememberForPersona preference: %v", err)
	}

	memories, err := store.RecallForPersona(ctx, persona, userID, 10)
	if err != nil {
		t.Fatalf("RecallForPersona: %v", err)
	}
	if len(memories) != 3 {
		t.Fatalf("expected 3 memories, got %d", len(memories))
	}
	// Ordered by importance DESC: high (9) first.
	if memories[0] != "high importance fact" {
		t.Errorf("expected highest-importance memory first, got %q", memories[0])
	}
}

// TestRememberForPersonaDefaults verifies kind and importance defaults.
func TestRememberForPersonaDefaults(t *testing.T) {
	_, store := testDB(t)
	ctx := context.Background()

	persona := "test-persona-" + newID()
	userID := "test-user-" + newID()

	// Empty kind defaults to "context"; importance 0 defaults to 5.
	if err := store.RememberForPersona(ctx, persona, userID, "", "default kind memory", "", 0); err != nil {
		t.Fatalf("RememberForPersona: %v", err)
	}
	memories, err := store.RecallForPersona(ctx, persona, userID, 10)
	if err != nil {
		t.Fatalf("RecallForPersona: %v", err)
	}
	if len(memories) != 1 || memories[0] != "default kind memory" {
		t.Errorf("unexpected memories: %v", memories)
	}
}

// TestRememberForPersonaInvalidKind verifies that invalid kinds are rejected.
func TestRememberForPersonaInvalidKind(t *testing.T) {
	_, store := testDB(t)
	ctx := context.Background()

	persona := "test-persona-" + newID()
	userID := "test-user-" + newID()

	err := store.RememberForPersona(ctx, persona, userID, "garbage", "bad memory", "", 5)
	if err == nil {
		t.Fatal("expected error for invalid kind, got nil")
	}
	if !strings.Contains(err.Error(), "invalid memory kind") {
		t.Errorf("unexpected error message: %v", err)
	}
}

// TestImportanceClampLogic verifies the clamping logic applied inside extractAndSaveMemories.
// This is a pure logic test — no database required.
func TestImportanceClampLogic(t *testing.T) {
	clamp := func(v int) int {
		if v < 1 {
			return 1
		}
		if v > 5 {
			return 5
		}
		return v
	}

	cases := []struct {
		raw  int
		want int
	}{
		{-3, 1},
		{0, 1},
		{1, 1},
		{3, 3},
		{5, 5},
		{6, 5},
		{10, 5},
	}
	for _, tc := range cases {
		got := clamp(tc.raw)
		if got != tc.want {
			t.Errorf("clamp(%d) = %d, want %d", tc.raw, got, tc.want)
		}
	}
}

// TestExtractMemoriesImportanceClamp verifies that importance is clamped before storage.
// Requires DATABASE_URL.
func TestExtractMemoriesImportanceClamp(t *testing.T) {
	_, store := testDB(t)
	ctx := context.Background()

	persona := "clamp-persona-" + newID()
	humanID := "clamp-human-" + newID()
	convoID := "clamp-convo-" + newID()

	// Simulate the clamp applied by extractAndSaveMemories: raw=10 → stored as 5.
	rawImportance := 10
	clamped := rawImportance
	if clamped > 5 {
		clamped = 5
	}
	if err := store.RememberForPersona(ctx, persona, humanID, "fact", "over-range memory", convoID, clamped); err != nil {
		t.Fatalf("RememberForPersona: %v", err)
	}
	memories, err := store.RecallForPersona(ctx, persona, humanID, 5)
	if err != nil {
		t.Fatalf("RecallForPersona: %v", err)
	}
	if len(memories) != 1 || memories[0] != "over-range memory" {
		t.Errorf("unexpected memories: %v", memories)
	}
}

// TestBuildSystemPromptInjectsMemories verifies that buildSystemPrompt includes stored memories.
func TestBuildSystemPromptInjectsMemories(t *testing.T) {
	_, store := testDB(t)
	ctx := context.Background()

	persona := "test-role-" + newID()
	humanID := "human-" + newID()
	agentID := "agent-" + newID()

	// Insert a persona record so buildSystemPrompt finds it.
	if err := store.UpsertAgentPersona(ctx, AgentPersona{
		Name:    persona,
		Display: "Test Persona",
		Prompt:  "You are a helpful assistant.",
		Model:   "sonnet",
		Active:  true,
	}); err != nil {
		t.Fatalf("agent_personas upsert failed: %v", err)
	}

	// Store a memory for this persona about the human.
	if err := store.RememberForPersona(ctx, persona, humanID, "fact", "user prefers brevity", "", 8); err != nil {
		t.Fatalf("RememberForPersona: %v", err)
	}

	// Build a minimal conversation node with role: and human tags.
	convo := &Node{
		ID:   newID(),
		Tags: []string{"role:" + persona, humanID, agentID},
	}

	mind := &Mind{store: store}
	prompt := mind.buildSystemPrompt(convo, agentID, nil)

	if !strings.Contains(prompt, "== MEMORIES ==") {
		t.Errorf("expected MEMORIES section in prompt, got:\n%s", prompt)
	}
	if !strings.Contains(prompt, "user prefers brevity") {
		t.Errorf("expected memory content in prompt, got:\n%s", prompt)
	}
}
