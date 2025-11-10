package format

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

const (
	ansiReset  = "\x1b[0m"
	ansiBold   = "\x1b[1m"
	ansiDim    = "\x1b[2m"
	ansiRed    = "\x1b[31m"
	ansiYellow = "\x1b[33m"
	ansiGreen  = "\x1b[32m"
	ansiBlue   = "\x1b[34m"
)

func HumanBytes(n int64) string {
	const unit = 1024
	if n < unit {
		return fmt.Sprintf("%d B", n)
	}
	div, exp := int64(unit), 0
	for n >= unit && exp < 6 {
		n /= div
		exp++
	}
	return fmt.Sprintf("%d %ciB", n, "KMGTPE"[exp-1])
}

func PrintTable(w io.Writer, headers []string, rows [][]string) {
	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = len(h)
	}
	for _, r := range rows {
		for i, c := range r {
			if len(c) > widths[i] {
				widths[i] = len(c)
			}
		}
	}
	// header
	for i, h := range headers {
		fmt.Fprintf(w, "%-*s", widths[i]+2, h)
	}
	fmt.Fprintln(w)
	// sep
	for _, ww := range widths {
		fmt.Fprint(w, strings.Repeat("-", ww), "  ")
	}
	fmt.Fprintln(w)
	// rows
	for _, r := range rows {
		for i, c := range r {
			fmt.Fprintf(w, "%-*s", widths[i]+2, c)
		}
		fmt.Fprintln(w)
	}
}

func PrintJSON(w io.Writer, v any) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

// Bold returns bold text if terminal supports ANSI (most do).
func Bold(s string) string { return ansiBold + s + ansiReset }

// Dim returns dim/secondary-colored text.
func Dim(s string) string { return ansiDim + s + ansiReset }

// Color wraps s in a given ANSI color code.
func Color(s, code string) string { return code + s + ansiReset }

// ColorSize color-codes size buckets: ≥1GiB red, ≥100MiB yellow, else green.
func ColorSize(s string, size int64) string {
	switch {
	case size >= 1<<30: // ≥ 1 GiB
		return Color(s, ansiRed)
	case size >= 100<<20: // ≥ 100 MiB
		return Color(s, ansiYellow)
	default:
		return Color(s, ansiGreen)
	}
}

// PathWithIcon prefixes files with an icon and gives a subtle color to the path.
func PathWithIcon(p string) string {
	return "🗎 " + Color(p, ansiBlue)
}
