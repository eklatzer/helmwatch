package helm

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	releaseName = "helmwatch"
)

type TemplateOptions struct {
	Chart      string
	Version    string
	ValuesFile string
	Exclusions []string
}

func Template(options TemplateOptions) (string, error) {
	dir, err := os.MkdirTemp("", "helmwatch-*")
	if err != nil {
		return "", err
	}

	args := []string{"template", releaseName, options.Chart, "-f", options.ValuesFile, "--output-dir", dir}

	if options.Version != "" {
		args = append(args, "--version", options.Version)
	}

	cmd := exec.Command("helm", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to template chart: %w\n%s", err, string(out))
	}

	if len(options.Exclusions) > 0 {
		err := applyExclusions(dir, options.Exclusions)
		if err != nil {
			return "", err
		}
	}

	return dir, nil
}

func applyExclusions(dir string, exclusions []string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}

		lines := strings.Split(string(data), "\n")
		var filtered []string

		for _, line := range lines {
			exclude := false
			for _, pattern := range exclusions {
				if strings.Contains(line, pattern) {
					exclude = true
					break
				}
			}

			if !exclude {
				filtered = append(filtered, line)
			}
		}

		return os.WriteFile(path, []byte(strings.Join(filtered, "\n")), info.Mode())
	})
}
