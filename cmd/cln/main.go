package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

const version = "0.1.0"

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}
	cmd := os.Args[1]
	args := os.Args[2:]
	switch cmd {
	case "scan":
		runScan(args)
	case "dup":
		runDup(args)
	case "suggest":
		runSuggest(args)
	case "rm":
		runRm(args)
	case "version", "--version", "-v":
		fmt.Println("cln", version)
	case "--help", "-h", "help":
		usage()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", cmd)
		usage()
		os.Exit(1)
	}
}

func usage() {
	fmt.Println(`cln — Tiny Disk Cleanup CLI
Usage:
  cln scan [PATH ...] --min 100MB [--json]
  cln dup  [PATH ...] --min 10MB  [--json]
  cln suggest [--json]
  cln rm [FILE ...] [--force] [--dry]
  cln version`)
}

// parseSize like "100MB", "1GB", "500K"
func parseSize(s string) (int64, error) {
	if s == "" {
		return 0, nil
	}
	units := map[string]int64{
		"B": 1, "K": 1024, "KB": 1024,
		"M": 1024 * 1024, "MB": 1024 * 1024,
		"G": 1024 * 1024 * 1024, "GB": 1024 * 1024 * 1024,
	}
	su := strings.TrimSpace(strings.ToUpper(s))
	for k, mul := range units {
		if strings.HasSuffix(su, k) {
			num := strings.TrimSuffix(su, k)
			var v float64
			if _, err := fmt.Sscanf(num, "%f", &v); err != nil {
				return 0, err
			}
			return int64(v * float64(mul)), nil
		}
	}
	var v int64
	_, err := fmt.Sscanf(su, "%d", &v)
	return v, err
}

func defaultPaths(paths []string) []string {
	if len(paths) > 0 {
		return paths
	}
	return []string{"."}
}

func runScan(args []string) {
	fs := flag.NewFlagSet("scan", flag.ExitOnError)
	minStr := fs.String("min", "100MB", "minimum size threshold (e.g., 100MB)")
	asJSON := fs.Bool("json", false, "output JSON")
	_ = fs.Parse(args)

	minBytes, err := parseSize(*minStr)
	if err != nil {
		fmt.Fprintln(os.Stderr, "bad --min:", err)
		os.Exit(1)
	}
	roots := defaultPaths(fs.Args())
	ctx := context.Background()
	ch, _ := scan.WalkLarge(ctx, roots, minBytes)

	var rows [][]string
	var items []map[string]any
	for fi := range ch {
		if *asJSON {
			items = append(items, map[string]any{
				"path": fi.Path, "size": fi.Size, "mod_time": fi.ModTime,
			})
		} else {
			rows = append(rows, []string{
				format.HumanBytes(fi.Size),
				fi.ModTime.Format(time.RFC3339),
				fi.Path,
			})
		}
	}
	if *asJSON {
		_ = format.PrintJSON(os.Stdout, items)
	} else {
		format.PrintTable(os.Stdout, []string{"SIZE", "MODTIME", "PATH"}, rows)
	}
}

func runDup(args []string) {
	fs := flag.NewFlagSet("dup", flag.ExitOnError)
	minStr := fs.String("min", "10MB", "minimum size threshold (e.g., 10MB)")
	asJSON := fs.Bool("json", false, "output JSON")
	chunk := fs.Int("chunk", 4*1024*1024, "hash chunk size in bytes")
	_ = fs.Parse(args)

	minBytes, err := parseSize(*minStr)
	if err != nil {
		fmt.Fprintln(os.Stderr, "bad --min:", err)
		os.Exit(1)
	}
	roots := defaultPaths(fs.Args())

	// phase 1: collect files >= min and group by size
	ctx := context.Background()
	ch, _ := scan.WalkLarge(ctx, roots, minBytes)
	bySize := map[int64][]string{}
	for fi := range ch {
		bySize[fi.Size] = append(bySize[fi.Size], fi.Path)
	}

	type dupGroup struct {
		Size  int64    `json:"size"`
		Files []string `json:"files"`
		Hash  string   `json:"hash"`
	}
	var groups []dupGroup

	for size, files := range bySize {
		if len(files) < 2 {
			continue
		}
		// compute hash per file
		hashes := map[string][]string{}
		for _, p := range files {
			h, err := hash.FileSHA256(p, *chunk)
			if err != nil {
				continue
			}
			hashes[h] = append(hashes[h], p)
		}
		for h, fslist := range hashes {
			if len(fslist) > 1 {
				sort.Strings(fslist)
				groups = append(groups, dupGroup{Size: size, Files: fslist, Hash: h})
			}
		}
	}

	sort.Slice(groups, func(i, j int) bool {
		if groups[i].Size == groups[j].Size {
			return groups[i].Hash < groups[j].Hash
		}
		return groups[i].Size > groups[j].Size
	})

	if *asJSON {
		_ = format.PrintJSON(os.Stdout, groups)
		return
	}

	var rows [][]string
	for _, g := range groups {
		rows = append(rows, []string{format.HumanBytes(g.Size), g.Hash[:12], g.Files[0]})
		for _, f := range g.Files[1:] {
			rows = append(rows, []string{"", "", f})
		}
	}
	if len(rows) == 0 {
		fmt.Println("No duplicates found.")
		return
	}
	format.PrintTable(os.Stdout, []string{"SIZE", "HASH", "FILE"}, rows)
}

func runSuggest(args []string) {
	fs := flag.NewFlagSet("suggest", flag.ExitOnError)
	asJSON := fs.Bool("json", false, "output JSON")
	_ = fs.Parse(args)

	sugg := junk.Suggest()
	if *asJSON {
		_ = format.PrintJSON(os.Stdout, sugg)
		return
	}
	var rows [][]string
	for _, s := range sugg {
		rows = append(rows, []string{s.Path, s.Note})
	}
	format.PrintTable(os.Stdout, []string{"PATH", "NOTE"}, rows)
}

func runRm(args []string) {
	fs := flag.NewFlagSet("rm", flag.ExitOnError)
	force := fs.Bool("force", false, "do not prompt")
	dry := fs.Bool("dry", false, "dry run (no deletion)")
	_ = fs.Parse(args)

	files := fs.Args()
	if len(files) == 0 {
		fmt.Fprintln(os.Stderr, "rm: provide at least one file")
		os.Exit(1)
	}
	if !*force {
		fmt.Printf("About to remove %d file(s). Proceed? [y/N]: ", len(files))
		rd := bufio.NewReader(os.Stdin)
		line, _ := rd.ReadString('\n')
		line = strings.TrimSpace(strings.ToLower(line))
		if line != "y" && line != "yes" {
			fmt.Println("Aborted.")
			return
		}
	}
	var errs []error
	for _, f := range files {
		if *dry {
			fmt.Println("[dry] rm", f)
			continue
		}
		if err := os.Remove(f); err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", f, err))
		} else {
			fmt.Println("removed", f)
		}
	}
	if len(errs) > 0 {
		fmt.Fprintln(os.Stderr, "Some errors:")
		for _, e := range errs {
			fmt.Fprintln(os.Stderr, " -", e)
		}
		os.Exit(1)
	}
}
