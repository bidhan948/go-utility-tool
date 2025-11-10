package format

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
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
