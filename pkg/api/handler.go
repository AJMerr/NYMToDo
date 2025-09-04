package api

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/AJMerr/MSE/pkg/store"
	"github.com/AJMerr/NYMToDo/pkg/todo"
	"github.com/AJMerr/gonk/pkg/jsonutil"
	"github.com/AJMerr/gonk/pkg/router"
)

type Handler struct {
	S *store.Store
}

func New(s *store.Store) *Handler {
	return &Handler{S: s}
}

func (h *Handler) RegisterRouter(r *router.Router) {
	r.POST("/todos", h.createToDoHandler)
	r.GET("/todos", h.getAllToDoHandler)
	r.GET("/todos/{id}", h.getToDoHandler)
}

// Helper functions
const indexKey = "todos:index"

func makeKey(id string) string { return "todo:" + id }

// 16 random bytes to hex
func newID() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return strconv.FormatInt(time.Now().UnixNano(), 10)
	}
	return hex.EncodeToString(b[:])
}

func readIndex(s *store.Store) ([]string, error) {
	b, ok := s.Get(indexKey)
	if !ok || len(b) == 0 {
		return []string{}, nil
	}
	var ids []string
	if err := json.Unmarshal(b, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

func writeIndex(s *store.Store, ids []string) error {
	b, err := json.Marshal(ids)
	if err != nil {
		return err
	}
	return s.Set(indexKey, b)
}

// Handlers

// POST
func (h *Handler) createToDoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	title := strings.TrimSpace(req.Title)
	if title == "" {
		http.Error(w, "Title cannot be empty", http.StatusBadRequest)
		return
	}
	description := strings.TrimSpace(req.Description)

	id := newID()
	t := todo.ToDo{
		ID:          id,
		Title:       title,
		Description: description,
		IsCompleted: false,
	}

	tb, err := json.Marshal(t)
	if err != nil {
		http.Error(w, "store_error", http.StatusInternalServerError)
		return
	}
	if err := h.S.Set(makeKey(id), tb); err != nil {
		http.Error(w, "store_error", http.StatusInternalServerError)
		return
	}

	ids, err := readIndex(h.S)
	if err != nil {
		http.Error(w, "read_index_error", http.StatusInternalServerError)
		return
	}
	ids = append(ids, id)
	if err := writeIndex(h.S, ids); err != nil {
		http.Error(w, "write_index_error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(t)
}

// GET all
func (h *Handler) getAllToDoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	ids, err := readIndex(h.S)
	if err != nil {
		http.Error(w, "read_index_error", http.StatusInternalServerError)
		return
	}

	out := make([]todo.ToDo, 0, len(ids))
	for _, id := range ids {
		if b, ok := h.S.Get(makeKey(id)); ok {
			var t todo.ToDo
			if err := json.Unmarshal(b, &t); err == nil {
				out = append(out, t)
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(out)
}

// GET by ID
func (h *Handler) getToDoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonutil.WriteError(w, http.StatusMethodNotAllowed, "{error: method_not_allowed}")
		return
	}

	id := r.PathValue("id")
	if id == "" {
		jsonutil.WriteError(w, http.StatusNotFound, "{error: todo_not_found}")
		return
	}

	key := makeKey(id)
	bytes, ok := h.S.Get(key)
	if !ok {
		jsonutil.WriteError(w, http.StatusNotFound, "{error: todo_not_found}")
		return
	}

	var t todo.ToDo
	err := json.Unmarshal(bytes, &t)
	if err != nil {
		jsonutil.WriteError(w, http.StatusInternalServerError, "{error: decode_error}")
		return
	}
	jsonutil.WriteJSON(w, http.StatusOK, t)
}
