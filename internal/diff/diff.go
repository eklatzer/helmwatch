package diff

import (
	"os/exec"
	"strings"
)

func Dirs(dir1, dir2 string) string {
	cmd := exec.Command("diff", "-ru", dir1, dir2)
	out, _ := cmd.CombinedOutput()
	if len(out) == 0 {
		return ""
	}

	return colorizeDiff(string(out))
}

func colorizeDiff(diff string) string {
	var out strings.Builder
	lines := strings.Split(diff, "\n")

	for _, line := range lines {
		switch {
		case strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++"):
			out.WriteString("\033[32m" + line + "\033[0m\n")
		case strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---"):
			out.WriteString("\033[31m" + line + "\033[0m\n")
		case strings.HasPrefix(line, "@@"):
			out.WriteString("\033[33m" + line + "\033[0m\n")
		case strings.HasPrefix(line, "+++"), strings.HasPrefix(line, "---"):
			out.WriteString("\033[34m" + line + "\033[0m\n")
		default:
			out.WriteString(line + "\n")
		}
	}

	return out.String()
}
