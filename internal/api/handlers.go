package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// Handler holds dependencies for HTTP handlers.
type Handler struct {
	NoAuth bool
}

// New creates a new Handler.
func New(noAuth bool) *Handler {
	return &Handler{NoAuth: noAuth}
}

// Routes returns the HTTP mux with all routes registered.
func (h *Handler) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/help", h.help)
	mux.HandleFunc("/.well-known/agent.md", h.help)
	mux.HandleFunc("/stats", h.stats)
	mux.HandleFunc("/count", h.count)
	mux.HandleFunc("/wrap", h.wrap)
	mux.HandleFunc("/number", h.numberLines)
	mux.HandleFunc("/dedup", h.dedup)
	mux.HandleFunc("/sort", h.sortLines)
	mux.HandleFunc("/trim", h.trimLines)
	mux.HandleFunc("/pad", h.pad)
	mux.HandleFunc("/reverse", h.reverse)
	mux.HandleFunc("/reverse-lines", h.reverseLines)
	mux.HandleFunc("/head", h.head)
	mux.HandleFunc("/tail", h.tail)
	mux.HandleFunc("/extract", h.extractLines)
	mux.HandleFunc("/find", h.find)
	mux.HandleFunc("/replace", h.replace)
	mux.HandleFunc("/squeeze", h.squeeze)
	mux.HandleFunc("/grep", h.grep)
	mux.HandleFunc("/unique", h.unique)
	mux.HandleFunc("/join", h.joinLines)
	mux.HandleFunc("/mcp", h.mcp)
	mux.HandleFunc("/", h.notFound)
	return mux
}

func wantsJSON(r *http.Request) bool {
	if r.URL.Query().Get("format") == "json" {
		return true
	}
	return strings.Contains(r.Header.Get("Accept"), "application/json")
}

func writeText(w http.ResponseWriter, text string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(w, text)
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, code int, msg, hint string) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "error: %s | hint: %s\n", msg, hint)
}

func getBody(r *http.Request) (string, error) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func getQueryParam(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}

func getIntParam(r *http.Request, key string, def int) int {
	v := r.URL.Query().Get(key)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}
