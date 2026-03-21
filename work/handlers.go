package work

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/lovyou-ai/site/auth"
)

// ViewUser holds user info for templates.
type ViewUser struct {
	Name    string
	Picture string
}

// Handlers serves the Work product HTTP endpoints.
type Handlers struct {
	store *Store
	wrap  func(http.HandlerFunc) http.Handler
}

// NewHandlers creates Work product handlers with auth middleware.
// The wrap function is applied to every route for authentication.
func NewHandlers(store *Store, wrap func(http.HandlerFunc) http.Handler) *Handlers {
	if wrap == nil {
		wrap = func(hf http.HandlerFunc) http.Handler { return hf }
	}
	return &Handlers{store: store, wrap: wrap}
}

// Register adds all /work routes to the mux.
func (h *Handlers) Register(mux *http.ServeMux) {
	// Pages.
	mux.Handle("GET /work", h.wrap(h.handleIndex))
	mux.Handle("GET /work/project/{id}", h.wrap(h.handleBoard))
	mux.Handle("GET /work/project/{id}/list", h.wrap(h.handleList))
	mux.Handle("GET /work/task/{id}", h.wrap(h.handleTaskDetail))

	// Mutations.
	mux.Handle("POST /work/project", h.wrap(h.handleCreateProject))
	mux.Handle("POST /work/task", h.wrap(h.handleCreateTask))
	mux.Handle("POST /work/task/{id}/state", h.wrap(h.handleTransitionTask))
	mux.Handle("POST /work/task/{id}/update", h.wrap(h.handleUpdateTask))
	mux.Handle("POST /work/task/{id}/comment", h.wrap(h.handleAddComment))
	mux.Handle("DELETE /work/task/{id}", h.wrap(h.handleDeleteTask))

	// API.
	mux.Handle("GET /work/api/tasks", h.wrap(h.handleAPITasks))
}

// ────────────────────────────────────────────────────────────────────
// Helpers
// ────────────────────────────────────────────────────────────────────

func (h *Handlers) viewUser(r *http.Request) ViewUser {
	u := auth.UserFromContext(r.Context())
	if u == nil {
		return ViewUser{Name: "Anonymous"}
	}
	return ViewUser{Name: u.Name, Picture: u.Picture}
}

func (h *Handlers) userID(r *http.Request) string {
	u := auth.UserFromContext(r.Context())
	if u == nil {
		return "anonymous"
	}
	return u.ID
}

func (h *Handlers) userName(r *http.Request) string {
	u := auth.UserFromContext(r.Context())
	if u == nil {
		return "anonymous"
	}
	return u.Name
}

// verifyProjectAccess checks that the user owns the project.
func (h *Handlers) verifyProjectAccess(r *http.Request, projectID string) (*Project, error) {
	project, err := h.store.GetProject(r.Context(), projectID)
	if err != nil {
		return nil, err
	}
	if project.OwnerID != "" && project.OwnerID != h.userID(r) {
		return nil, ErrNotFound
	}
	return project, nil
}

// ────────────────────────────────────────────────────────────────────
// Page handlers
// ────────────────────────────────────────────────────────────────────

func (h *Handlers) handleIndex(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	projects, err := h.store.ListProjects(ctx, h.userID(r))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(projects) == 0 {
		WorkOnboarding(h.viewUser(r)).Render(ctx, w)
		return
	}

	http.Redirect(w, r, "/work/project/"+projects[0].ID, http.StatusSeeOther)
}

func (h *Handlers) handleBoard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	projectID := r.PathValue("id")

	project, err := h.verifyProjectAccess(r, projectID)
	if errors.Is(err, ErrNotFound) {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	projects, err := h.store.ListProjects(ctx, h.userID(r))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tasks, err := h.store.ListTasks(ctx, ListTasksParams{
		ProjectID: projectID,
		ParentID:  "root",
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	columns := groupByState(tasks)
	WorkBoard(*project, projects, columns, h.viewUser(r)).Render(ctx, w)
}

func (h *Handlers) handleList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	projectID := r.PathValue("id")

	project, err := h.verifyProjectAccess(r, projectID)
	if errors.Is(err, ErrNotFound) {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	projects, err := h.store.ListProjects(ctx, h.userID(r))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	params := ListTasksParams{
		ProjectID: projectID,
		ParentID:  "root",
	}
	if state := r.URL.Query().Get("state"); state != "" {
		params.State = state
	}
	if assignee := r.URL.Query().Get("assignee"); assignee != "" {
		params.Assignee = assignee
	}

	tasks, err := h.store.ListTasks(ctx, params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	WorkList(*project, projects, tasks, h.viewUser(r)).Render(ctx, w)
}

func (h *Handlers) handleTaskDetail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	taskID := r.PathValue("id")

	task, err := h.store.GetTask(ctx, taskID)
	if errors.Is(err, ErrNotFound) {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	project, err := h.verifyProjectAccess(r, task.ProjectID)
	if errors.Is(err, ErrNotFound) {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	comments, err := h.store.ListComments(ctx, taskID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	subtasks, err := h.store.ListTasks(ctx, ListTasksParams{
		ProjectID: task.ProjectID,
		ParentID:  taskID,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	blockers, err := h.store.ListBlockers(ctx, taskID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	WorkTaskDetail(*project, *task, comments, subtasks, blockers, h.viewUser(r)).Render(ctx, w)
}

// ────────────────────────────────────────────────────────────────────
// Mutation handlers
// ────────────────────────────────────────────────────────────────────

func (h *Handlers) handleCreateProject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	name := strings.TrimSpace(r.FormValue("name"))
	description := strings.TrimSpace(r.FormValue("description"))

	if name == "" {
		http.Error(w, "project name is required", http.StatusBadRequest)
		return
	}

	project, err := h.store.CreateProject(ctx, name, description, h.userName(r), h.userID(r))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/work/project/"+project.ID, http.StatusSeeOther)
}

func (h *Handlers) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	title := strings.TrimSpace(r.FormValue("title"))
	if title == "" {
		http.Error(w, "task title is required", http.StatusBadRequest)
		return
	}

	projectID := r.FormValue("project_id")
	if _, err := h.verifyProjectAccess(r, projectID); err != nil {
		http.Error(w, "project not found", http.StatusNotFound)
		return
	}

	priority := r.FormValue("priority")
	if priority == "" {
		priority = PriorityMedium
	}

	params := CreateTaskParams{
		Title:       title,
		Description: strings.TrimSpace(r.FormValue("description")),
		Priority:    priority,
		ProjectID:   projectID,
		ParentID:    r.FormValue("parent_id"),
		Assignee:    strings.TrimSpace(r.FormValue("assignee")),
		CreatedBy:   h.userName(r),
	}

	task, err := h.store.CreateTask(ctx, params)
	if err != nil {
		log.Printf("work: create task: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if isHTMX(r) {
		TaskCard(*task).Render(ctx, w)
		return
	}

	http.Redirect(w, r, "/work/project/"+task.ProjectID, http.StatusSeeOther)
}

func (h *Handlers) handleTransitionTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	taskID := r.PathValue("id")
	newState := r.FormValue("state")

	// Verify access.
	task, err := h.store.GetTask(ctx, taskID)
	if errors.Is(err, ErrNotFound) {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err := h.verifyProjectAccess(r, task.ProjectID); err != nil {
		http.NotFound(w, r)
		return
	}

	if err := h.store.TransitionTask(ctx, taskID, newState); err != nil {
		if errors.Is(err, ErrNotFound) {
			http.NotFound(w, r)
			return
		}
		if errors.Is(err, ErrInvalidState) || errors.Is(err, ErrInvalidTransition) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	task, err = h.store.GetTask(ctx, taskID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if isHTMX(r) {
		TaskCard(*task).Render(ctx, w)
		return
	}
	http.Redirect(w, r, "/work/project/"+task.ProjectID, http.StatusSeeOther)
}

func (h *Handlers) handleUpdateTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	taskID := r.PathValue("id")

	// Verify access.
	task, err := h.store.GetTask(ctx, taskID)
	if errors.Is(err, ErrNotFound) {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err := h.verifyProjectAccess(r, task.ProjectID); err != nil {
		http.NotFound(w, r)
		return
	}

	params := UpdateTaskParams{}
	if v := r.FormValue("title"); v != "" {
		params.Title = &v
	}
	if v := r.FormValue("description"); r.Form != nil && r.Form.Has("description") {
		params.Description = &v
	}
	if v := r.FormValue("priority"); v != "" {
		params.Priority = &v
	}
	if r.Form != nil && r.Form.Has("assignee") {
		v := r.FormValue("assignee")
		params.Assignee = &v
	}

	if err := h.store.UpdateTask(ctx, taskID, params); err != nil {
		if errors.Is(err, ErrNotFound) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	task, err = h.store.GetTask(ctx, taskID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if isHTMX(r) {
		TaskCard(*task).Render(ctx, w)
		return
	}
	http.Redirect(w, r, "/work/task/"+taskID, http.StatusSeeOther)
}

func (h *Handlers) handleAddComment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	taskID := r.PathValue("id")

	body := strings.TrimSpace(r.FormValue("body"))
	if body == "" {
		http.Error(w, "comment body is required", http.StatusBadRequest)
		return
	}

	// Verify access.
	task, err := h.store.GetTask(ctx, taskID)
	if errors.Is(err, ErrNotFound) {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err := h.verifyProjectAccess(r, task.ProjectID); err != nil {
		http.NotFound(w, r)
		return
	}

	comment, err := h.store.AddComment(ctx, taskID, h.userName(r), body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if isHTMX(r) {
		CommentItem(*comment).Render(ctx, w)
		return
	}
	http.Redirect(w, r, "/work/task/"+taskID, http.StatusSeeOther)
}

func (h *Handlers) handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	taskID := r.PathValue("id")

	task, err := h.store.GetTask(ctx, taskID)
	if errors.Is(err, ErrNotFound) {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err := h.verifyProjectAccess(r, task.ProjectID); err != nil {
		http.NotFound(w, r)
		return
	}

	if err := h.store.DeleteTask(ctx, task.ID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if isHTMX(r) {
		w.Header().Set("HX-Redirect", "/work/project/"+task.ProjectID)
		w.WriteHeader(http.StatusOK)
		return
	}
	http.Redirect(w, r, "/work/project/"+task.ProjectID, http.StatusSeeOther)
}

// ────────────────────────────────────────────────────────────────────
// API handlers
// ────────────────────────────────────────────────────────────────────

func (h *Handlers) handleAPITasks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	projectID := r.URL.Query().Get("project")
	if projectID == "" {
		http.Error(w, "project query param required", http.StatusBadRequest)
		return
	}

	if _, err := h.verifyProjectAccess(r, projectID); err != nil {
		http.Error(w, "project not found", http.StatusNotFound)
		return
	}

	params := ListTasksParams{
		ProjectID: projectID,
		ParentID:  "root",
	}
	if state := r.URL.Query().Get("state"); state != "" {
		params.State = state
	}

	tasks, err := h.store.ListTasks(ctx, params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// ────────────────────────────────────────────────────────────────────
// Board helpers
// ────────────────────────────────────────────────────────────────────

// BoardColumn holds tasks grouped by state for the kanban board.
type BoardColumn struct {
	State string
	Label string
	Tasks []Task
}

func groupByState(tasks []Task) []BoardColumn {
	columns := []BoardColumn{
		{State: StateBacklog, Label: "Backlog"},
		{State: StateTodo, Label: "To Do"},
		{State: StateDoing, Label: "Doing"},
		{State: StateReview, Label: "Review"},
		{State: StateDone, Label: "Done"},
	}
	byState := map[string]*BoardColumn{}
	for i := range columns {
		byState[columns[i].State] = &columns[i]
	}
	for _, t := range tasks {
		if col, ok := byState[t.State]; ok {
			col.Tasks = append(col.Tasks, t)
		}
	}
	return columns
}

func isHTMX(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}
