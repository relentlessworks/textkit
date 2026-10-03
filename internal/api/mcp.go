package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/relentlessworks/textkit/internal/model"
)

type mcpRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
}

type mcpResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *mcpError   `json:"error,omitempty"`
}

type mcpError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type mcpTool struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	InputSchema map[string]any `json:"inputSchema"`
}

func (h *Handler) mcp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, 405, "method not allowed", "POST JSON-RPC 2.0 requests to /mcp")
		return
	}

	var req mcpRequest
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, 400, "cannot read body", "send a JSON-RPC 2.0 request")
		return
	}
	if err := json.Unmarshal(body, &req); err != nil {
		writeError(w, 400, "invalid JSON-RPC", "send a valid JSON-RPC 2.0 request")
		return
	}

	w.Header().Set("Content-Type", "application/json")

	switch req.Method {
	case "initialize":
		json.NewEncoder(w).Encode(mcpResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: map[string]any{
				"protocolVersion": "2024-11-05",
				"serverInfo":      map[string]string{"name": "textkit", "version": "0.1.0"},
				"capabilities":   map[string]any{"tools": map[string]any{}},
			},
		})

	case "tools/list":
		json.NewEncoder(w).Encode(mcpResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: map[string]any{"tools:": mcpTools()},
		})

	case "tools/call":
		var params struct {
			Name      string         `json:"name"`
			Arguments map[string]any `json:"arguments"`
		}
		if err := json.Unmarshal(req.Params, &params); err != nil {
			json.NewEncoder(w).Encode(mcpResponse{
				JSONRPC: "2.0", ID: req.ID,
				Error: &mcpError{Code: -32602, Message: "invalid params"},
			})
			return
		}
		result := h.callMCPTool(params.Name, params.Arguments)
		json.NewEncoder(w).Encode(mcpResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: map[string]any{
				"content": []map[string]any{
					{"type": "text", "text": result},
				},
			},
		})

	default:
		json.NewEncoder(w).Encode(mcpResponse{
			JSONRPC: "2.0", ID: req.ID,
			Error: &mcpError{Code: -32601, Message: "method not found"},
		})
	}
}

func mcpTools() []mcpTool {
	strSchema := func(desc string) map[string]any {
		return map[string]any{
			"type": "object",
			"properties": map[string]any{
				"text": map[string]any{"type": "string", "description": desc},
			},
			"required": []string{"text"},
		}
	}
	strSchemaWithExtra := func(desc string, extra map[string]any) map[string]any {
		props := map[string]any{
			"text": map[string]any{"type": "string", "description": desc},
		}
		for k, v := range extra {
			props[k] = v
		}
		return map[string]any{
			"type":       "object",
			"properties":  props,
			"required":   []string{"text"},
		}
	}

	return []mcpTool{
		{Name: "stats", Description: "Get text statistics (words, chars, lines, sentences, paragraphs, bytes, runes, reading/speaking time)", InputSchema: strSchema("Text to analyze")},
		{Name: "count", Description: "Count words, chars, bytes, or lines", InputSchema: strSchemaWithExtra("Text to count", map[string]any{
			"type": map[string]any{"type": "string", "description": "Count type: words, chars, bytes, lines (default: words)"},
		})},
		{Name: "wrap", Description: "Wrap text to a given width", InputSchema: strSchemaWithExtra("Text to wrap", map[string]any{
			"width": map[string]any{"type": "integer", "description": "Wrap width (default: 80)"},
		})},
		{Name: "number", Description: "Prepend line numbers", InputSchema: strSchemaWithExtra("Text to number", map[string]any{
			"start": map[string]any{"type": "integer", "description": "Starting line number (default: 1)"},
		})},
		{Name: "dedup", Description: "Remove duplicate lines", InputSchema: strSchema("Text to deduplicate")},
		{Name: "sort", Description: "Sort lines ascending or descending", InputSchema: strSchemaWithExtra("Text to sort", map[string]any{
			"desc": map[string]any{"type": "boolean", "description": "Sort descending (default: false)"},
		})},
		{Name: "trim", Description: "Trim whitespace from lines", InputSchema: strSchemaWithExtra("Text to trim", map[string]any{
			"mode": map[string]any{"type": "string", "description": "Trim mode: left, right, both (default: both)"},
		})},
		{Name: "pad", Description: "Pad text to a width", InputSchema: strSchemaWithExtra("Text to pad", map[string]any{
			"width": map[string]any{"type": "integer", "description": "Target width (default: 80)"},
			"align": map[string]any{"type": "string", "description": "Alignment: left, right, center (default: left)"},
			"fill":  map[string]any{"type": "string", "description": "Fill character (default: space)"},
		})},
		{Name: "reverse", Description: "Reverse characters in text", InputSchema: strSchema("Text to reverse")},
		{Name: "reverse_lines", Description: "Reverse line order", InputSchema: strSchema("Text to reverse lines")},
		{Name: "head", Description: "Get first N lines", InputSchema: strSchemaWithExtra("Text", map[string]any{
			"n": map[string]any{"type": "integer", "description": "Number of lines (default: 10)"},
		})},
		{Name: "tail", Description: "Get last N lines", InputSchema: strSchemaWithExtra("Text", map[string]any{
			"n": map[string]any{"type": "integer", "description": "Number of lines (default: 10)"},
		})},
		{Name: "extract", Description: "Extract lines by range", InputSchema: strSchemaWithExtra("Text", map[string]any{
			"start": map[string]any{"type": "integer", "description": "Start line (1-indexed, default: 1)"},
			"end":   map[string]any{"type": "integer", "description": "End line (1-indexed, default: 1)"},
		})},
		{Name: "find", Description: "Find lines containing a string", InputSchema: strSchemaWithExtra("Text to search", map[string]any{
			"q": map[string]any{"type": "string", "description": "Search string"},
		})},
		{Name: "replace", Description: "Replace all occurrences", InputSchema: strSchemaWithExtra("Text", map[string]any{
			"old": map[string]any{"type": "string", "description": "Text to find"},
			"new": map[string]any{"type": "string", "description": "Replacement text"},
		})},
		{Name: "squeeze", Description: "Collapse consecutive blank lines", InputSchema: strSchema("Text to squeeze")},
		{Name: "grep", Description: "Filter lines matching a string", InputSchema: strSchemaWithExtra("Text to filter", map[string]any{
			"q":      map[string]any{"type": "string", "description": "Search string"},
			"invert": map[string]any{"type": "boolean", "description": "Invert match (default: false)"},
		})},
		{Name: "join", Description: "Join lines with a separator", InputSchema: strSchemaWithExtra("Text", map[string]any{
			"sep": map[string]any{"type": "string", "description": "Separator (default: space)"},
		})},
	}
}

func (h *Handler) callMCPTool(name string, args map[string]any) string {
	text, _ := args["text"].(string)
	getStr := func(k string) string {
		v, _ := args[k].(string)
		return v
	}
	getInt := func(k string, def int) int {
		v, ok := args[k].(float64)
		if !ok {
			return def
		}
		return int(v)
	}
	getBool := func(k string) bool {
		v, _ := args[k].(bool)
		return v
	}

	switch name {
	case "stats":
		s := model.ComputeStats(text)
		return fmt.Sprintf("chars=%d chars_no_space=%d words=%d lines=%d sentences=%d paragraphs=%d bytes=%d runes=%d reading_min=%d speaking_min=%d",
			s.Chars, s.CharsNoSpace, s.Words, s.Lines, s.Sentences, s.Paragraphs, s.Bytes, s.Runes,
			model.ReadingTime(s.Words), model.SpeakingTime(s.Words))
	case "count":
		ct := getStr("type")
		if ct == "" {
			ct = "words"
		}
		var n int
		switch ct {
		case "words":
			n = model.WordCount(text)
		case "chars":
			n = model.CountChars(text)
		case "bytes":
			n = model.CountBytes(text)
		case "lines":
			n = model.CountLines(text)
		default:
			return "error: unknown count type | hint: use type=words, type=chars, type=bytes, or type=lines"
		}
		return fmt.Sprintf("%s=%d", ct, n)
	case "wrap":
		return model.Wrap(text, getInt("width", 80))
	case "number":
		return model.NumberLines(text, getInt("start", 1))
	case "dedup":
		return model.Dedup(text)
	case "sort":
		return model.SortLines(text, getBool("desc"))
	case "trim":
		mode := getStr("mode")
		if mode == "" {
			mode = "both"
		}
		return model.TrimLines(text, mode)
	case "pad":
		fill := getStr("fill")
		if fill == "" {
			fill = " "
		}
		align := getStr("align")
		if align == "" {
			align = "left"
		}
		return model.Pad(text, getInt("width", 80), fill, align)
	case "reverse":
		return model.Reverse(text)
	case "reverse_lines":
		return model.ReverseLines(text)
	case "head":
		return model.Head(text, getInt("n", 10))
	case "tail":
		return model.Tail(text, getInt("n", 10))
	case "extract":
		return model.ExtractLines(text, getInt("start", 1), getInt("end", 1))
	case "find":
		q := getStr("q")
		if q == "" {
			return "error: missing search query | hint: provide q parameter"
		}
		matches := model.Find(text, q)
		return fmt.Sprintf("matches=%d lines=%s", len(matches), strings.Trim(strings.Join(strings.Fields(fmt.Sprint(matches)), ","), "[]"))
	case "replace":
		old := getStr("old")
		if old == "" {
			return "error: missing old parameter | hint: provide old parameter"
		}
		return model.Replace(text, old, getStr("new"))
	case "squeeze":
		return model.SqueezeBlankLines(text)
	case "grep":
		q := getStr("q")
		if q == "" {
			return "error: missing search query | hint: provide q parameter"
		}
		return model.Grep(text, q, getBool("invert"))
	case "join":
		sep := getStr("sep")
		if sep == "" {
			sep = " "
		}
		return model.JoinLines(text, sep)
	default:
		return "error: unknown tool | hint: call tools/list to see available tools"
	}
}
