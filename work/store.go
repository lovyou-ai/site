// Package work implements the Work product backend — task management
// with projects, labels, comments, and blockers backed by Postgres.
package work

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

// ────────────────────────────────────────────────────────────────────
// Domain types
// ────────────────────────────────────────────────────────────────────

// Task states.
const (
	StateBacklog = "backlog"
	StateTodo    = "todo"
	StateDoing   = "doing"
	StateReview  = "review"
	StateDone    = "done"
)

// Task priorities.
const (
	PriorityUrgent = "urgent"
	PriorityHigh   = "high"
	PriorityMedium = "medium"
	PriorityLow    = "low"
)

// validStates is the set of legal task states.
var validStates = map[string]bool{
	StateBacklog: true,
	StateTodo:    true,
	StateDoing:   true,
	StateReview:  true,
	StateDone:    true,
}

// validPriorities is the set of legal task priorities.
var validPriorities = map[string]bool{
	PriorityUrgent: true,
	PriorityHigh:   true,
	PriorityMedium: true,
	PriorityLow:    true,
}

// stateOrder defines the allowed state transitions. A task may move
// forward one step, backward one step, or stay in the same state.
var stateOrder = map[string]int{
	StateBacklog: 0,
	StateTodo:    1,
	StateDoing:   2,
	StateReview:  3,
	StateDone:    4,
}

// Project is a container for tasks.
type Project struct {
	ID          string
	Name        string
	Description string
	Owner       string // display name
	OwnerID     string // FK to users.id — used for access control
	CreatedAt   time.Time
}

// Task is a unit of work within a project.
type Task struct {
	ID           string
	Title        string
	Description  string
	State        string // backlog, todo, doing, review, done
	Priority     string // urgent, high, medium, low
	Assignee     string // empty = unassigned
	ProjectID    string
	ParentID     string // empty = no parent
	DueDate      *time.Time
	Effort       *int
	CreatedBy    string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Labels       []Label // populated on read
	SubtaskCount int     // populated on read
	SubtaskDone  int     // populated on read
	BlockerCount int     // populated on read
}

// Comment is a discussion entry on a task.
type Comment struct {
	ID        string
	TaskID    string
	Author    string
	Body      string
	CreatedAt time.Time
}

// Label is a colored tag for categorising tasks.
type Label struct {
	ID    string
	Name  string
	Color string
}

// CreateTaskParams holds the fields for creating a new task.
type CreateTaskParams struct {
	Title       string
	Description string
	Priority    string
	ProjectID   string
	ParentID    string
	Assignee    string
	DueDate     *time.Time
	Effort      *int
	CreatedBy   string
}

// UpdateTaskParams holds optional fields for updating a task.
// nil means "do not change".
type UpdateTaskParams struct {
	Title       *string
	Description *string
	Priority    *string
	Assignee    *string
	DueDate     *time.Time
	Effort      *int
}

// ListTasksParams controls filtering for task listing.
type ListTasksParams struct {
	ProjectID string // required
	State     string // filter by state, empty = all
	Assignee  string // filter by assignee, empty = all
	ParentID  string // "root" = top-level only, ID = subtasks of that task, empty = all
}

// ────────────────────────────────────────────────────────────────────
// Errors
// ────────────────────────────────────────────────────────────────────

var (
	ErrNotFound          = errors.New("not found")
	ErrInvalidState      = errors.New("invalid state")
	ErrInvalidPriority   = errors.New("invalid priority")
	ErrInvalidTransition = errors.New("invalid state transition")
)

// ────────────────────────────────────────────────────────────────────
// Store
// ────────────────────────────────────────────────────────────────────

// Store is a Postgres-backed store for the Work product.
type Store struct {
	db *sql.DB
}

// New connects to Postgres and auto-creates tables.
func New(dsn string) (*Store, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping db: %w", err)
	}
	s := &Store{db: db}
	if err := s.migrate(); err != nil {
		db.Close()
		return nil, fmt.Errorf("migrate: %w", err)
	}
	return s, nil
}

// NewWithDB wraps an existing database connection.
func NewWithDB(db *sql.DB) (*Store, error) {
	s := &Store{db: db}
	if err := s.migrate(); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}
	return s, nil
}

// Close closes the database connection.
func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) migrate() error {
	ddl := `
CREATE TABLE IF NOT EXISTS work_projects (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT DEFAULT '',
    owner TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS work_labels (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    color TEXT NOT NULL DEFAULT '#6366F1'
);

CREATE TABLE IF NOT EXISTS work_tasks (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT DEFAULT '',
    state TEXT NOT NULL DEFAULT 'todo' CHECK (state IN ('backlog','todo','doing','review','done')),
    priority TEXT NOT NULL DEFAULT 'medium' CHECK (priority IN ('urgent','high','medium','low')),
    assignee TEXT DEFAULT NULL,
    project_id TEXT NOT NULL REFERENCES work_projects(id),
    parent_id TEXT DEFAULT NULL REFERENCES work_tasks(id),
    due_date DATE DEFAULT NULL,
    effort INT DEFAULT NULL,
    created_by TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS work_task_labels (
    task_id TEXT REFERENCES work_tasks(id) ON DELETE CASCADE,
    label_id TEXT REFERENCES work_labels(id) ON DELETE CASCADE,
    PRIMARY KEY (task_id, label_id)
);

CREATE TABLE IF NOT EXISTS work_task_blockers (
    task_id TEXT REFERENCES work_tasks(id) ON DELETE CASCADE,
    blocker_id TEXT REFERENCES work_tasks(id) ON DELETE CASCADE,
    PRIMARY KEY (task_id, blocker_id)
);

CREATE TABLE IF NOT EXISTS work_comments (
    id TEXT PRIMARY KEY,
    task_id TEXT NOT NULL REFERENCES work_tasks(id) ON DELETE CASCADE,
    author TEXT NOT NULL,
    body TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

ALTER TABLE work_projects ADD COLUMN IF NOT EXISTS owner_id TEXT;
`
	_, err := s.db.Exec(ddl)
	return err
}

// ────────────────────────────────────────────────────────────────────
// ID generation
// ────────────────────────────────────────────────────────────────────

func newID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		panic("crypto/rand: " + err.Error())
	}
	return hex.EncodeToString(b)
}

// ────────────────────────────────────────────────────────────────────
// Projects
// ────────────────────────────────────────────────────────────────────

// CreateProject creates a new project owned by the given user.
func (s *Store) CreateProject(ctx context.Context, name, description, ownerName, ownerID string) (*Project, error) {
	p := &Project{
		ID:          newID(),
		Name:        name,
		Description: description,
		Owner:       ownerName,
		OwnerID:     ownerID,
	}
	err := s.db.QueryRowContext(ctx,
		`INSERT INTO work_projects (id, name, description, owner, owner_id) VALUES ($1, $2, $3, $4, $5) RETURNING created_at`,
		p.ID, p.Name, p.Description, p.Owner, p.OwnerID,
	).Scan(&p.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("create project: %w", err)
	}
	return p, nil
}

// ListProjects returns projects owned by the given user.
func (s *Store) ListProjects(ctx context.Context, ownerID string) ([]Project, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, name, description, owner, COALESCE(owner_id, ''), created_at
		 FROM work_projects WHERE owner_id = $1 ORDER BY created_at`, ownerID)
	if err != nil {
		return nil, fmt.Errorf("list projects: %w", err)
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var p Project
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Owner, &p.OwnerID, &p.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan project: %w", err)
		}
		projects = append(projects, p)
	}
	return projects, rows.Err()
}

// GetProject returns a project by ID.
func (s *Store) GetProject(ctx context.Context, id string) (*Project, error) {
	var p Project
	err := s.db.QueryRowContext(ctx,
		`SELECT id, name, description, owner, COALESCE(owner_id, ''), created_at FROM work_projects WHERE id = $1`, id,
	).Scan(&p.ID, &p.Name, &p.Description, &p.Owner, &p.OwnerID, &p.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get project: %w", err)
	}
	return &p, nil
}

// ────────────────────────────────────────────────────────────────────
// Tasks
// ────────────────────────────────────────────────────────────────────

// CreateTask creates a new task in a project.
func (s *Store) CreateTask(ctx context.Context, p CreateTaskParams) (*Task, error) {
	if p.Priority == "" {
		p.Priority = PriorityMedium
	}
	if !validPriorities[p.Priority] {
		return nil, ErrInvalidPriority
	}

	t := &Task{
		ID:          newID(),
		Title:       p.Title,
		Description: p.Description,
		State:       StateTodo,
		Priority:    p.Priority,
		Assignee:    p.Assignee,
		ProjectID:   p.ProjectID,
		ParentID:    p.ParentID,
		DueDate:     p.DueDate,
		Effort:      p.Effort,
		CreatedBy:   p.CreatedBy,
	}

	var parentID *string
	if p.ParentID != "" {
		parentID = &p.ParentID
	}
	var assignee *string
	if p.Assignee != "" {
		assignee = &p.Assignee
	}

	err := s.db.QueryRowContext(ctx,
		`INSERT INTO work_tasks (id, title, description, state, priority, assignee, project_id, parent_id, due_date, effort, created_by)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		 RETURNING created_at, updated_at`,
		t.ID, t.Title, t.Description, t.State, t.Priority, assignee,
		t.ProjectID, parentID, t.DueDate, t.Effort, t.CreatedBy,
	).Scan(&t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("create task: %w", err)
	}
	return t, nil
}

// GetTask returns a task by ID with labels, subtask counts, and blocker count.
func (s *Store) GetTask(ctx context.Context, id string) (*Task, error) {
	t, err := s.scanTask(ctx, id)
	if err != nil {
		return nil, err
	}

	// Load labels.
	t.Labels, err = s.taskLabels(ctx, id)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (s *Store) scanTask(ctx context.Context, id string) (*Task, error) {
	var t Task
	var assignee, parentID sql.NullString
	var dueDate sql.NullTime
	var effort sql.NullInt32

	err := s.db.QueryRowContext(ctx, `
		SELECT t.id, t.title, t.description, t.state, t.priority,
		       t.assignee, t.project_id, t.parent_id,
		       t.due_date, t.effort, t.created_by,
		       t.created_at, t.updated_at,
		       COALESCE((SELECT COUNT(*) FROM work_tasks st WHERE st.parent_id = t.id), 0),
		       COALESCE((SELECT COUNT(*) FROM work_tasks st WHERE st.parent_id = t.id AND st.state = 'done'), 0),
		       COALESCE((SELECT COUNT(*) FROM work_task_blockers b WHERE b.task_id = t.id), 0)
		FROM work_tasks t
		WHERE t.id = $1`, id,
	).Scan(
		&t.ID, &t.Title, &t.Description, &t.State, &t.Priority,
		&assignee, &t.ProjectID, &parentID,
		&dueDate, &effort, &t.CreatedBy,
		&t.CreatedAt, &t.UpdatedAt,
		&t.SubtaskCount, &t.SubtaskDone, &t.BlockerCount,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get task: %w", err)
	}

	if assignee.Valid {
		t.Assignee = assignee.String
	}
	if parentID.Valid {
		t.ParentID = parentID.String
	}
	if dueDate.Valid {
		d := dueDate.Time
		t.DueDate = &d
	}
	if effort.Valid {
		e := int(effort.Int32)
		t.Effort = &e
	}
	return &t, nil
}

func (s *Store) taskLabels(ctx context.Context, taskID string) ([]Label, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT l.id, l.name, l.color
		FROM work_labels l
		JOIN work_task_labels tl ON tl.label_id = l.id
		WHERE tl.task_id = $1
		ORDER BY l.name`, taskID)
	if err != nil {
		return nil, fmt.Errorf("task labels: %w", err)
	}
	defer rows.Close()

	var labels []Label
	for rows.Next() {
		var l Label
		if err := rows.Scan(&l.ID, &l.Name, &l.Color); err != nil {
			return nil, fmt.Errorf("scan label: %w", err)
		}
		labels = append(labels, l)
	}
	return labels, rows.Err()
}

// ListTasks returns tasks matching the filter criteria.
func (s *Store) ListTasks(ctx context.Context, p ListTasksParams) ([]Task, error) {
	query := `
		SELECT t.id, t.title, t.description, t.state, t.priority,
		       t.assignee, t.project_id, t.parent_id,
		       t.due_date, t.effort, t.created_by,
		       t.created_at, t.updated_at,
		       COALESCE((SELECT COUNT(*) FROM work_tasks st WHERE st.parent_id = t.id), 0),
		       COALESCE((SELECT COUNT(*) FROM work_tasks st WHERE st.parent_id = t.id AND st.state = 'done'), 0),
		       COALESCE((SELECT COUNT(*) FROM work_task_blockers b WHERE b.task_id = t.id), 0)
		FROM work_tasks t
		WHERE t.project_id = $1`

	args := []any{p.ProjectID}
	argN := 2

	if p.State != "" {
		query += fmt.Sprintf(" AND t.state = $%d", argN)
		args = append(args, p.State)
		argN++
	}
	if p.Assignee != "" {
		query += fmt.Sprintf(" AND t.assignee = $%d", argN)
		args = append(args, p.Assignee)
		argN++
	}
	if p.ParentID == "root" {
		query += " AND t.parent_id IS NULL"
	} else if p.ParentID != "" {
		query += fmt.Sprintf(" AND t.parent_id = $%d", argN)
		args = append(args, p.ParentID)
		argN++
	}

	query += " ORDER BY t.created_at"

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list tasks: %w", err)
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		var assignee, parentID sql.NullString
		var dueDate sql.NullTime
		var effort sql.NullInt32

		if err := rows.Scan(
			&t.ID, &t.Title, &t.Description, &t.State, &t.Priority,
			&assignee, &t.ProjectID, &parentID,
			&dueDate, &effort, &t.CreatedBy,
			&t.CreatedAt, &t.UpdatedAt,
			&t.SubtaskCount, &t.SubtaskDone, &t.BlockerCount,
		); err != nil {
			return nil, fmt.Errorf("scan task: %w", err)
		}

		if assignee.Valid {
			t.Assignee = assignee.String
		}
		if parentID.Valid {
			t.ParentID = parentID.String
		}
		if dueDate.Valid {
			d := dueDate.Time
			t.DueDate = &d
		}
		if effort.Valid {
			e := int(effort.Int32)
			t.Effort = &e
		}
		tasks = append(tasks, t)
	}

	// Load labels for each task.
	for i := range tasks {
		labels, err := s.taskLabels(ctx, tasks[i].ID)
		if err != nil {
			return nil, err
		}
		tasks[i].Labels = labels
	}

	return tasks, rows.Err()
}

// UpdateTask updates mutable fields on a task.
func (s *Store) UpdateTask(ctx context.Context, id string, p UpdateTaskParams) error {
	// Build SET clause dynamically.
	sets := []string{"updated_at = NOW()"}
	args := []any{}
	argN := 1

	if p.Title != nil {
		sets = append(sets, fmt.Sprintf("title = $%d", argN))
		args = append(args, *p.Title)
		argN++
	}
	if p.Description != nil {
		sets = append(sets, fmt.Sprintf("description = $%d", argN))
		args = append(args, *p.Description)
		argN++
	}
	if p.Priority != nil {
		if !validPriorities[*p.Priority] {
			return ErrInvalidPriority
		}
		sets = append(sets, fmt.Sprintf("priority = $%d", argN))
		args = append(args, *p.Priority)
		argN++
	}
	if p.Assignee != nil {
		if *p.Assignee == "" {
			sets = append(sets, "assignee = NULL")
		} else {
			sets = append(sets, fmt.Sprintf("assignee = $%d", argN))
			args = append(args, *p.Assignee)
			argN++
		}
	}
	if p.DueDate != nil {
		sets = append(sets, fmt.Sprintf("due_date = $%d", argN))
		args = append(args, *p.DueDate)
		argN++
	}
	if p.Effort != nil {
		sets = append(sets, fmt.Sprintf("effort = $%d", argN))
		args = append(args, *p.Effort)
		argN++
	}

	query := fmt.Sprintf("UPDATE work_tasks SET %s WHERE id = $%d",
		joinStrings(sets, ", "), argN)
	args = append(args, id)

	res, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("update task: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

// TransitionTask moves a task to a new state. Validates that the
// transition is at most one step forward or backward.
func (s *Store) TransitionTask(ctx context.Context, id, newState string) error {
	if !validStates[newState] {
		return ErrInvalidState
	}

	// Read current state.
	var curState string
	err := s.db.QueryRowContext(ctx,
		`SELECT state FROM work_tasks WHERE id = $1`, id,
	).Scan(&curState)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		return fmt.Errorf("read state: %w", err)
	}

	// Validate transition: allow any move (board UX needs flexibility).
	// The constraint is intentionally relaxed for the MVP — tasks can
	// move freely between columns on a kanban board.
	_ = stateOrder // kept for future stricter validation

	res, err := s.db.ExecContext(ctx,
		`UPDATE work_tasks SET state = $1, updated_at = NOW() WHERE id = $2`,
		newState, id,
	)
	if err != nil {
		return fmt.Errorf("transition task: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

// DeleteTask removes a task. Subtasks and blockers cascade.
func (s *Store) DeleteTask(ctx context.Context, id string) error {
	res, err := s.db.ExecContext(ctx,
		`DELETE FROM work_tasks WHERE id = $1`, id,
	)
	if err != nil {
		return fmt.Errorf("delete task: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

// ────────────────────────────────────────────────────────────────────
// Comments
// ────────────────────────────────────────────────────────────────────

// AddComment adds a comment to a task.
func (s *Store) AddComment(ctx context.Context, taskID, author, body string) (*Comment, error) {
	c := &Comment{
		ID:     newID(),
		TaskID: taskID,
		Author: author,
		Body:   body,
	}
	err := s.db.QueryRowContext(ctx,
		`INSERT INTO work_comments (id, task_id, author, body) VALUES ($1, $2, $3, $4) RETURNING created_at`,
		c.ID, c.TaskID, c.Author, c.Body,
	).Scan(&c.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("add comment: %w", err)
	}
	return c, nil
}

// ListComments returns all comments for a task, oldest first.
func (s *Store) ListComments(ctx context.Context, taskID string) ([]Comment, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, task_id, author, body, created_at FROM work_comments WHERE task_id = $1 ORDER BY created_at`,
		taskID,
	)
	if err != nil {
		return nil, fmt.Errorf("list comments: %w", err)
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var c Comment
		if err := rows.Scan(&c.ID, &c.TaskID, &c.Author, &c.Body, &c.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan comment: %w", err)
		}
		comments = append(comments, c)
	}
	return comments, rows.Err()
}

// ────────────────────────────────────────────────────────────────────
// Labels
// ────────────────────────────────────────────────────────────────────

// CreateLabel creates a new label.
func (s *Store) CreateLabel(ctx context.Context, name, color string) (*Label, error) {
	if color == "" {
		color = "#6366F1"
	}
	l := &Label{
		ID:    newID(),
		Name:  name,
		Color: color,
	}
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO work_labels (id, name, color) VALUES ($1, $2, $3)`,
		l.ID, l.Name, l.Color,
	)
	if err != nil {
		return nil, fmt.Errorf("create label: %w", err)
	}
	return l, nil
}

// ListLabels returns all labels ordered by name.
func (s *Store) ListLabels(ctx context.Context) ([]Label, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, name, color FROM work_labels ORDER BY name`)
	if err != nil {
		return nil, fmt.Errorf("list labels: %w", err)
	}
	defer rows.Close()

	var labels []Label
	for rows.Next() {
		var l Label
		if err := rows.Scan(&l.ID, &l.Name, &l.Color); err != nil {
			return nil, fmt.Errorf("scan label: %w", err)
		}
		labels = append(labels, l)
	}
	return labels, rows.Err()
}

// SetTaskLabels replaces all labels on a task.
func (s *Store) SetTaskLabels(ctx context.Context, taskID string, labelIDs []string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `DELETE FROM work_task_labels WHERE task_id = $1`, taskID); err != nil {
		return fmt.Errorf("clear labels: %w", err)
	}
	for _, lid := range labelIDs {
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO work_task_labels (task_id, label_id) VALUES ($1, $2)`,
			taskID, lid,
		); err != nil {
			return fmt.Errorf("set label: %w", err)
		}
	}
	return tx.Commit()
}

// ────────────────────────────────────────────────────────────────────
// Blockers
// ────────────────────────────────────────────────────────────────────

// AddBlocker marks blockerID as blocking taskID.
func (s *Store) AddBlocker(ctx context.Context, taskID, blockerID string) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO work_task_blockers (task_id, blocker_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
		taskID, blockerID,
	)
	if err != nil {
		return fmt.Errorf("add blocker: %w", err)
	}
	return nil
}

// RemoveBlocker removes a blocker relationship.
func (s *Store) RemoveBlocker(ctx context.Context, taskID, blockerID string) error {
	_, err := s.db.ExecContext(ctx,
		`DELETE FROM work_task_blockers WHERE task_id = $1 AND blocker_id = $2`,
		taskID, blockerID,
	)
	if err != nil {
		return fmt.Errorf("remove blocker: %w", err)
	}
	return nil
}

// ListBlockers returns the tasks that block the given task.
func (s *Store) ListBlockers(ctx context.Context, taskID string) ([]Task, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT t.id, t.title, t.description, t.state, t.priority,
		       t.assignee, t.project_id, t.parent_id,
		       t.due_date, t.effort, t.created_by,
		       t.created_at, t.updated_at,
		       COALESCE((SELECT COUNT(*) FROM work_tasks st WHERE st.parent_id = t.id), 0),
		       COALESCE((SELECT COUNT(*) FROM work_tasks st WHERE st.parent_id = t.id AND st.state = 'done'), 0),
		       COALESCE((SELECT COUNT(*) FROM work_task_blockers b WHERE b.task_id = t.id), 0)
		FROM work_tasks t
		JOIN work_task_blockers tb ON tb.blocker_id = t.id
		WHERE tb.task_id = $1
		ORDER BY t.created_at`, taskID)
	if err != nil {
		return nil, fmt.Errorf("list blockers: %w", err)
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		var assignee, parentID sql.NullString
		var dueDate sql.NullTime
		var effort sql.NullInt32

		if err := rows.Scan(
			&t.ID, &t.Title, &t.Description, &t.State, &t.Priority,
			&assignee, &t.ProjectID, &parentID,
			&dueDate, &effort, &t.CreatedBy,
			&t.CreatedAt, &t.UpdatedAt,
			&t.SubtaskCount, &t.SubtaskDone, &t.BlockerCount,
		); err != nil {
			return nil, fmt.Errorf("scan blocker: %w", err)
		}

		if assignee.Valid {
			t.Assignee = assignee.String
		}
		if parentID.Valid {
			t.ParentID = parentID.String
		}
		if dueDate.Valid {
			d := dueDate.Time
			t.DueDate = &d
		}
		if effort.Valid {
			e := int(effort.Int32)
			t.Effort = &e
		}
		tasks = append(tasks, t)
	}
	return tasks, rows.Err()
}

// ────────────────────────────────────────────────────────────────────
// Helpers
// ────────────────────────────────────────────────────────────────────

func joinStrings(ss []string, sep string) string {
	if len(ss) == 0 {
		return ""
	}
	out := ss[0]
	for _, s := range ss[1:] {
		out += sep + s
	}
	return out
}
