package main

import (
	"bufio"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	pwl "github.com/justjanne/powerline-go/powerline"
)

var (
	pythonVersionRe    = regexp.MustCompile(`^\d+\.\d+(\.\d+)?`)
	pythonExeVersionRe = regexp.MustCompile(`^python(\d+\.\d+)$`)
	pythonPatchlevelRe = regexp.MustCompile(`^#define\s+PY_VERSION\s+"([^"]+)"`)
)

// pythonVersionFromPyvenvCfg reads the version recorded in a virtualenv's
// pyvenv.cfg ("version" for the venv module, "version_info" for uv and
// virtualenv). It returns "" when there is no such file or key.
func pythonVersionFromPyvenvCfg(path string) string {
	f, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		key, value, ok := strings.Cut(scanner.Text(), "=")
		if !ok {
			continue
		}
		switch strings.TrimSpace(key) {
		case "version", "version_info":
			return pythonVersionRe.FindString(strings.TrimSpace(value))
		}
	}
	return ""
}

// pythonVersionFromPatchlevel reads PY_VERSION from a patchlevel.h header
// (only present when the Python development headers are installed).
func pythonVersionFromPatchlevel(path string) string {
	f, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if m := pythonPatchlevelRe.FindStringSubmatch(scanner.Text()); m != nil {
			return m[1]
		}
	}
	return ""
}

// pythonVersion determines the version of the given python3 executable
// without running it. It tries, in order:
//  1. pyvenv.cfg when the executable lives in a virtualenv (full version)
//  2. include/python3.X/patchlevel.h next to the real interpreter (full version)
//  3. the name of the real interpreter, e.g. python3.14 (major.minor only)
func pythonVersion(exe string) string {
	if v := pythonVersionFromPyvenvCfg(filepath.Join(filepath.Dir(filepath.Dir(exe)), "pyvenv.cfg")); v != "" {
		return v
	}

	resolved, err := filepath.EvalSymlinks(exe)
	if err != nil {
		return ""
	}
	prefix := filepath.Dir(filepath.Dir(resolved))

	minor := ""
	if m := pythonExeVersionRe.FindStringSubmatch(filepath.Base(resolved)); m != nil {
		minor = m[1]
	}

	if minor != "" {
		if v := pythonVersionFromPatchlevel(filepath.Join(prefix, "include", "python"+minor, "patchlevel.h")); v != "" {
			return v
		}
		return minor
	}

	// The interpreter isn't named python3.X; use the headers if there is
	// exactly one set of them under its prefix.
	headers, _ := filepath.Glob(filepath.Join(prefix, "include", "python3.*", "patchlevel.h"))
	if len(headers) == 1 {
		return pythonVersionFromPatchlevel(headers[0])
	}
	return ""
}

// segmentPython shows the version of the python3 found on PATH. Its
// background is darker when that python3 lives in ~/.local/bin.
func segmentPython(p *powerline) []pwl.Segment {
	exe, err := exec.LookPath("python3")
	if err != nil {
		return []pwl.Segment{}
	}
	version := pythonVersion(exe)
	if version == "" {
		return []pwl.Segment{}
	}

	bg := p.theme.PythonBg
	if filepath.Base(filepath.Dir(filepath.Dir(exe))) == ".local" {
		bg = p.theme.PythonLocalBg
	}

	return []pwl.Segment{{
		Name:       "python",
		Content:    version,
		Foreground: p.theme.PythonFg,
		Background: bg,
	}}
}
