package junk

import (
	"os"
	"path/filepath"
	"runtime"
)

type Suggestion struct {
	Path string `json:"path"`
	Note string `json:"note"`
}

func Suggest() []Suggestion {
	var out []Suggestion
	home, _ := os.UserHomeDir()
	switch runtime.GOOS {
	case "darwin":
		out = append(out,
			Suggestion{Path: "/private/var/folders", Note: "macOS app caches"},
			Suggestion{Path: filepath.Join(home, "Library", "Caches"), Note: "User caches"},
			Suggestion{Path: filepath.Join(home, "Library", "Logs"), Note: "Old logs"},
		)
	case "linux":
		out = append(out,
			Suggestion{Path: "/var/tmp", Note: "System temp"},
			Suggestion{Path: filepath.Join(home, ".cache"), Note: "User cache"},
			Suggestion{Path: filepath.Join(home, ".local", "share", "Trash"), Note: "Trash"},
		)
	case "windows":
		out = append(out,
			Suggestion{Path: filepath.Join(home, "AppData", "Local", "Temp"), Note: "User temp"},
			Suggestion{Path: "C:\\Windows\\Temp", Note: "System temp"},
			Suggestion{Path: filepath.Join(home, "AppData", "Local", "Microsoft", "Windows", "INetCache"), Note: "IE/Edge cache"},
		)
	}
	// common browser caches (cross-check per OS)
	out = append(out,
		Suggestion{Path: filepath.Join(home, ".cache", "google-chrome"), Note: "Chrome cache (Linux)"},
		Suggestion{Path: filepath.Join(home, "AppData", "Local", "Google", "Chrome", "User Data", "Default", "Cache"), Note: "Chrome cache (Windows)"},
	)
	return out
}
