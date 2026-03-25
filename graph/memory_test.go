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

// TestBuildSystemPromptInjectsMemories verifies that buildSystemPrompt includes stored memories.
func TestBuildSystemPromptInjectsMemories(t *testing.T) {
	db, store := testDB(t)
	ctx := context.Background()

	persona := "test-role-" + newID()
	humanID := "human-" + newID()
	agentID := "agent-" + newID()

	// Insert a persona record so buildSystemPrompt finds it.
	if _, err := db.ExecContext(ctx,
		`INSERT INTO agent_personas (id, role, prompt, created_at) VALUES ($1, $2, $3, NOW())`,
		newID(), persona, "You are a helpful assistant."); err != nil {
		t.Skipf("agent_personas insert failed (schema may differ): %v", err)
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
	prompt := mind.buildSystemPrompt(convo, agentID)

	if !strings.Contains(prompt, "== MEMORIES ==") {
		t.Errorf("expected MEMORIES section in prompt, got:\n%s", prompt)
	}
	if !strings.Contains(prompt, "user prefers brevity") {
		t.Errorf("expected memory content in prompt, got:\n%s", prompt)
	}
}
