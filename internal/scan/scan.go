package scan

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

type FileInfo struct {
	Path    string    `json:"path"`
	Size    int64     `json:"size"`
	ModTime time.Time `json:"mod_time"`
}

func WalkLarge(ctx context.Context, roots []string, minBytes int64) (<-chan FileInfo, <-chan error) {
	out := make(chan FileInfo, 128)
	errc := make(chan error, 1)

	go func() {
		defer close(out)
		defer close(errc)

		workers := runtime.NumCPU()
		paths := make(chan string, 1024)

		var wg sync.WaitGroup
		for i := 0; i < workers; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for p := range paths {
					select {
					case <-ctx.Done():
						return
					default:
					}
					fi, err := os.Lstat(p)
					if err != nil {
						continue
					}
					if fi.Mode()&os.ModeSymlink != 0 {
						// skip symlinks
						continue
					}
					if fi.Mode().IsRegular() && fi.Size() >= minBytes {
						out <- FileInfo{Path: p, Size: fi.Size(), ModTime: fi.ModTime()}
					}
				}
			}()
		}

		walkFn := func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return nil // skip unreadable
			}
			if d.Type()&os.ModeSymlink != 0 {
				return nil
			}
			if d.IsDir() {
				return nil
			}
			select {
			case <-ctx.Done():
				return errors.New("canceled")
			case paths <- path:
				return nil
			}
		}

		for _, r := range roots {
			_ = filepath.WalkDir(r, walkFn)
		}
		close(paths)
		wg.Wait()
	}()

	return out, errc
}
