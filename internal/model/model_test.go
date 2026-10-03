package model

import "testing"

func TestWordCount(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{"", 0},
		{"hello", 1},
		{"hello world", 2},
		{"  hello   world  ", 2},
		{"one two three four five", 5},
	}
	for _, tt := range tests {
		got := WordCount(tt.input)
		if got != tt.want {
			t.Errorf("WordCount(%q) = %d, want %d", tt.input, got, tt.want)
		}
	}
}

func TestComputeStats(t *testing.T) {
	text := "Hello world.\nThis is a test.\n\nAnother paragraph."
	s := ComputeStats(text)
	if s.Words != 8 {
		t.Errorf("Words = %d, want 8", s.Words)
	}
	if s.Lines != 4 {
		t.Errorf("Lines = %d, want 4", s.Lines)
	}
	if s.Sentences != 3 {
		t.Errorf("Sentences = %d, want 3", s.Sentences)
	}
	if s.Paragraphs != 2 {
		t.Errorf("Paragraphs = %d, want 2", s.Paragraphs)
	}
}

func TestWrap(t *testing.T) {
	text := "The quick brown fox jumps over the lazy dog"
	wrapped := Wrap(text, 20)
	lines := splitLines(wrapped)
	for _, line := range lines {
		if len(line) > 20 {
			t.Errorf("Line too long: %q (%d chars)", line, len(line))
		}
	}
}

func TestNumberLines(t *testing.T) {
	text := "a\nb\nc"
	result := NumberLines(text, 1)
	lines := splitLines(result)
	if len(lines) != 3 {
		t.Fatalf("Expected 3 lines, got %d", len(lines))
	}
	if lines[0] != "1  a" {
		t.Errorf("First line = %q, want %q", lines[0], "1  a")
	}
}

func TestDedup(t *testing.T) {
	text := "a\nb\na\nc\nb"
	result := Dedup(text)
	if result != "a\nb\nc" {
		t.Errorf("Dedup = %q, want %q", result, "a\nb\nc")
	}
}

func TestSortLines(t *testing.T) {
	text := "banana\napple\ncherry"
	result := SortLines(text, false)
	if result != "apple\nbanana\ncherry" {
		t.Errorf("SortLines asc = %q, want %q", result, "apple\nbanana\ncherry")
	}
	result = SortLines(text, true)
	if result != "cherry\nbanana\napple" {
		t.Errorf("SortLines desc = %q, want %q", result, "cherry\nbanana\napple")
	}
}

func TestTrimLines(t *testing.T) {
	text := "  hello  \n  world  "
	result := TrimLines(text, "both")
	if result != "hello\nworld" {
		t.Errorf("TrimLines = %q, want %q", result, "hello\nworld")
	}
}

func TestPad(t *testing.T) {
	result := Pad("hi", 5, " ", "left")
	if result != "hi   " {
		t.Errorf("Pad left = %q, want %q", result, "hi   ")
	}
	result = Pad("hi", 5, " ", "right")
	if result != "   hi" {
		t.Errorf("Pad right = %q, want %q", result, "   hi")
	}
	result = Pad("hi", 6, " ", "center")
	if result != "  hi  " {
		t.Errorf("Pad center = %q, want %q", result, "  hi  ")
	}
}

func TestReverse(t *testing.T) {
	result := Reverse("hello")
	if result != "olleh" {
		t.Errorf("Reverse = %q, want %q", result, "olleh")
	}
	result = Reverse("héllo")
	if result != "olléh" {
		t.Errorf("Reverse unicode = %q, want %q", result, "olléh")
	}
}

func TestReverseLines(t *testing.T) {
	text := "a\nb\nc"
	result := ReverseLines(text)
	if result != "c\nb\na" {
		t.Errorf("ReverseLines = %q, want %q", result, "c\nb\na")
	}
}

func TestHead(t *testing.T) {
	text := "a\nb\nc\nd\ne"
	result := Head(text, 2)
	if result != "a\nb" {
		t.Errorf("Head = %q, want %q", result, "a\nb")
	}
}

func TestTail(t *testing.T) {
	text := "a\nb\nc\nd\ne"
	result := Tail(text, 2)
	if result != "d\ne" {
		t.Errorf("Tail = %q, want %q", result, "d\ne")
	}
}

func TestExtractLines(t *testing.T) {
	text := "a\nb\nc\nd\ne"
	result := ExtractLines(text, 2, 4)
	if result != "b\nc\nd" {
		t.Errorf("ExtractLines = %q, want %q", result, "b\nc\nd")
	}
}

func TestFind(t *testing.T) {
	text := "hello world\nfoo bar\nhello there"
	matches := Find(text, "hello")
	if len(matches) != 2 {
		t.Fatalf("Expected 2 matches, got %d", len(matches))
	}
	if matches[0] != 1 || matches[1] != 3 {
		t.Errorf("Matches = %v, want [1 3]", matches)
	}
}

func TestReplace(t *testing.T) {
	result := Replace("hello world", "world", "there")
	if result != "hello there" {
		t.Errorf("Replace = %q, want %q", result, "hello there")
	}
}

func TestCountLines(t *testing.T) {
	if CountLines("a\nb\nc") != 3 {
		t.Errorf("CountLines = %d, want 3", CountLines("a\nb\nc"))
	}
	if CountLines("a\nb\nc\n") != 3 {
		t.Errorf("CountLines trailing = %d, want 3", CountLines("a\nb\nc\n"))
	}
	if CountLines("") != 0 {
		t.Errorf("CountLines empty = %d, want 0", CountLines(""))
	}
}

func TestSqueezeBlankLines(t *testing.T) {
	text := "a\n\n\n\nb\n\n\nc"
	result := SqueezeBlankLines(text)
	if result != "a\n\nb\n\nc" {
		t.Errorf("SqueezeBlankLines = %q, want %q", result, "a\n\nb\n\nc")
	}
}

func TestGrep(t *testing.T) {
	text := "hello world\nfoo bar\nhello there"
	result := Grep(text, "hello", false)
	if result != "hello world\nhello there" {
		t.Errorf("Grep = %q, want %q", result, "hello world\nhello there")
	}
	result = Grep(text, "hello", true)
	if result != "foo bar" {
		t.Errorf("Grep invert = %q, want %q", result, "foo bar")
	}
}

func TestJoinLines(t *testing.T) {
	text := "a\nb\nc"
	result := JoinLines(text, ", ")
	if result != "a, b, c" {
		t.Errorf("JoinLines = %q, want %q", result, "a, b, c")
	}
}

func TestReadingTime(t *testing.T) {
	if ReadingTime(200) != 1 {
		t.Errorf("ReadingTime(200) = %d, want 1", ReadingTime(200))
	}
	if ReadingTime(201) != 2 {
		t.Errorf("ReadingTime(201) = %d, want 2", ReadingTime(201))
	}
}

func TestSpeakingTime(t *testing.T) {
	if SpeakingTime(130) != 1 {
		t.Errorf("SpeakingTime(130) = %d, want 1", SpeakingTime(130))
	}
	if SpeakingTime(131) != 2 {
		t.Errorf("SpeakingTime(131) = %d, want 2", SpeakingTime(131))
	}
}

func splitLines(s string) []string {
	if s == "" {
		return nil
	}
	var result []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			result = append(result, s[start:i])
			start = i + 1
		}
	}
	result = append(result, s[start:])
	return result
}
