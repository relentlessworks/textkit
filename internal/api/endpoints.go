package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/relentlessworks/textkit/internal/model"
)

func (h *Handler) help(w http.ResponseWriter, r *http.Request) {
	text := `textkit — agentic-first text analysis and manipulation service

AUTH:
  No auth required (stateless computation service).

ENDPOINTS:
  POST /stats          — text statistics (words, chars, lines, sentences, paragraphs, bytes, runes)
  POST /count          — count words, chars, bytes, or lines
  POST /wrap           — wrap text to a given width
  POST /number         — prepend line numbers
  POST /dedup          — remove duplicate lines
  POST /sort           — sort lines (asc/desc)
  POST /trim           — trim whitespace from lines (left/right/both)
  POST /pad            — pad text to width (left/right/center)
  POST /reverse        — reverse characters
  POST /reverse-lines  — reverse line order
  POST /head           — first N lines
  POST /tail           — last N lines
  POST /extract        — extract lines by range (start, end)
  POST /find           — find lines containing a string
  POST /replace        — replace all occurrences
  POST /squeeze        — collapse consecutive blank lines
  POST /grep           — filter lines matching a string
  POST /unique         — remove duplicate lines (alias of dedup)
  POST /join           — join lines with a separator
  POST /mcp            — MCP JSON-RPC 2.0 endpoint

PARAMETERS (query string):
  width=N    — wrap width (default 80)
  start=N    — line number start (default 1)
  end=N      — line number end
  n=N        — number of lines (head/tail)
  desc=true  — sort descending
  mode=left  — trim mode (left, right, both)
  align=left — pad alignment (left, right, center)
  fill=X     — pad fill character (default space)
  q=STRING   — search string (find, grep)
  old=STRING — text to replace
  new=STRING — replacement text
  sep=STRING — line separator (join)
  invert=true — invert grep match
  type=X     — count type (words, chars, bytes, lines)

RESPONSE FORMAT:
  Plain text by default. JSON via Accept: application/json or ?format=json.

EXAMPLES:
  curl -X POST localhost:8080/stats -d "Hello world."
  curl -X POST "localhost:8080/wrap?width=20" -d "The quick brown fox"
  curl -X POST "localhost:8080/head?n=3" -d "$(cat file.txt)"
  curl -X POST "localhost:8080/sort?desc=true" -d "banana\napple\ncherry"
  curl -X POST "localhost:8080/replace?old=foo&new=bar" -d "foo bar foo"
`
	writeText(w, text)
}

func (h *Handler) stats(w http.ResponseWriter, r *http.Request) {
	body, err := getBody(r)
	if err != nil {
		writeError(w, 400, "cannot read body", "send text in the request body")
		return
	}
	s := model.ComputeStats(body)
	if wantsJSON(r) {
		writeJSON(w, map[string]int{
			"chars":         s.Chars,
			"chars_no_space": s.CharsNoSpace,
			"words":          s.Words,
			"lines":          s.Lines,
			"sentences":      s.Sentences,
			"paragraphs":     s.Paragraphs,
			"bytes":          s.Bytes,
			"runes":          s.Runes,
			"reading_min":    model.ReadingTime(s.Words),
			"speaking_min":   model.SpeakingTime(s.Words),
		})
		return
	}
	fmt.Fprintf(w, "chars=%d chars_no_space=%d words=%d lines=%d sentences=%d paragraphs=%d bytes=%d runes=%d reading_min=%d speaking_min=%d\n",
		s.Chars, s.CharsNoSpace, s.Words, s.Lines, s.Sentences, s.Paragraphs, s.Bytes, s.Runes,
		model.ReadingTime(s.Words), model.SpeakingTime(s.Words))
}

func (h *Handler) count(w http.ResponseWriter, r *http.Request) {
	body, err := getBody(r)
	if err != nil {
		writeError(w, 400, "cannot read body", "send text in the request body")
		return
	}
	countType := getQueryParam(r, "type")
	if countType == "" {
		countType = "words"
	}
	var n int
	switch countType {
	case "words":
		n = model.WordCount(body)
	case "chars":
		n = model.CountChars(body)
	case "bytes":
		n = model.CountBytes(body)
	case "lines":
		n = model.CountLines(body)
	default:
		writeError(w, 400, "unknown count type", "use type=words, type=chars, type=bytes, or type=lines")
		return
	}
	if wantsJSON(r) {
		writeJSON(w, map[string]int{countType: n})
		return
	}
	fmt.Fprintf(w, "%s=%d\n", countType, n)
}

func (h *Handler) wrap(w http.ResponseWriter, r *http.Request) {
	body, err := getBody(r)
	if err != nil {
		writeError(w, 400, "cannot read body", "send text in the request body")
		return
	}
	width := getIntParam(r, "width", 80)
	result := model.Wrap(body, width)
	if wantsJSON(r) {
		writeJSON(w, map[string]string{"result": result})
		return
	}
	writeText(w, result)
}

func (h *Handler) numberLines(w http.ResponseWriter, r *http.Request) {
	body, err := getBody(r)
	if err != nil {
		writeError(w, 400, "cannot read body", "send text in the request body")
		return
	}
	start := getIntParam(r, "start", 1)
	result := model.NumberLines(body, start)
	if wantsJSON(r) {
		writeJSON(w, map[string]string{"result": result})
		return
	}
	writeText(w, result)
}

func (h *Handler) dedup(w http.ResponseWriter, r *http.Request) {
	body, err := getBody(r)
	if err != nil {
		writeError(w, 400, "cannot read body", "send text in the request body")
		return
	}
	result := model.Dedup(body)
	if wantsJSON(r) {
		writeJSON(w, map[string]string{"result": result})
		return
	}
	writeText(w, result)
}

func (h *Handler) unique(w http.ResponseWriter, r *http.Request) {
	h.dedup(w, r)
}

func (h *Handler) sortLines(w http.ResponseWriter, r *http.Request) {
	body, err := getBody(r)
	if err != nil {
		writeError(w, 400, "cannot read body", "send text in the request body")
		return
	}
	desc := getQueryParam(r, "desc") == "true"
	result := model.SortLines(body, desc)
	if wantsJSON(r) {
		writeJSON(w, map[string]string{"result": result})
		return
	}
	writeText(w, result)
}

func (h *Handler) trimLines(w http.ResponseWriter, r *http.Request) {
	body, err := getBody(r)
	if err != nil {
		writeError(w, 400, "cannot read body", "send text in the request body")
		return
	}
	mode := getQueryParam(r, "mode")
	if mode == "" {
		mode = "both"
	}
	if mode != "left" && mode != "right" && mode != "both" {
		writeError(w, 400, "invalid mode", "use mode=left, mode=right, or mode=both")
		return
	}
	result := model.TrimLines(body, mode)
	if wantsJSON(r) {
		writeJSON(w, map[string]string{"result": result})
		return
	}
	writeText(w, result)
}

func (h *Handler) pad(w http.ResponseWriter, r *http.Request) {
	body, err := getBody(r)
	if err != nil {
		writeError(w, 400, "cannot read body", "send text in the request body")
		return
	}
	width := getIntParam(r, "width", 80)
	fill := getQueryParam(r, "fill")
	if fill == "" {
		fill = " "
	}
	align := getQueryParam(r, "align")
	if align == "" {
		align = "left"
	}
	if align != "left" && align != "right" && align != "center" {
		writeError(w, 400, "invalid align", "use align=left, align=right, or align=center")
		return
	}
	result := model.Pad(body, width, fill, align)
	if wantsJSON(r) {
		writeJSON(w, map[string]string{"result": result})
		return
	}
	writeText(w, result)
}

func (h *Handler) reverse(w http.ResponseWriter, r *http.Request) {
	body, err := getBody(r)
	if err != nil {
		writeError(w, 400, "cannot read body", "send text in the request body")
		return
	}
	result := model.Reverse(body)
	if wantsJSON(r) {
		writeJSON(w, map[string]string{"result": result})
		return
	}
	writeText(w, result)
}

func (h *Handler) reverseLines(w http.ResponseWriter, r *http.Request) {
	body, err := getBody(r)
	if err != nil {
		writeError(w, 400, "cannot read body", "send text in the request body")
		return
	}
	result := model.ReverseLines(body)
	if wantsJSON(r) {
		writeJSON(w, map[string]string{"result": result})
		return
	}
	writeText(w, result)
}

func (h *Handler) head(w http.ResponseWriter, r *http.Request) {
	body, err := getBody(r)
	if err != nil {
		writeError(w, 400, "cannot read body", "send text in the request body")
		return
	}
	n := getIntParam(r, "n", 10)
	result := model.Head(body, n)
	if wantsJSON(r) {
		writeJSON(w, map[string]string{"result": result})
		return
	}
	writeText(w, result)
}

func (h *Handler) tail(w http.ResponseWriter, r *http.Request) {
	body, err := getBody(r)
	if err != nil {
		writeError(w, 400, "cannot read body", "send text in the request body")
		return
	}
	n := getIntParam(r, "n", 10)
	result := model.Tail(body, n)
	if wantsJSON(r) {
		writeJSON(w, map[string]string{"result": result})
		return
	}
	writeText(w, result)
}

func (h *Handler) extractLines(w http.ResponseWriter, r *http.Request) {
	body, err := getBody(r)
	if err != nil {
		writeError(w, 400, "cannot read body", "send text in the request body")
		return
	}
	start := getIntParam(r, "start", 1)
	end := getIntParam(r, "end", 1)
	result := model.ExtractLines(body, start, end)
	if wantsJSON(r) {
		writeJSON(w, map[string]string{"result": result})
		return
	}
	writeText(w, result)
}

func (h *Handler) find(w http.ResponseWriter, r *http.Request) {
	body, err := getBody(r)
	if err != nil {
		writeError(w, 400, "cannot read body", "send text in the request body")
		return
	}
	q := getQueryParam(r, "q")
	if q == "" {
		writeError(w, 400, "missing search query", "add ?q=search_string to the URL")
		return
	}
	matches := model.Find(body, q)
	if wantsJSON(r) {
		writeJSON(w, map[string]interface{}{"matches": matches, "count": len(matches)})
		return
	}
	if len(matches) == 0 {
		fmt.Fprintln(w, "matches=0")
		return
	}
	fmt.Fprintf(w, "matches=%d lines=%s\n", len(matches), strings.Trim(strings.Join(strings.Fields(fmt.Sprint(matches)), ","), "[]"))
}

func (h *Handler) replace(w http.ResponseWriter, r *http.Request) {
	body, err := getBody(r)
	if err != nil {
		writeError(w, 400, "cannot read body", "send text in the request body")
		return
	}
	old := getQueryParam(r, "old")
	new := getQueryParam(r, "new")
	if old == "" {
		writeError(w, 400, "missing old parameter", "add ?old=text_to_replace to the URL")
		return
	}
	result := model.Replace(body, old, new)
	if wantsJSON(r) {
		writeJSON(w, map[string]string{"result": result})
		return
	}
	writeText(w, result)
}

func (h *Handler) squeeze(w http.ResponseWriter, r *http.Request) {
	body, err := getBody(r)
	if err != nil {
		writeError(w, 400, "cannot read body", "send text in the request body")
		return
	}
	result := model.SqueezeBlankLines(body)
	if wantsJSON(r) {
		writeJSON(w, map[string]string{"result": result})
		return
	}
	writeText(w, result)
}

func (h *Handler) grep(w http.ResponseWriter, r *http.Request) {
	body, err := getBody(r)
	if err != nil {
		writeError(w, 400, "cannot read body", "send text in the request body")
		return
	}
	q := getQueryParam(r, "q")
	if q == "" {
		writeError(w, 400, "missing search query", "add ?q=search_string to the URL")
		return
	}
	invert := getQueryParam(r, "invert") == "true"
	result := model.Grep(body, q, invert)
	if wantsJSON(r) {
		writeJSON(w, map[string]string{"result": result})
		return
	}
	writeText(w, result)
}

func (h *Handler) joinLines(w http.ResponseWriter, r *http.Request) {
	body, err := getBody(r)
	if err != nil {
		writeError(w, 400, "cannot read body", "send text in the request body")
		return
	}
	sep := getQueryParam(r, "sep")
	if sep == "" {
		sep = " "
	}
	result := model.JoinLines(body, sep)
	if wantsJSON(r) {
		writeJSON(w, map[string]string{"result": result})
		return
	}
	writeText(w, result)
}

func (h *Handler) notFound(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		writeText(w, "textkit — text analysis and manipulation service. GET /help for usage.")
		return
	}
	writeError(w, 404, "endpoint not found", "GET /help to see available endpoints")
}
