package xpath

import (
	"os"
	"path/filepath"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
)

func Abs(inPath string) string {
	home, err := homedir.Dir()

	if err != nil {
		return ""
	}

	if strings.HasPrefix(inPath, "$HOME") {
		inPath = home + inPath[5:]
	}

	if strings.HasPrefix(inPath, "$") {
		end := strings.Index(inPath, string(os.PathSeparator))
		inPath = os.Getenv(inPath[1:end]) + inPath[end:]
	}

	if filepath.IsAbs(inPath) {
		return filepath.Clean(inPath)
	}

	p, err := filepath.Abs(inPath)

	if err != nil {
		return ""
	}

	return filepath.Clean(p)
}
