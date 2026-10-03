package model

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// Stats holds text statistics.
type Stats struct {
	Chars        int
	CharsNoSpace int
	Words        int
	Lines        int
	Sentences    int
	Paragraphs   int
	Bytes        int
	Runes        int
}

// ComputeStats computes text statistics.
func ComputeStats(text string) Stats {
	s := Stats{
		Bytes: len(text),
		Runes: utf8.RuneCountInString(text),
		Chars: utf8.RuneCountInString(text),
	}

	for _, r := range text {
		if !unicode.IsSpace(r) {
			s.CharsNoSpace++
		}
	}

	// Words
	s.Words = WordCount(text)

	// Lines
	s.Lines = strings.Count(text, "\n")
	if len(text) > 0 && !strings.HasSuffix(text, "\n") {
		s.Lines++
	}

	// Sentences (count sentence-ending punctuation)
	s.Sentences = countSentences(text)

	// Paragraphs (blocks separated by blank lines)
	s.Paragraphs = countParagraphs(text)

	return s
}

// WordCount counts words in text.
func WordCount(text string) int {
	return len(strings.Fields(text))
}

func countSentences(text string) int {
	count := 0
	inSentence := false
	for _, r := range text {
		if isSentenceEnd(r) {
			if inSentence {
				count++
				inSentence = false
			}
		} else if !unicode.IsSpace(r) {
			inSentence = true
		}
	}
	if inSentence {
		count++
	}
	return count
}

func isSentenceEnd(r rune) bool {
	return r == '.' || r == '!' || r == '?' || r == '。' || r == '！' || r == '？'
}

func countParagraphs(text string) int {
	paragraphs := 0
	inPara := false
	for _, line := range strings.Split(text, "\n") {
		if strings.TrimSpace(line) == "" {
			if inPara {
				paragraphs++
				inPara = false
			}
		} else {
			inPara = true
		}
	}
	if inPara {
		paragraphs++
	}
	if paragraphs == 0 && strings.TrimSpace(text) != "" {
		paragraphs = 1
	}
	return paragraphs
}

// ReadingTime returns reading time in minutes (200 wpm).
func ReadingTime(words int) int {
	minutes := words / 200
	if words%200 > 0 {
		minutes++
	}
	if minutes == 0 {
		minutes = 1
	}
	return minutes
}

// SpeakingTime returns speaking time in minutes (130 wpm).
func SpeakingTime(words int) int {
	minutes := words / 130
	if words%130 > 0 {
		minutes++
	}
	if minutes == 0 {
		minutes = 1
	}
	return minutes
}

// Wrap wraps text to a given width.
func Wrap(text string, width int) string {
	if width <= 0 {
		width = 80
	}
	var result strings.Builder
	for _, line := range strings.Split(text, "\n") {
		wrapped := wrapLine(line, width)
		if result.Len() > 0 {
			result.WriteString("\n")
		}
		result.WriteString(wrapped)
	}
	return result.String()
}

func wrapLine(line string, width int) string {
	if len(line) <= width {
		return line
	}
	var result strings.Builder
	words := strings.Fields(line)
	if len(words) == 0 {
		return line
	}
	lineLen := 0
	for i, word := range words {
		wLen := len(word)
		if i == 0 {
			result.WriteString(word)
			lineLen = wLen
		} else if lineLen+1+wLen <= width {
			result.WriteString(" ")
			result.WriteString(word)
			lineLen += 1 + wLen
		} else {
			result.WriteString("\n")
			result.WriteString(word)
			lineLen = wLen
		}
	}
	return result.String()
}

// NumberLines prepends line numbers to each line.
func NumberLines(text string, start int) string {
	lines := strings.Split(text, "\n")
	var result strings.Builder
	for i, line := range lines {
		if result.Len() > 0 {
			result.WriteString("\n")
		}
		result.WriteString(itoaPad(i+start, len(itoa(len(lines)+start-1))))
		result.WriteString("  ")
		result.WriteString(line)
	}
	return result.String()
}

func itoaPad(n, width int) string {
	s := itoa(n)
	for len(s) < width {
		s = " " + s
	}
	return s
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := false
	if n < 0 {
		neg = true
		n = -n
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}

// Dedup removes duplicate lines, keeping first occurrence.
func Dedup(text string) string {
	seen := make(map[string]bool)
	lines := strings.Split(text, "\n")
	var result strings.Builder
	for _, line := range lines {
		if !seen[line] {
			seen[line] = true
			if result.Len() > 0 {
				result.WriteString("\n")
			}
			result.WriteString(line)
		}
	}
	return result.String()
}

// SortLines sorts lines ascending or descending.
func SortLines(text string, descending bool) string {
	lines := strings.Split(text, "\n")
	// Simple sort using insertion sort to avoid importing sort
	for i := 1; i < len(lines); i++ {
		for j := i; j > 0; j-- {
			if descending {
				if lines[j] > lines[j-1] {
					lines[j], lines[j-1] = lines[j-1], lines[j]
				} else {
					break
				}
			} else {
				if lines[j] < lines[j-1] {
					lines[j], lines[j-1] = lines[j-1], lines[j]
				} else {
					break
				}
			}
		}
	}
	return strings.Join(lines, "\n")
}

// TrimLines trims leading/trailing whitespace from each line.
func TrimLines(text string, mode string) string {
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		switch mode {
		case "left":
			lines[i] = strings.TrimLeft(line, " \t\r")
		case "right":
			lines[i] = strings.TrimRight(line, " \t\r")
		case "both", "":
			lines[i] = strings.TrimSpace(line)
		}
	}
	return strings.Join(lines, "\n")
}

// Pad pads text to a given width with optional fill character.
func Pad(text string, width int, fill string, align string) string {
	if fill == "" {
		fill = " "
	}
	f := fill[0]
	textLen := utf8.RuneCountInString(text)
	if textLen >= width {
		return text
	}
	padLen := width - textLen
	padding := strings.Repeat(string(f), padLen)
	switch align {
	case "right":
		return padding + text
	case "center":
		left := padLen / 2
		right := padLen - left
		return strings.Repeat(string(f), left) + text + strings.Repeat(string(f), right)
	default: // "left"
		return text + padding
	}
}

// Reverse reverses a string (rune-aware).
func Reverse(text string) string {
	runes := []rune(text)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// ReverseLines reverses the order of lines.
func ReverseLines(text string) string {
	lines := strings.Split(text, "\n")
	for i, j := 0, len(lines)-1; i < j; i, j = i+1, j-1 {
		lines[i], lines[j] = lines[j], lines[i]
	}
	return strings.Join(lines, "\n")
}

// Head returns the first n lines.
func Head(text string, n int) string {
	if n <= 0 {
		return ""
	}
	lines := strings.Split(text, "\n")
	if n > len(lines) {
		n = len(lines)
	}
	return strings.Join(lines[:n], "\n")
}

// Tail returns the last n lines.
func Tail(text string, n int) string {
	if n <= 0 {
		return ""
	}
	lines := strings.Split(text, "\n")
	if n > len(lines) {
		n = len(lines)
	}
	return strings.Join(lines[len(lines)-n:], "\n")
}

// ExtractLines returns lines from start to end (1-indexed, inclusive).
func ExtractLines(text string, start, end int) string {
	lines := strings.Split(text, "\n")
	if start < 1 {
		start = 1
	}
	if end > len(lines) {
		end = len(lines)
	}
	if start > end {
		return ""
	}
	return strings.Join(lines[start-1:end], "\n")
}

// Find returns line numbers that contain the search string.
func Find(text string, search string) []int {
	if search == "" {
		return nil
	}
	lines := strings.Split(text, "\n")
	var matches []int
	for i, line := range lines {
		if strings.Contains(line, search) {
			matches = append(matches, i+1)
		}
	}
	return matches
}

// Replace replaces all occurrences of old with new in text.
func Replace(text, old, new string) string {
	if old == "" {
		return text
	}
	return strings.ReplaceAll(text, old, new)
}

// CountLines returns the number of lines.
func CountLines(text string) int {
	if text == "" {
		return 0
	}
	n := strings.Count(text, "\n")
	if !strings.HasSuffix(text, "\n") {
		n++
	}
	return n
}

// CountChars returns the number of characters (runes).
func CountChars(text string) int {
	return utf8.RuneCountInString(text)
}

// CountBytes returns the number of bytes.
func CountBytes(text string) int {
	return len(text)
}

// SqueezeBlankLines collapses multiple consecutive blank lines into one.
func SqueezeBlankLines(text string) string {
	lines := strings.Split(text, "\n")
	var result []string
	prevBlank := false
	for _, line := range lines {
		blank := strings.TrimSpace(line) == ""
		if blank && prevBlank {
			continue
		}
		result = append(result, line)
		prevBlank = blank
	}
	return strings.Join(result, "\n")
}

// Grep filters lines matching a substring.
func Grep(text, pattern string, invert bool) string {
	lines := strings.Split(text, "\n")
	var result []string
	for _, line := range lines {
		match := strings.Contains(line, pattern)
		if invert {
			match = !match
		}
		if match {
			result = append(result, line)
		}
	}
	return strings.Join(result, "\n")
}

// Unique returns only unique lines (same as Dedup but sorted output).
func Unique(text string) string {
	return Dedup(text)
}

// JoinLines joins lines with a custom separator.
func JoinLines(text, sep string) string {
	lines := strings.Split(text, "\n")
	return strings.Join(lines, sep)
}

// Center centers text in a field of given width.
func Center(text string, width int, fill string) string {
	return Pad(text, width, fill, "center")
}
