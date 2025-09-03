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
