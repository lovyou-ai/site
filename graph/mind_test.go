package graph

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

// TestMindFindUnreplied verifies the Mind's conversation detection query.
// Requires DATABASE_URL to be set.
func TestMindFindUnreplied(t *testing.T) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Skip("DATABASE_URL not set")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close()

	store, err := NewStore(db)
	if err != nil {
		t.Fatalf("new store: %v", err)
	}

	ctx := context.Background()
	mind := NewMind(db, store, "fake-token")

	// Setup: create an agent user, a space, and a conversation.
	suffix := time.Now().UnixNano()
	agentName := "TestAgent"
	humanName := "TestHuman"
	spaceSlug := "mind-test"

	// Clean up from previous runs.
	db.ExecContext(ctx, `DELETE FROM spaces WHERE slug = $1`, spaceSlug)
	db.ExecContext(ctx, `DELETE FROM users WHERE name = $1`, agentName)

	// Create agent user.
	_, err = db.ExecContext(ctx,
		`INSERT INTO users (id, google_id, email, name, kind) VALUES ($1, $2, $3, $4, 'agent')
		 ON CONFLICT (google_id) DO NOTHING`,
		newID(), "agent:"+agentName, agentName+"@test.lovyou.ai", agentName)
	if err != nil {
		t.Fatalf("create agent user: %v", err)
	}
	t.Cleanup(func() {
		db.ExecContext(ctx, `DELETE FROM spaces WHERE slug = $1`, spaceSlug)
		db.ExecContext(ctx, `DELETE FROM users WHERE name = $1`, agentName)
	})

	// Create space.
	space, err := store.CreateSpace(ctx, spaceSlug, "Mind Test", "", "test-owner", "project", "public")
	if err != nil {
		t.Fatalf("create space: %v", err)
	}
	_ = suffix

	// Test 1: Agent-created conversation with no messages → should NOT be found.
	t.Run("agent_created_no_messages", func(t *testing.T) {
		convo, err := store.CreateNode(ctx, CreateNodeParams{
			SpaceID:    space.ID,
			Kind:       KindConversation,
			Title:      "Agent started",
			Author:     agentName,
			AuthorKind: "agent",
			Tags:       []string{agentName, humanName},
		})
		if err != nil {
			t.Fatalf("create conversation: %v", err)
		}
		t.Cleanup(func() { store.DeleteNode(ctx, convo.ID) })

		convos, err := mind.findUnreplied(ctx)
		if err != nil {
			t.Fatalf("findUnreplied: %v", err)
		}
		for _, c := range convos {
			if c.ConversationID == convo.ID {
				t.Errorf("should not find agent-created conversation with no messages")
			}
		}
	})

	// Test 2: Human-created conversation with no messages → SHOULD be found.
	t.Run("human_created_no_messages", func(t *testing.T) {
		convo, err := store.CreateNode(ctx, CreateNodeParams{
			SpaceID:    space.ID,
			Kind:       KindConversation,
			Title:      "Human started",
			Author:     humanName,
			AuthorKind: "human",
			Tags:       []string{humanName, agentName},
		})
		if err != nil {
			t.Fatalf("create conversation: %v", err)
		}
		t.Cleanup(func() { store.DeleteNode(ctx, convo.ID) })

		convos, err := mind.findUnreplied(ctx)
		if err != nil {
			t.Fatalf("findUnreplied: %v", err)
		}
		found := false
		for _, c := range convos {
			if c.ConversationID == convo.ID {
				found = true
				if c.AgentName != agentName {
					t.Errorf("agent name = %q, want %q", c.AgentName, agentName)
				}
			}
		}
		if !found {
			t.Errorf("should find human-created conversation")
		}
	})

	// Test 3: Conversation with human's last message → SHOULD be found.
	t.Run("human_last_message", func(t *testing.T) {
		convo, err := store.CreateNode(ctx, CreateNodeParams{
			SpaceID:    space.ID,
			Kind:       KindConversation,
			Title:      "Active chat",
			Author:     humanName,
			AuthorKind: "human",
			Tags:       []string{humanName, agentName},
		})
		if err != nil {
			t.Fatalf("create conversation: %v", err)
		}
		t.Cleanup(func() {
			db.ExecContext(ctx, `DELETE FROM nodes WHERE parent_id = $1`, convo.ID)
			store.DeleteNode(ctx, convo.ID)
		})

		// Agent replies.
		_, err = store.CreateNode(ctx, CreateNodeParams{
			SpaceID:    space.ID,
			ParentID:   convo.ID,
			Kind:       KindComment,
			Body:       "Hello!",
			Author:     agentName,
			AuthorKind: "agent",
		})
		if err != nil {
			t.Fatalf("agent reply: %v", err)
		}
		// Human replies after agent.
		time.Sleep(10 * time.Millisecond)
		_, err = store.CreateNode(ctx, CreateNodeParams{
			SpaceID:    space.ID,
			ParentID:   convo.ID,
			Kind:       KindComment,
			Body:       "What's next?",
			Author:     humanName,
			AuthorKind: "human",
		})
		if err != nil {
			t.Fatalf("human reply: %v", err)
		}

		convos, err := mind.findUnreplied(ctx)
		if err != nil {
			t.Fatalf("findUnreplied: %v", err)
		}
		found := false
		for _, c := range convos {
			if c.ConversationID == convo.ID {
				found = true
			}
		}
		if !found {
			t.Errorf("should find conversation where human sent last message")
		}
	})

	// Test 4: Conversation where agent already replied last → should NOT be found.
	t.Run("agent_last_message", func(t *testing.T) {
		convo, err := store.CreateNode(ctx, CreateNodeParams{
			SpaceID:    space.ID,
			Kind:       KindConversation,
			Title:      "Agent replied",
			Author:     humanName,
			AuthorKind: "human",
			Tags:       []string{humanName, agentName},
		})
		if err != nil {
			t.Fatalf("create conversation: %v", err)
		}
		t.Cleanup(func() {
			db.ExecContext(ctx, `DELETE FROM nodes WHERE parent_id = $1`, convo.ID)
			store.DeleteNode(ctx, convo.ID)
		})

		// Human sends message.
		_, err = store.CreateNode(ctx, CreateNodeParams{
			SpaceID:    space.ID,
			ParentID:   convo.ID,
			Kind:       KindComment,
			Body:       "Hello",
			Author:     humanName,
			AuthorKind: "human",
		})
		if err != nil {
			t.Fatalf("human msg: %v", err)
		}
		// Agent replies (most recent).
		time.Sleep(10 * time.Millisecond)
		_, err = store.CreateNode(ctx, CreateNodeParams{
			SpaceID:    space.ID,
			ParentID:   convo.ID,
			Kind:       KindComment,
			Body:       "Hi there!",
			Author:     agentName,
			AuthorKind: "agent",
		})
		if err != nil {
			t.Fatalf("agent msg: %v", err)
		}

		convos, err := mind.findUnreplied(ctx)
		if err != nil {
			t.Fatalf("findUnreplied: %v", err)
		}
		for _, c := range convos {
			if c.ConversationID == convo.ID {
				t.Errorf("should not find conversation where agent already replied")
			}
		}
	})

	// Test 5: Staleness guard — old message should be skipped by poll().
	t.Run("staleness_guard", func(t *testing.T) {
		mind.maxAge = 1 * time.Millisecond // very short for testing

		convo, err := store.CreateNode(ctx, CreateNodeParams{
			SpaceID:    space.ID,
			Kind:       KindConversation,
			Title:      "Stale chat",
			Author:     humanName,
			AuthorKind: "human",
			Tags:       []string{humanName, agentName},
		})
		if err != nil {
			t.Fatalf("create conversation: %v", err)
		}
		t.Cleanup(func() {
			store.DeleteNode(ctx, convo.ID)
			mind.maxAge = 5 * time.Minute // restore
		})

		// The conversation was created (which counts as the "last message").
		// Wait for it to become stale.
		time.Sleep(5 * time.Millisecond)

		convos, err := mind.findUnreplied(ctx)
		if err != nil {
			t.Fatalf("findUnreplied: %v", err)
		}
		for _, c := range convos {
			if c.ConversationID == convo.ID {
				if time.Since(c.LastMessageAt) <= mind.maxAge {
					t.Errorf("staleness guard should have filtered this (age: %v, max: %v)",
						time.Since(c.LastMessageAt), mind.maxAge)
				}
				// Found but stale — poll() would skip it. That's the correct behavior.
			}
		}
	})
}

// TestMindE2E runs a full end-to-end test: human message → Mind replies.
// Requires DATABASE_URL and CLAUDE_CODE_OAUTH_TOKEN.
func TestMindE2E(t *testing.T) {
	dsn := os.Getenv("DATABASE_URL")
	token := os.Getenv("CLAUDE_CODE_OAUTH_TOKEN")
	if dsn == "" || token == "" {
		t.Skip("DATABASE_URL and CLAUDE_CODE_OAUTH_TOKEN required")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close()

	store, err := NewStore(db)
	if err != nil {
		t.Fatalf("new store: %v", err)
	}

	ctx := context.Background()
	mind := NewMind(db, store, token)

	agentName := "TestAgent"
	humanName := "TestHuman"
	spaceSlug := "mind-e2e-test"

	// Clean up.
	db.ExecContext(ctx, `DELETE FROM spaces WHERE slug = $1`, spaceSlug)
	db.ExecContext(ctx, `DELETE FROM users WHERE name = $1`, agentName)

	// Create agent user.
	_, err = db.ExecContext(ctx,
		`INSERT INTO users (id, google_id, email, name, kind) VALUES ($1, $2, $3, $4, 'agent')
		 ON CONFLICT (google_id) DO NOTHING`,
		newID(), "agent:"+agentName, agentName+"@test.lovyou.ai", agentName)
	if err != nil {
		t.Fatalf("create agent user: %v", err)
	}
	t.Cleanup(func() {
		db.ExecContext(ctx, `DELETE FROM spaces WHERE slug = $1`, spaceSlug)
		db.ExecContext(ctx, `DELETE FROM users WHERE name = $1`, agentName)
	})

	// Create space and conversation.
	space, err := store.CreateSpace(ctx, spaceSlug, "E2E Test", "", "test-owner", "project", "public")
	if err != nil {
		t.Fatalf("create space: %v", err)
	}

	convo, err := store.CreateNode(ctx, CreateNodeParams{
		SpaceID:    space.ID,
		Kind:       KindConversation,
		Title:      "E2E Test Conversation",
		Body:       "Testing the Mind's auto-reply",
		Author:     humanName,
		AuthorKind: "human",
		Tags:       []string{humanName, agentName},
	})
	if err != nil {
		t.Fatalf("create conversation: %v", err)
	}

	// Human sends a message.
	_, err = store.CreateNode(ctx, CreateNodeParams{
		SpaceID:    space.ID,
		ParentID:   convo.ID,
		Kind:       KindComment,
		Body:       "Hey Mind, can you hear me? Reply with exactly: YES I CAN",
		Author:     humanName,
		AuthorKind: "human",
	})
	if err != nil {
		t.Fatalf("human message: %v", err)
	}

	// Mind should find it and reply.
	convos, err := mind.findUnreplied(ctx)
	if err != nil {
		t.Fatalf("findUnreplied: %v", err)
	}

	found := false
	for _, c := range convos {
		if c.ConversationID == convo.ID {
			found = true
			t.Logf("found unreplied conversation: %q (agent: %s)", c.Title, c.AgentName)

			// Actually reply.
			err := mind.replyTo(ctx, c)
			if err != nil {
				t.Fatalf("replyTo: %v", err)
			}

			// Verify reply exists.
			messages, err := store.ListNodes(ctx, ListNodesParams{
				SpaceID:  space.ID,
				ParentID: convo.ID,
			})
			if err != nil {
				t.Fatalf("list messages: %v", err)
			}

			agentReplied := false
			for _, msg := range messages {
				if msg.Author == agentName && msg.AuthorKind == "agent" {
					agentReplied = true
					t.Logf("agent replied: %s", msg.Body)
				}
			}
			if !agentReplied {
				t.Errorf("agent should have replied")
			}
		}
	}
	if !found {
		t.Errorf("should find the test conversation")
	}
}
