package api

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/AJMerr/MSE/pkg/store"
	"github.com/AJMerr/gonk/pkg/router"
)

type Handler struct {
	S *store.Store
}

func New(s *store.Store) *Handler

func (h *Handler) RegisterRouter(r *router.Router) {
	r.GET("/todos", h.createToDoHandler)
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

}
